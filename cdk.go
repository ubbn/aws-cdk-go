package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
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

	// The code that defines your stack goes here
	apigw := awsapigatewayv2.NewHttpApi(stack, jsii.String(prefix+"http-apigateway"), &awsapigatewayv2.HttpApiProps{})

	function := awslambda.NewFunction(stack, jsii.String(prefix+"func"),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Code: awslambda.Code_FromAsset(jsii.String("lambda/app"), &awss3assets.AssetOptions{
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

	functionIntg := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String(prefix+"integration"), function, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{})

	apigw.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_POST},
		Integration: functionIntg,
	})

	// Output API gateway URL
	awscdk.NewCfnOutput(
		stack, jsii.String(prefix+"apigw URL"),
		&awscdk.CfnOutputProps{Value: apigw.Url(),
			Description: jsii.String("API Gateway endpoint")},
	)

	return stack
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
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
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
