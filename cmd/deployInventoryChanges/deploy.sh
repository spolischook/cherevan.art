#!/bin/bash

# Set variables
LAMBDA_FUNCTION_NAME="cherevan_art_deploy"
AWS_REGION="eu-central-1"
ZIP_FILE="deploy.zip"
AWS_PROFILE="deploy-lambda"

# Build the Go binary for Linux (required for AWS Lambda)
echo "Building Go binary for Lambda..."
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# Zip the binary
echo "Zipping the binary..."
zip $ZIP_FILE bootstrap

# Deploy the Lambda function using AWS CLI
echo "Deploying to AWS Lambda..."
aws lambda update-function-code \
    --region $AWS_REGION \
    --function-name arn:aws:lambda:$AWS_REGION:805655309022:function:$LAMBDA_FUNCTION_NAME \
    --zip-file fileb://$ZIP_FILE \
    --profile $AWS_PROFILE

# Cleanup
echo "Cleaning up..."
rm bootstrap $ZIP_FILE

echo "Deployment complete."
