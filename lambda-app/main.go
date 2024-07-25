package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	prompt := req.Body
	log.Println("input -", prompt)

	return events.APIGatewayV2HTTPResponse{
		StatusCode:      http.StatusOK,
		Body:            "image",
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST,OPTIONS",
		},
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
