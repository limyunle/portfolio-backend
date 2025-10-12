package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Storage interface {
	UploadJSON(ctx context.Context, bucket, key string, data interface{}) error
	GetJSON(ctx context.Context, bucket, key string, out interface{}) error
}

type Service struct {
	Client S3API
	Bucket string
}

func NewService(client S3API, bucket string) *Service {
	return &Service{
		Client: client,
		Bucket: bucket,
	}
}

func (s *Service) UploadJSON(ctx context.Context, bucket, key string, data interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return fmt.Errorf("json encode failed: %w", err)
	}

	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   buf,
	})
	if err != nil {
		return fmt.Errorf("s3 upload failed: %w", err)
	}

	return nil
}

func (s *Service) GetJSON(ctx context.Context, bucket, key string, out interface{}) error {
	obj, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 get object failed: %w", err)
	}
	defer obj.Body.Close()

	if err := json.NewDecoder(obj.Body).Decode(out); err != nil {
		return fmt.Errorf("json decode failed: %w", err)
	}
	return nil
}
