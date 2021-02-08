package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"net/http"
	"strings"
)

const BucketName = "vangogh"

var (
	BucketStorage *minio.Client
)

func UploadPhoto(fileName string, size int64, buffer []byte) (string, error) {

	contentType := http.DetectContentType(buffer)

	fileNameParts := strings.Split(fileName, ".")
	extension := fileNameParts[len(fileNameParts)-1]
	key := fmt.Sprintf("%s.%s", uuid.New().String(), extension)

	uploadInfo, err := BucketStorage.PutObject(context.TODO(),
		BucketName,
		key,
		bytes.NewBuffer(buffer),
		size,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{"FileName": fileName},
			UserTags:     nil,
			ContentType:  contentType,
			Internal:     minio.AdvancedPutOptions{},
		})
	if err != nil {
		return "", err
	}

	return uploadInfo.Key, nil
}
