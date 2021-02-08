package database

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Title         string `gorm:"title" json:"title"`
	FileName      string `gorm:"file_name" json:"file_name"`
	StorageKey    string `gorm:"storage_key" json:"storage_key"`
	StorageBucket string `gorm:"storage_bucket" json:"storage_bucket"`
	ThumbnailURL  string `gorm:"thumbnail_url" json:"thumbnail_url"`
	Size          int64  `gorm:"size" json:"size"`
	ExifData      string `gorm:"exif_data" json:"exif_data"`
	AlbumRefer    uint   `gorm:"album_refer" json:"album_refer"`
	Album         Album  `gorm:"foreignKey:AlbumRefer" json:"album"`
	OwnerRefer    uint   `json:"owner_refer"`
	Owner         User   `gorm:"foreignKey:OwnerRefer" json:"owner"`
}
