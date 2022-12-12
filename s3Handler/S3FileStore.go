package s3Handler

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
	"log"
	"os"
	"sync"
)

type S3FileStore struct {
	AccessKeyId     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Session         *session.Session
	mutex           sync.RWMutex
}

var singletonS3FileStore *S3FileStore
var onceS3FileStore sync.Once

func GetS3FileStore() *S3FileStore {
	onceS3FileStore.Do(func() {
		log.Println("Init S3FileStore...")
		accessKeyID, secretAccessKey, region, bucket, session := ConnectAws()
		s3FileStore := S3FileStore{
			AccessKeyId:     accessKeyID,
			SecretAccessKey: secretAccessKey,
			Region:          region,
			Bucket:          bucket,
			Session:         session,
		}

		singletonS3FileStore = &s3FileStore
	})
	return singletonS3FileStore
}

func ConnectAws() (string, string, string, string, *session.Session) {
	accessKeyID := os.Getenv("S3_ACK")
	secretAccessKey := os.Getenv("S3_SCK")
	region := os.Getenv("S3_REGION")
	bucket := os.Getenv("S3_BUCKET")
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		log.Fatalf("Failed to connect to AWS: %v", err)
	}
	return accessKeyID, secretAccessKey, region, bucket, session
}

func (s3FileStore *S3FileStore) UploadFile(fileName string, fileData bytes.Buffer, mimetype *mimetype.MIME) (*s3manager.UploadOutput, error) {
	log.Println("Upload File to S3...")
	//s3FileStore.mutex.Lock()
	//defer s3FileStore.mutex.Unlock()

	uploader := s3manager.NewUploader(s3FileStore.Session)
	//upload to the s3 bucket
	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3FileStore.Bucket),
		//ACL:    aws.String("public-read"),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(fileData.Bytes()),
		ContentType: aws.String(mimetype.String()),
	})
}

func (s3FileStore *S3FileStore) GetObject(fileName string) (*s3.GetObjectOutput, error) {
	//log.Println("GetObject from S3...")
	downloader := s3manager.NewDownloader(s3FileStore.Session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3FileStore.Bucket),
		Key:    aws.String(fileName),
	}

	return downloader.S3.GetObject(input)
}
