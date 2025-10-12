package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	putCalled bool
	getCalled bool
}

func (m *mockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.putCalled = true

	if input.Body == nil || input.Bucket == nil || input.Key == nil {
		return nil, fmt.Errorf("invalid input")
	}

	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	m.getCalled = true

	data := map[string]string{"hello": "world"}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	return &s3.GetObjectOutput{
		Body: io.NopCloser(buf),
	}, nil
}

func TestUploadJSON(t *testing.T) {
	mock := &mockS3Client{}
	service := NewService(mock, "test-bucket")

	data := map[string]string{"foo": "bar"}
	err := service.UploadJSON(context.Background(), "test-bucket", "key", data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.putCalled {
		t.Fatalf("expected PutObject to be called")
	}
}

func TestGetJSON(t *testing.T) {
	mock := &mockS3Client{}
	service := NewService(mock, "test-bucket")

	var out map[string]string
	err := service.GetJSON(context.Background(), "test-bucket", "key", &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.getCalled {
		t.Fatalf("expected GetObject to be called")
	}

	if out["hello"] != "world" {
		t.Errorf("expected hello=world, got %+v", out)
	}
}
