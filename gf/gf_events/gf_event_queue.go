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
    "encoding/json"
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
func queueSQSinit(pSQSqueueName string) (*GFevenstQueueInfo, error) {



    //----------------------------
    // SESSION
	/*sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))*/

    sess, err := session.NewSession()
    if err != nil {
        return nil, err
    }
    svc := sqs.New(sess)

    //----------------------------
    // GET_QUEUE_URL
    result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
        QueueName: aws.String(pSQSqueueName),
    })
    if err != nil {
        return nil, err
    }
    SQSqueueURL := *result.QueueUrl

    //----------------------------

    queueInfo := &GFevenstQueueInfo{
        awsSQSclient:   svc,
        awsSQSqueueURL: SQSqueueURL,
    }

    return queueInfo, nil
}

//-----------------------------------------------------------------
func queueSQSpushEvent(pEvent interface{},
    pQueueInfo *GFevenstQueueInfo) error {


    eventDataJSONencoded, err := json.Marshal(pEvent)
    if err != nil {
        return err
    }

    fmt.Println("SENDING AWS SQS msg----")
	result, err := pQueueInfo.awsSQSclient.SendMessage(&sqs.SendMessageInput{
        MessageBody: aws.String(string(eventDataJSONencoded)),
        QueueUrl:    &pQueueInfo.awsSQSqueueURL,
        // DelaySeconds: aws.Int64(10),
        /*MessageAttributes: map[string]*sqs.MessageAttributeValue{

			"time_sec": &sqs.MessageAttributeValue{
                DataType:    aws.String("String"),
                StringValue: aws.String(fmt.Sprintf("%f", pEvent.TimeSec)),
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
        },*/
    })

    if err != nil {
        fmt.Println("FAILED TO SEND SQS MSG")
        fmt.Println(fmt.Sprint(err))
        return err
    }

    fmt.Println("Success", *result.MessageId)

	return nil
}