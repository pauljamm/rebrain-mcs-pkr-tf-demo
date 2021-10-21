package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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


	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
	    body := "Do work"
	    _, err = svc.SendMessage(&sqs.SendMessageInput{
                  QueueUrl:    queue.QueueUrl,
                  MessageBody: aws.String(body),
            })
	    failOnError(err, "Failed to publish a message")
	    fmt.Fprintf(w, "Published")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "OK")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		failOnError(err, "Failed to start HTTP server")
	}
}
