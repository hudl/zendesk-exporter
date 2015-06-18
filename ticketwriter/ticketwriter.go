// Package ticketwriter is responsible for writing tickets to their destination.
// Initially the destination is a combination of a local file, which then gets
// written to S3. More possible destinations (Kinesis, Sumologic) will be added
// in the future.
package ticketwriter

import (
	"github.com/goamz/goamz/aws"
	"github.com/hudl/ZeGo/zego"

	"encoding/json"
	"fmt"
	"os"
	"time"
)

// A TicketWriter contains information needed to write tickets to files
// and to S3. This includes AWS	auth info, bucket locations, and filename
type TicketWriter struct {
	filePrefix string

	awsAuth        aws.Auth
	s3BucketName   string
	s3KeyPrefix    string
	ksisStreamName string

	maxFileSz       uint64
	currentFileName string
}

// New(...) creates a new instance of a TicketWriter from the given parameters.
// Using this method instead of using a struct initializer prevents clients from
// having to import github.com/goamz/aws
func New(maxFileSz uint64, filePrefix, s3BucketName, s3KeyPrefix, awsAccessKey, awsSecretKey, ksisStream string) TicketWriter {
	tickWriter := TicketWriter{}
	tickWriter.filePrefix = filePrefix

	auth, err := aws.GetAuth(awsAccessKey, awsSecretKey, "", time.Time{})
	if err != nil {
		log.Error("Error trying to get AWS Auth: %+v", err)
		panic(err)
	}
	tickWriter.awsAuth = auth
	tickWriter.s3BucketName = s3BucketName
	tickWriter.s3KeyPrefix = s3KeyPrefix
	tickWriter.maxFileSz = maxFileSz
	tickWriter.ksisStreamName = ksisStream

	return tickWriter
}

// Write takes a slice of tickets and their startTime and writes them to a local file.
// If that file grows above a certain threshold (tickWrt.maxFileSz) then the file is
// uploaded to S3, deleted locally, then a new file will be created the next time Write
// is called.
func (tickWrt *TicketWriter) Write(ticks []zego.Ticket, startTime string) {
	if tickWrt.currentFileName == "" {
		tickWrt.currentFileName = tickWrt.filePrefix + startTime + ".json"
	}
	tickFile, _ := os.OpenFile(tickWrt.currentFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	for _, t := range ticks {
		jsonBytes, err := json.Marshal(t)
		if err != nil {
			log.Error("Error marshalling ticket to json. Ticket: %+v \nError: %+v", t, err)
			continue
		}
		_, err = tickFile.WriteString(string(jsonBytes))
		if err != nil {
			log.Error("Error writing to file: %+v", err)
		}
		tickFile.WriteString("\n")

		tickWrt.uploadToKinesis(jsonBytes, "partionKey-arbitrary")

		logString := fmt.Sprintf("Zendesk Ticket Id=%d Tags=%v CreatedAt=%q Subject=%q Submitter=%v Assignee=%v Group_Id=%v",
			t.Id, t.Tags, t.CreatedAt, t.Subject, t.SubmitterId, t.AssigneeId, t.GroupId)
		for _, field := range t.Custom_Fields {
			logString += fmt.Sprintf(" %d=%q", field.Id, field.Value)
		}

		logString += fmt.Sprintf(" Channel=%q From=%q To=%q", t.Via.Channel, t.Via.Source.From, t.Via.Source.To)
		log.Info(logString)
	}

	// Get size of file to find out if we need to upload and delete it.
	st, _ := tickFile.Stat()
	fileSz := uint64(st.Size())
	log.Info("File is %dB large", fileSz)

	if fileSz >= tickWrt.maxFileSz {
		log.Info("File is above max size. Uploading to S3")
		tickFile.Close()
		tickWrt.uploadToS3()

		os.Remove(tickWrt.currentFileName)

		//Clear current file name so the next one is unique
		tickWrt.currentFileName = ""
	}
}
