package constants

import "cloudview/app/src/types"

var OAuth = types.OauthProviders{
	GOOGLE: "google",
	GITHUB: "github",
	LOCAL:  "local",
}

// The same cookie name is set on the client side.
var COOKIE_NAME = "cv-token"

// Make sure to match the values in webapp
var (
	EC2          = "ec2"
	APIGATEWAYV2 = "apigatewayV2"
	ROUTE53      = "route53"
	LAMBDA       = "lambda"
	EKS          = "eks"
	ELBV2        = "elbV2"
	S3           = "s3"
	EFS          = "efs"
	RDS          = "rds"
	CLOUDFRONT   = "cloudfront"
	CLOUDWATCH   = "cloudwatch"
	DYNAMODB     = "dynamodb"
)
