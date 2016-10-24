package awsSession

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/unchartedsoftware/plog"
)

var awsSession *session.Session = nil

func Get() (*session.Session) {
  return awsSession
}

func Create() (*session.Session) {
	sess, err := session.NewSession(&aws.Config{Credentials: credentials.NewEnvCredentials()})
	if err != nil {
    awsSession = nil
		log.Error(err)
	} else {
    awsSession = sess
  }

  return awsSession
}
