// Package aws helps developers to use AWS cloud API more friendly
package aws

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	app "github.com/pottava/golang-microservices/app-aws/app/config"
)

func config() *aws.Config {
	log := aws.LogLevel(aws.LogOff)
	cfg := app.NewConfig()
	if cfg.AwsLog {
		log = aws.LogLevel(aws.LogDebug)
	}
	return &aws.Config{
		Credentials: credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&ec2rolecreds.EC2RoleProvider{ExpiryWindow: cfg.AwsRoleExpiry * time.Minute},
			}),
		Region:   aws.String(os.Getenv("AWS_REGION")),
		LogLevel: log,
	}
}
