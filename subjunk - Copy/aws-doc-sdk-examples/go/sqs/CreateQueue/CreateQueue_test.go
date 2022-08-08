/*
   Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/

package main

import (
    "encoding/json"
    "io/ioutil"
    "testing"

    "github.com/google/uuid"

    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

// Config defines a set of configuration values
type Config struct {
    Queue string `json:"Queue"`
}

// configFile defines the name of the file containing configuration values
var configFileName = "config.json"

// globalConfig contains the configuration values
var globalConfig Config

func populateConfiguration(t *testing.T) error {
    // Get configuration from config.json

    // Get entire file as a JSON string
    content, err := ioutil.ReadFile(configFileName)
    if err != nil {
        return err
    }

    // Convert []byte to string
    text := string(content)

    // Marshall JSON string in text into global struct
    err = json.Unmarshal([]byte(text), &globalConfig)
    if err != nil {
        return err
    }

    t.Log("Queue: " + globalConfig.Queue)

    return nil
}

func deleteQueue(sess *session.Session, queueURL *string) error {
    // Create an SQS service client
    svc := sqs.New(sess)

    _, err := svc.DeleteQueue(&sqs.DeleteQueueInput{
        QueueUrl: queueURL,
    })
    if err != nil {
        return err
    }

    return nil
}

func TestCreateQueue(t *testing.T) {
    err := populateConfiguration(t)
    if err != nil {
        t.Fatal(err)
    }

    queueCreated := false

    if globalConfig.Queue == "" {
        // Create a unique, random queue name
        id := uuid.New()
        globalConfig.Queue = "myqueue-" + id.String()
        queueCreated = true
    }

    // Create a session using credentials from ~/.aws/credentials
    // and the Region from ~/.aws/config
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    result, err := CreateQueue(sess, &globalConfig.Queue)
    if err != nil {
        t.Fatal(err)
    }

    t.Log("Got URL " + *result.QueueUrl + " for queue " + globalConfig.Queue)

    if queueCreated {
        err = deleteQueue(sess, result.QueueUrl)
        if err != nil {
            t.Log("You'll have to delete queue " + globalConfig.Queue + " yourself")
            t.Fatal(err)
        }

        t.Log("Deleted queue " + globalConfig.Queue)
    }
}
