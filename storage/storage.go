package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
)

type Storage struct {
	svc *s3.S3
}

func NewStorage(region, endpoint string) *Storage {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess, &aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})

	return &Storage{svc}
}

func (s *Storage) CreateBucket(bucket string) error {

	// Try to check bucket exist
	_, err := s.svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		return fmt.Errorf("the bucket %q already exist", bucket)
	}

	fmt.Printf("Successfully read bucket %q\n", bucket)

	// Create the S3 Bucket
	_, err = s.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("unable to create bucket %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = s.svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("error occurred while waiting for bucket to be created, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)
	return nil
}

func (s *Storage) Upload(bucket, filename string, file multipart.File) error {
	sess, err := session.NewSession(&s.svc.Config)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		// Print the error and exit.
		return fmt.Errorf("unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
	return nil
}

func (s *Storage) PutObjectTag(bucket, object string, name string) error {
	params := &s3.PutObjectTaggingInput{
		Bucket: &bucket,
		Key:    &object,
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{
				{
					Key:   aws.String("name"),
					Value: aws.String(name),
				},
			},
		},
	}

	// Set bucket ACL
	_, err := s.svc.PutObjectTagging(params)
	if err != nil {
		return err
	}

	fmt.Println("Successfully put tags on object")
	return nil
}

func (s *Storage) GetObjectTag(bucket, object string) (*s3.GetObjectTaggingOutput, error) {
	// Get object tagging
	return s.svc.GetObjectTagging(&s3.GetObjectTaggingInput{Bucket: &bucket, Key: &object})
}

func (s *Storage) GetListObject(bucket string) ([]*s3.Object, error) {
	// Get the list of items
	resp, err := s.svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return nil,
			fmt.Errorf("unable to list items in bucket %q, %v", bucket, err)
	}

	fmt.Printf("found %d items in bucket %q\n", len(resp.Contents), bucket)

	return resp.Contents, nil

}

func (s *Storage) Download(bucket, item string) (*aws.WriteAtBuffer, error) {
	sess, err := session.NewSession(&s.svc.Config)
	if err != nil {
		return nil, err
	}
	downloader := s3manager.NewDownloader(sess)

	buffer := aws.NewWriteAtBuffer([]byte{})
	numBytes, err := downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		return nil, fmt.Errorf("unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded ", numBytes, " bytes")

	return buffer, nil
}
