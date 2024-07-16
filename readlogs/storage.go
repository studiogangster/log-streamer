package readlogs

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
)

func Read(prefix string, offset int64, length int64, latestFirst bool) ([]byte, error) {
	// Load environment variables from .env file (optional)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	minioURL := os.Getenv("MINIOURL")
	accessKeyID := os.Getenv("MINIO_ACCESSKEY")
	secretAccessKey := os.Getenv("MINIO_SECRET")
	isSecureStr := os.Getenv("MINIO_SECURE")
	minioRegion := os.Getenv("MINIO_REGION")
	minioBucket := os.Getenv("MINIO_BUCKET") // MinIO bucket name from environment

	// Parse MINIO_SECURE environment variable to boolean
	isSecure, err := strconv.ParseBool(isSecureStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing MINIO_SECURE: %v", err)
	}

	// Initialize MinIO client
	minioClient, err := minio.New(minioURL, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: isSecure,
		Region: minioRegion,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing MinIO client: %v", err)
	}

	// List objects in the bucket with the given prefix
	ctx := context.Background()
	objectCh := minioClient.ListObjects(ctx, minioBucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	})

	// Collect objects
	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %v", object.Err)
		}
		objects = append(objects, object)
	}

	// Sort objects by sequence number extracted from file names
	sort.Slice(objects, func(i, j int) bool {
		seqNoI := extractSeqNo(objects[i].Key, prefix)
		seqNoJ := extractSeqNo(objects[j].Key, prefix)

		if latestFirst {
			return seqNoI > seqNoJ
		}
		return seqNoI < seqNoJ

	})

	// Prepare to concatenate content
	var concatenatedContent []byte
	var currentOffset int64

	// Iterate over sorted objects to find and concatenate content within offset and length
	for _, obj := range objects {
		extractSeqNo(obj.Key, prefix)
		objectSize := obj.Size

		// Check if current object contributes to the offset and length range
		if currentOffset < offset+length && currentOffset+objectSize > offset {
			// Calculate start and end positions within the current object
			start := int64(0)
			end := objectSize

			if currentOffset < offset {
				start = offset - currentOffset
			}
			if currentOffset+objectSize > offset+length {
				end = offset + length - currentOffset
			}

			getObjOptions := minio.GetObjectOptions{}
			getObjOptions.SetRange(start, end-start)
			// Retrieve object content within the specified range
			objectContent, err := minioClient.GetObject(ctx, minioBucket, obj.Key, getObjOptions)
			if err != nil {
				return nil, fmt.Errorf("error retrieving object %s: %v", obj.Key, err)
			}

			// Read and append object content to the concatenated content
			objectBytes, err := io.ReadAll(objectContent)
			objectContent.Close()
			if err != nil {
				return nil, fmt.Errorf("error reading object %s content: %v", obj.Key, err)
			}

			concatenatedContent = append(concatenatedContent, objectBytes...)
		}

		// Update current offset
		currentOffset += objectSize

		// Stop processing if we've collected enough content
		if currentOffset >= offset+length {
			break
		}
	}

	return concatenatedContent, nil
}

// extractSeqNo extracts the sequential number from a file name.
func extractSeqNo(fileName string, prefix string) int {
	// Assuming file name format is prefix<number>.[extension]
	fileName = strings.Split(fileName, prefix)[1]
	parts := strings.Split(fileName, ".")

	if len(parts) < 2 {
		return 0 // Invalid file name format
	}

	seqStr := strings.TrimPrefix(parts[0], prefix) // Adjust prefix as needed

	seqNo, err := strconv.Atoi(seqStr)
	if err != nil {
		return 0 // Error converting sequential number
	}

	return seqNo
}
