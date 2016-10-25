package awsSession

import (
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/unchartedsoftware/plog"
)

var awsSession *session.Session = nil

func Get() (*session.Session) {
  return awsSession
}

func Create() (*session.Session) {
  // AWS looks for credentials in the following places:
  // 1) Environment variables (AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY_ID)
  // 2) Credentials file (Shared/Profile Specific)
  // 3) IAM roles if running on EC2
	sess, err := session.NewSession()
	if err != nil {
    awsSession = nil
		log.Error(err)
	} else {
    awsSession = sess
  }

  return awsSession
}
