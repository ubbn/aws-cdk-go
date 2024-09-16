package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Create struct to hold info about new item
type Item struct {
	RequestId string
	SourceIp  string
	Auth      string
	Epoch     int64
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	log.Println("Full req: ", req)
	log.Println("Full ctx: ", req.RequestContext)
	log.Println("RequestId: ", req.RequestContext.RequestID)
	log.Println("SourceIP: ", req.RequestContext.HTTP.SourceIP)
	log.Println("Epoch: ", req.RequestContext.TimeEpoch)
	log.Println("Authentication part: ", req.Headers["authorization"])

	tableName := os.Getenv("JWTLOG")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	item := Item{
		RequestId: req.RequestContext.RequestID,
		SourceIp:  req.RequestContext.HTTP.SourceIP,
		Auth:      req.Headers["authorization"],
		Epoch:     req.RequestContext.TimeEpoch,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		// Context: map[string]interface{}{
		// 	"Access-Control-Allow-Origin":  "*",
		// 	"Access-Control-Allow-Methods": "POST,OPTIONS",
		// },
	}, nil
}

func main() {
	lambda.Start(handler)
}

type Request struct {
	TextPrompts []TextPrompt `json:"text_prompts"`
	CfgScale    float64      `json:"cfg_scale"`
	Steps       int          `json:"steps"`
	Seed        int          `json:"seed"`
}

type TextPrompt struct {
	Text string `json:"text"`
}

type Response struct {
	Result    string     `json:"result"`
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	Base64       string `json:"base64"`
	FinishReason string `json:"finishReason"`
}
