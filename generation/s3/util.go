package s3

import (
	"runtime"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	mutex      = sync.Mutex{}
	awsSession *session.Session
)

// NewS3Client returns a new S3 client using the aws session
func NewS3Client() (*s3.S3, error) {
	mutex.Lock()
	if awsSession == nil {
		// You will need aws credentials (access key id, secret access key) and region information
		// AWS looks for these in the following areas
		// 1) Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION)
		// 2) Credentials file (Shared/Profile Specific)
		// 3) IAM roles if running on EC2
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
