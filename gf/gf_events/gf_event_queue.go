// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package gf_events

import (
    "fmt"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

//-----------------------------------------------------------------
type GFevenstQueueInfo struct {
    awsSQSclient   *sqs.SQS
    awsSQSqueueURL string
}

//-----------------------------------------------------------------
func queueSQSinit() (*GFevenstQueueInfo, error) {



	sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    svc := sqs.New(sess)

    // URL to our queue
    qURL := "QueueURL"




    queueInfo := &GFevenstQueueInfo{
        awsSQSclient:   svc,
        awsSQSqueueURL: qURL,
    }

    return queueInfo, nil
}

//-----------------------------------------------------------------
func pushEvent(pEvent GFeventMsg,
    pQueueInfo *GFevenstQueueInfo) error {

	result, err := pQueueInfo.awsSQSclient.SendMessage(&sqs.SendMessageInput{
        DelaySeconds:      aws.Int64(10),
        MessageAttributes: map[string]*sqs.MessageAttributeValue{

			"time_sec": &sqs.MessageAttributeValue{
                DataType:    aws.String("String"),
                StringValue: aws.String(fmt.Sprint(pEvent.TimeSec)),
            },
			"module": &sqs.MessageAttributeValue{
                DataType:    aws.String("String"),
                StringValue: aws.String(pEvent.Module),
            },
            "type": &sqs.MessageAttributeValue{
                DataType:    aws.String("String"),
                StringValue: aws.String(pEvent.Type),
            },
            "msg": &sqs.MessageAttributeValue{
                DataType:    aws.String("String"),
                StringValue: aws.String(pEvent.Msg),
            },
        },
        MessageBody: aws.String(""),
        QueueUrl:    &pQueueInfo.awsSQSqueueURL,
    })

    if err != nil {
        return err
    }

    fmt.Println("Success", *result.MessageId)

	return nil
}