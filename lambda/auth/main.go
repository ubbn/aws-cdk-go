package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	log.Println("Full req: ", req)
	log.Println("Full ctx: ", req.RequestContext)
	log.Println("RequestId: ", req.RequestContext.RequestID)
	log.Println("SourceIP: ", req.RequestContext.HTTP.SourceIP)
	log.Println("Epoch: ", req.RequestContext.TimeEpoch)
	log.Println("Authentication part: ", req.Headers["authorization"])

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
