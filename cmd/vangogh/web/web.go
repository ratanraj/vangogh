package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const DEBUG = true

func RunServer(addr string) {
	r := gin.New()
	r.Use(gin.Logger(),
		gin.Recovery())

	AccessSecret := os.Getenv("ACCESS_SECRET")
	store := cookie.NewStore([]byte(AccessSecret))
	r.Use(sessions.Sessions("vangoghsession", store))

	r.LoadHTMLGlob("web/templates/*")

	r.GET("/", Home)
	r.GET("/new-user", NewUserPage)
	r.POST("/new-user", NewUserPage)

	r.StaticFile("/favicon.ico", "./web/static/favicon.ico")
	r.Static("/static", "./web/static/")

	{ // OAuth Routes
		r.GET("/login", LoginController)
		r.GET("/callback", CallbackController)
	}
	api := r.Group("/api", UserMiddleware())
	{
		api.GET("/me", Me)

		api.GET("/album", ListAlbums)
		api.PUT("/album", CreateAlbum)
		api.DELETE("/album/:id", DeleteAlbum)

		api.GET("/album/:album_id/photo", ListPhotos)
		api.PUT("/album/:album_id/photo", UploadPhoto)

		api.GET("/photo/:id", GetPhoto)
		api.DELETE("/photo/:id", DeletePhoto)
	}

	log.Fatal(r.Run(addr))
}