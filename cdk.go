package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	authorizers "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2authorizers"
	integrations "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	dynamodb "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	prefix        = "jwt-"
	defaultRegion = "eu-west-1"
	burstLimit    = 50  // Limit concurrent requests processing
	rateLimit     = 100 // Limit incoming requests when burst is taking an effect
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := dynamodb.NewTableV2(stack, jsii.String(prefix+"history"), &dynamodb.TablePropsV2{
		PartitionKey: &dynamodb.Attribute{
			Name: jsii.String("RequestId"),
			Type: dynamodb.AttributeType_STRING,
		},
		RemovalPolicy: "DESTROY",
	})

	envVars := map[string]*string{
		"JWTTABLE": table.TableName(),
	}

	// Http Apigateway
	apigw := apigateway.NewHttpApi(stack, jsii.String(prefix+"authorized-http-apigateway"), &apigateway.HttpApiProps{
		CreateDefaultStage: jsii.Bool(false), // Create own stage below
	})

	// Add new stage, instead of default one
	stage := apigw.AddStage(jsii.String(prefix+"stage"), &apigateway.HttpStageOptions{
		StageName:   jsii.String("prod"),
		AutoDeploy:  jsii.Bool(true),
		Description: jsii.String("Production stage"),
		Throttle: &apigateway.ThrottleSettings{
			BurstLimit: jsii.Number(burstLimit),
			RateLimit:  jsii.Number(rateLimit),
		},
	})

	// Create lambda authorizer function
	authorizer := createFunc(stack, "authorizer", "auth", &envVars)

	lambdaAUthorizer := authorizers.NewHttpLambdaAuthorizer(jsii.String("authorizer-func"), authorizer,
		&authorizers.HttpLambdaAuthorizerProps{ResponseTypes: &[]authorizers.HttpLambdaResponseType{
			authorizers.HttpLambdaResponseType_SIMPLE,
		}})

	handler := createFunc(stack, "handler", "app", &envVars)

	functionIntg := integrations.NewHttpLambdaIntegration(jsii.String(prefix+"integration"), handler,
		&integrations.HttpLambdaIntegrationProps{
			Timeout: awscdk.Duration_Seconds(jsii.Number(25)),
		})

	apigw.AddRoutes(&apigateway.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]apigateway.HttpMethod{apigateway.HttpMethod_GET, apigateway.HttpMethod_POST},
		Integration: functionIntg,
		Authorizer:  lambdaAUthorizer,
	})

	table.GrantReadWriteData(authorizer)
	table.GrantReadWriteData(handler)

	// Output API gateway URL
	awscdk.NewCfnOutput(
		stack, jsii.String(prefix+"apigw URL"),
		&awscdk.CfnOutputProps{Value: stage.Url(),
			Description: jsii.String("API Gateway endpoint")},
	)

	// Output DynamoDB name
	awscdk.NewCfnOutput(
		stack, jsii.String(prefix+"dynamodb-name"),
		&awscdk.CfnOutputProps{Value: table.TableName(),
			Description: jsii.String("DynamoDB name")},
	)
	return stack
}

// Create aws lambda function from given source adn with given name
func createFunc(stack awscdk.Stack, name string, sourcePath string, env *map[string]*string) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String(prefix+name),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Code: awslambda.Code_FromAsset(jsii.String("lambda"), &awss3assets.AssetOptions{
				Bundling: &awscdk.BundlingOptions{
					Image: awslambda.Runtime_PROVIDED_AL2023().BundlingImage(),
					Environment: &map[string]*string{
						"GOCACHE": jsii.String("/tmp/.cache"),
						"GOPATH":  jsii.String("/tmp/go"),
					},
					Command: &[]*string{
						jsii.String("bash"),
						jsii.String("-c"),
						jsii.String("cd " + sourcePath + " && go build -o /asset-output"),
					},
				},
			}),
			Handler:     jsii.String("index.main"),
			Timeout:     awscdk.Duration_Seconds(jsii.Number(30)),
			Environment: env,
		},
	)
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(&awscdk.AppProps{})

	NewCdkStack(app, prefix+"stack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment region in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// Read region from custom environment variable AWS_REGION
	region, isPresent := os.LookupEnv("AWS_REGION")
	if !isPresent {
		// Read region of chosen cdk CLI profile or fallback to default profile
		region = defaultRegion
	}

	fmt.Println("Applying to region: ", region)

	return &awscdk.Environment{
		Region: jsii.String(region),
	}
}
