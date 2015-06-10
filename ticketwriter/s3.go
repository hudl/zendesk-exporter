package ticketwriter

import (
	"bufio"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"net/http"

	"os"
)

// uploadToS3 uploads the contents of tickWrt.currenFileName to S3
func (tickWrt *TicketWriter) uploadToS3() {
	s3Conn := s3.New(tickWrt.awsAuth, aws.USEast)
	bucket := s3Conn.Bucket(tickWrt.s3BucketName)
	path := tickWrt.s3KeyPrefix + tickWrt.currentFileName

	file, err := os.Open(tickWrt.currentFileName)
	if err != nil {
		log.Error("Error opening file %s for upload to s3: %+v", tickWrt.currentFileName, err)
		panic(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	//read into buffer
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	filetype := http.DetectContentType(bytes)
	err = bucket.Put(path, bytes, filetype, s3.AuthenticatedRead, s3.Options{})
	if err != nil {
		log.Error("Error uploading file to S3: %+v", err)
		panic(err)
	}
}
