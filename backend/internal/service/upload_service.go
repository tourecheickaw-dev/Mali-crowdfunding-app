package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const maxUploadBytes = 5 * 1024 * 1024 // 5 MB

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

var (
	safeOwnerID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	fileNameRe  = regexp.MustCompile(`^[a-f0-9]{32}\.(jpg|png|webp)$`)
)

type PresignedUpload struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
	PublicURL string `json:"public_url"`
	ExpiresIn int    `json:"expires_in"`
}

type UploadService struct {
	presigner *s3.PresignClient
	bucket    string
	region    string
}

func NewUploadService(ctx context.Context) (*UploadService, error) {
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return nil, errors.New("S3_BUCKET is not set")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("aws config: %w", err)
	}
	return &UploadService{
		presigner: s3.NewPresignClient(s3.NewFromConfig(cfg)),
		bucket:    bucket,
		region:    region,
	}, nil
}

func (s *UploadService) ownerPrefix(ownerID string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/uploads/%s/", s.bucket, s.region, ownerID)
}

// IsOwnedURL reports whether url is a photo uploaded by this owner in our bucket.
func (s *UploadService) IsOwnedURL(ownerID, url string) bool {
	if !safeOwnerID.MatchString(ownerID) {
		return false
	}
	prefix := s.ownerPrefix(ownerID)
	if !strings.HasPrefix(url, prefix) {
		return false
	}
	return fileNameRe.MatchString(strings.TrimPrefix(url, prefix))
}

// PresignImage returns a temporary URL the browser can use to upload one image.
// Validation errors are plain errors (400); AWS errors are wrapped with %w (500).
func (s *UploadService) PresignImage(ctx context.Context, ownerID, contentType string, size int64) (*PresignedUpload, error) {
	if !safeOwnerID.MatchString(ownerID) {
		return nil, errors.New("invalid identifier")
	}
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return nil, errors.New("unsupported format: use JPEG, PNG or WebP")
	}
	if size <= 0 || size > maxUploadBytes {
		return nil, errors.New("the photo must be smaller than 5 MB")
	}

	rnd := make([]byte, 16)
	if _, err := rand.Read(rnd); err != nil {
		return nil, fmt.Errorf("random: %w", err)
	}
	name := hex.EncodeToString(rnd) + ext
	key := fmt.Sprintf("uploads/%s/%s", ownerID, name)

	req, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("presign: %w", err)
	}

	return &PresignedUpload{
		UploadURL: req.URL,
		Key:       key,
		PublicURL: s.ownerPrefix(ownerID) + name,
		ExpiresIn: 300,
	}, nil
}
