package main

import (
	"log"
	"os"
	"time"

        "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/sqs"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	url := os.Getenv("SQS_ENDPOINT_URL")

        sess,err := session.NewSession(&aws.Config{
                Endpoint:    aws.String(url),
		Region:      aws.String("RegionOne"),
		MaxRetries:  aws.Int(5),
	})
	failOnError(err, "Failed to connect to SQS API")

        svc := sqs.New(sess)

        queue, err := svc.CreateQueue(&sqs.CreateQueueInput{
              QueueName: aws.String("sample_app_queue"),
              Attributes: map[string]*string{
                  "DelaySeconds":           aws.String("0"),
                  "MessageRetentionPeriod": aws.String("86400"),
              },
        })
	failOnError(err, "Failed to create a queue")

	log.Printf("Waiting for messages")

	for {
            msg, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
                QueueUrl:    queue.QueueUrl,
                MaxNumberOfMessages: aws.Int64(1),
            })
	    failOnError(err, "Failed to get a message")

	    if len(msg.Messages) > 0 {
	        log.Printf("Received a message: %s and doing work for 10 seconds", *msg.Messages[0].Body)
		time.Sleep(10 * time.Second)
                _, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
                    QueueUrl:      queue.QueueUrl,
                    ReceiptHandle: msg.Messages[0].ReceiptHandle,
                })
	        failOnError(err, "Failed to delete a message")

	    }
        }
}
