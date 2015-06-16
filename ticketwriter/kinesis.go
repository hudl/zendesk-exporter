package ticketwriter

import (
	kinesis "github.com/sendgridlabs/go-kinesis"
)

// uploadToS3 uploads the contents of tickWrt.currenFileName to S3
func (tickWrt *TicketWriter) uploadToKinesis(data []byte, partitionKey string) {
	auth := kinesis.NewAuth(tickWrt.awsAuth.AccessKey, tickWrt.awsAuth.SecretKey)
	ksis := kinesis.New(auth, kinesis.USEast1)
	args := kinesis.NewArgs()
	args.AddRecord(data, partitionKey)
	args.Add("StreamName", tickWrt.ksisStreamName)
	resp, err := ksis.PutRecord(args)
	if err != nil {
		log.Error("Error uploading record to Kinesis: %+v. Response: %+v", err, resp)
		return
	}
	log.Info("Response from Kinesis: %+v", resp)
}
