# Firebase JWT Handler

This project builds a cloud backend stack in AWS using [Golang CDK](https://docs.aws.amazon.com/cdk/v2/guide/work-with-cdk-go.html) for handling requests containing Firebase JWT tokens. A part of or entire codebase of the project can be used as a reference for building serverless authentication system in AWS that uses [Firebase Authentication](https://firebase.google.com/docs/auth) in Go programming language.

## Resource Topology

The stack provisions following AWS resources:
- [API Gateway](https://aws.amazon.com/api-gateway/)
- [Lambda authorizer function](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html)
- [Lambda serverless function](https://aws.amazon.com/lambda/)
- [DynamoDB](https://aws.amazon.com/dynamodb/)
- [Simple Email Service](https://aws.amazon.com/ses/)

![Resource topology in stack](stack-topology.svg "Stack Topology")

### Flow in Stack
- API Gateway recieves requests, also it terminates TLS and guards its backend with custom throttling. It forwards the request to authorizer function and responds back to API caller. 
- Lambda authorizer function validates authenitcation token in request's header. And 
- Lambda function processes the valid requests and writes content of token into DynamoDB, and sends notification email to owner of the token.

## Prepare for development

### Prerequisites tools

For developing and deploying resources in the project, you will need following prerequisites:

* **Go, v1.20 or newer**   
Follow instructions on its official website [https://go.dev/doc/install](https://go.dev/doc/install).

* **aws-cdk CLI tool**    
CDK npm package needs to be installed. Assuming nodejs is already installed, run below:     
`npm install -g cdk`

* **Docker**   
Docker image [AL2003](https://docs.aws.amazon.com/linux/al2023/ug/go.html), Amazon Linux Image for Go, is used for building lambda functions with Go toolchain. Read more about [Building lambda function in Go](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)   
For docker installation, follow offical guide at [https://docs.docker.com/engine/install](https://docs.docker.com/engine/install)

### Install dependencies

 * `go mod download`
 * `go get`

### Configure AWS environment

This CDK stack is deployed to AWS environment which is setup in local AWS configuration. For more details about environment configuration, refer to AWS official guide on [Environments for the AWS CDK](https://docs.aws.amazon.com/cdk/v2/guide/environments.html)

## Development

### Deploy the stack

 * `cdk bootstrap`   bootstrap stack, only run once at first time
 * `cdk deploy`      deploy this stack to AWS, run whenever code changes

### Other useful commands

 * `cdk diff`        compare deployed stack with current state
 * `cdk destroy`     clean up provisioned AWS resources
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests


### Troubleshoot

```bash
# Deploy the stack in hot-swap for live updates in case of code changes 
cdk watch

# Send a request and monitor its live logs from watch command's output
curl -X POST  https://9h3r3yi8mh.execute-api.eu-west-1.amazonaws.com/  --header "Authorization: Bearer abcdefgh"
```
