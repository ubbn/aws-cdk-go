package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	apigateway "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	authorizers "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2authorizers"
	integrations "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	prefix        = "jwt-"
	defaultRegion = "eu-west-1"
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

	// Http Apigateway
	apigw := apigateway.NewHttpApi(stack, jsii.String(prefix+"authorized-http-apigateway"), &apigateway.HttpApiProps{})

	// Create lambda authorizer function
	authorizer := createFunc(stack, "authorizer", "lambda/auth")

	lambdaAUthorizer := authorizers.NewHttpLambdaAuthorizer(jsii.String("authorizer-func"), authorizer,
		&authorizers.HttpLambdaAuthorizerProps{ResponseTypes: &[]authorizers.HttpLambdaResponseType{
			authorizers.HttpLambdaResponseType_SIMPLE,
		}})

	handler := createFunc(stack, "handler", "lambda/app")

	functionIntg := integrations.NewHttpLambdaIntegration(jsii.String(prefix+"integration"), handler,
		&integrations.HttpLambdaIntegrationProps{
			Timeout: awscdk.Duration_Seconds(jsii.Number(25)),
		})

	apigw.AddRoutes(&apigateway.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]apigateway.HttpMethod{apigateway.HttpMethod_POST},
		Integration: functionIntg,
		Authorizer:  lambdaAUthorizer,
	})

	// Output API gateway URL
	awscdk.NewCfnOutput(
		stack, jsii.String(prefix+"apigw URL"),
		&awscdk.CfnOutputProps{Value: apigw.Url(),
			Description: jsii.String("API Gateway endpoint")},
	)

	return stack
}

// Create aws lambda function from given source adn with given name
func createFunc(stack awscdk.Stack, name string, sourcePath string) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String(prefix+name),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Code: awslambda.Code_FromAsset(jsii.String(sourcePath), &awss3assets.AssetOptions{
				Bundling: &awscdk.BundlingOptions{
					Image: awslambda.Runtime_PROVIDED_AL2023().BundlingImage(),
					Environment: &map[string]*string{
						"GOCACHE": jsii.String("/tmp/.cache"),
						"GOPATH":  jsii.String("/tmp/go"),
					},
					Command: &[]*string{
						jsii.String("bash"),
						jsii.String("-c"),
						jsii.String("go build -o /asset-output"),
					},
				},
			}),
			Handler: jsii.String("index.main"),
			Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
		},
	)
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "ValidatorStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String(defaultRegion),
	}

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
