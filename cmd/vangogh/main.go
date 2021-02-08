package main

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ratanraj/vangogh/cmd/vangogh/web"
	"github.com/ratanraj/vangogh/database"
	"github.com/ratanraj/vangogh/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func initDB(url string) {
	var err error

	MinioId := os.Getenv("MINIO_ID")
	MinioSecret := os.Getenv("MINIO_SECRET")
	MinioEndpoint := os.Getenv("MINIO_ENDPOINT")

	storage.BucketStorage, err = minio.New(MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioId, MinioSecret, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	database.DBConn, err = gorm.Open(sqlite.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = database.DBConn.AutoMigrate(
		&database.User{},
		&database.Album{},
		database.Photo{},
		&database.Session{},
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB("sqlite3.db")

	web.RunServer()
}
