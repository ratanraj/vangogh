package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func RunServer() {
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

		api.GET("/ping", func(c *gin.Context) {
			if uid,ok := c.Get("user"); ok {
				c.JSON(http.StatusOK, gin.H{"message": "Pong!", "user":uid})
			}
		})

		api.GET("/album", ListAlbums)
		api.PUT("/album", CreateAlbum)
		api.DELETE("/album/:id", DeleteAlbum)

		api.GET("/album/:album_id/photo", ListPhotos)
		api.PUT("/album/:album_id/photo", UploadPhoto)

		api.GET("/photo/:id", GetPhoto)
		api.DELETE("/photo/:id", DeletePhoto)
	}

	log.Fatal(r.Run())
}