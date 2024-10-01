module bootstrap

go 1.22.3

replace ubbn.com/dynamodb => ../dynamodb

require (
	github.com/aws/aws-lambda-go v1.47.0
	ubbn.com/dynamodb v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
