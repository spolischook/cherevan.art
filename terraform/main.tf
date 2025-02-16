terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region  = "eu-central-1"
  profile = var.aws_profile
}

provider "aws" {
  alias   = "eu-central-1"
  region  = "eu-central-1"
  profile = var.aws_profile
}

variable "aws_profile" {
  type        = string
  description = "AWS profile to use"
  default     = "personal"
}

variable "google_client_id" {
  type        = string
  description = "Google OAuth Client ID"
}

variable "google_client_secret" {
  type        = string
  description = "Google OAuth Client Secret"
}

variable "google_redirect_url" {
  type        = string
  description = "Google OAuth Redirect URL"
}

# IAM role for Lambda
resource "aws_iam_role" "lambda_role" {
  name = "google-auth-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Attach basic Lambda execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}

# Lambda function
resource "aws_lambda_function" "google_auth" {
  filename         = "../google-auth/function.zip"
  function_name    = "google-auth-lambda"
  role             = aws_iam_role.lambda_role.arn
  handler          = "main"
  source_code_hash = filebase64sha256("../google-auth/function.zip")
  runtime          = "provided.al2"

  environment {
    variables = {
      GOOGLE_CLIENT_ID     = var.google_client_id
      GOOGLE_CLIENT_SECRET = var.google_client_secret
      GOOGLE_REDIRECT_URL  = "https://auth.cherevan.art/callback"
    }
  }
}

# API Gateway
resource "aws_api_gateway_rest_api" "google_auth" {
  name = "google-auth-api"
}

# /auth endpoint
resource "aws_api_gateway_resource" "auth" {
  rest_api_id = aws_api_gateway_rest_api.google_auth.id
  parent_id   = aws_api_gateway_rest_api.google_auth.root_resource_id
  path_part   = "auth"
}

resource "aws_api_gateway_method" "auth" {
  rest_api_id   = aws_api_gateway_rest_api.google_auth.id
  resource_id   = aws_api_gateway_resource.auth.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "auth" {
  rest_api_id             = aws_api_gateway_rest_api.google_auth.id
  resource_id             = aws_api_gateway_resource.auth.id
  http_method             = aws_api_gateway_method.auth.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.google_auth.invoke_arn
}

# /callback endpoint
resource "aws_api_gateway_resource" "callback" {
  rest_api_id = aws_api_gateway_rest_api.google_auth.id
  parent_id   = aws_api_gateway_rest_api.google_auth.root_resource_id
  path_part   = "callback"
}

resource "aws_api_gateway_method" "callback" {
  rest_api_id   = aws_api_gateway_rest_api.google_auth.id
  resource_id   = aws_api_gateway_resource.callback.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "callback" {
  rest_api_id             = aws_api_gateway_rest_api.google_auth.id
  resource_id             = aws_api_gateway_resource.callback.id
  http_method             = aws_api_gateway_method.callback.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.google_auth.invoke_arn
}

# API Gateway deployment
resource "aws_api_gateway_deployment" "google_auth" {
  rest_api_id = aws_api_gateway_rest_api.google_auth.id

  depends_on = [
    aws_api_gateway_integration.auth,
    aws_api_gateway_integration.callback
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "prod" {
  deployment_id = aws_api_gateway_deployment.google_auth.id
  rest_api_id   = aws_api_gateway_rest_api.google_auth.id
  stage_name    = "prod"
}

# Lambda permission for API Gateway
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.google_auth.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.google_auth.execution_arn}/*/*"
}

# Custom domain configuration
resource "aws_acm_certificate" "auth" {
  provider = aws.eu-central-1
  domain_name = "auth.cherevan.art"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_domain_name" "auth" {
  domain_name = "auth.cherevan.art"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.auth.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  depends_on = [aws_acm_certificate.auth]
}

resource "aws_apigatewayv2_api_mapping" "auth" {
  api_id      = aws_api_gateway_rest_api.google_auth.id
  domain_name = aws_apigatewayv2_domain_name.auth.id
  stage       = aws_api_gateway_stage.prod.stage_name
}

output "certificate_validation_records" {
  value = {
    for dvo in aws_acm_certificate.auth.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }
}

output "domain_name_target" {
  value = try(aws_apigatewayv2_domain_name.auth.domain_name_configuration[0].target_domain_name, "")
}

output "api_gateway_url" {
  value = aws_api_gateway_stage.prod.invoke_url
}
