package aws

import (
	"runtime"
	"sync"

  "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	mutex   = sync.Mutex{}
	awsSession *session.Session = nil
)

func NewS3Client() (*s3.S3, error) {
  mutex.Lock()
  if awsSession == nil {
    sess, err := session.NewSession()
    if err != nil {
      mutex.Unlock()
      runtime.Gosched()
      return nil, err
    }
    awsSession = sess
  }
  mutex.Unlock()
  runtime.Gosched()
  return s3.New(awsSession), nil
}
