package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/ratanraj/vangogh/database"
	"github.com/ratanraj/vangogh/storage"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func ListAlbums(c *gin.Context) {
	if uid, ok := c.Get("user"); ok {
		var albums []database.Album
		database.DBConn.Where("owner_refer = ?", uid.(uint)).Preload("Owner").Find(&albums)
		c.JSON(http.StatusOK, gin.H{"albums": albums})
	} else {
		log.Println("no uid")
	}
}

func CreateAlbum(c *gin.Context) {
	albumData := struct {
		AlbumTitle string `json:"album_title"`
	}{}
	err := c.ShouldBindJSON(&albumData)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if uid, ok := c.Get("user"); ok {
		album := database.Album{
			Title:      albumData.AlbumTitle,
			OwnerRefer: uid.(uint),
		}
		database.DBConn.Create(&album)
	}
}

func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	database.DBConn.Delete(&database.Album{}, id)
}

// PHOTOS

func ListPhotos(c *gin.Context) {
	if uid, ok := c.Get("user"); ok {
		albumID := c.Param("album_id")
		var photos []database.Photo
		_=uid
		database.DBConn.
			//Where("owner_refer = ?", uid.(uint)).
			Where("album_refer = ?", albumID).
			Preload("Owner").
			Preload("Album").
			Find(&photos)
		c.JSON(http.StatusOK, gin.H{"photos":photos})
	}
}

func UploadPhoto(c *gin.Context) {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		panic(err)
	}
	fp, err := fileHeader.Open()
	if err != nil {
		panic(err)
	}

	photo := database.Photo{
		Title:         filepath.Base(fileHeader.Filename),
		FileName:      fileHeader.Filename,
		StorageKey:    "",
		StorageBucket: "",
		ThumbnailURL:  "",
		Size:          fileHeader.Size,
		AlbumRefer:    0,
		OwnerRefer:    0,
	}
	if uid, ok := c.Get("user"); ok {
		photo.OwnerRefer = uid.(uint)
	} else {
		panic(err)
	}
	albumID := c.Param("album_id")
	albumIDUINT,err := strconv.ParseInt(albumID, 10, 64)
	if err != nil {
		panic(err)
	}
	photo.AlbumRefer = uint(albumIDUINT)

	buffer := make([]byte, fileHeader.Size)
	_, err = fp.Read(buffer)
	if err != nil {
		panic(err)
	}


	x, err := exif.Decode(bytes.NewReader(buffer))
	if err == nil {
		b, err := x.MarshalJSON()
		if err == nil {
			photo.ExifData = string(b)
		}
	}

	key,err := storage.UploadPhoto(fileHeader.Filename, fileHeader.Size, buffer)
	if err != nil {
		panic(err)
	}
	photo.StorageKey = key

	database.DBConn.Create(&photo)
}

func GetPhoto(c *gin.Context) {

}

func DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	database.DBConn.Delete(&database.Photo{}, id)
}
