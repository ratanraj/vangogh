package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ratanraj/vangogh/database"
	"log"
	"net/http"
)

func Me(c *gin.Context) {
	if email, ok := c.Get("email"); ok {
		var user database.User
		database.DBConn.Where("email = ?", email.(string)).Find(&user)
		c.JSON(http.StatusOK, gin.H{
			"id": user.ID,
			"email": user.Email,
			"first_name": user.FirstName,
			"last_name": user.LastName,
		})
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func Home(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email")

	emailString := ""
	if email != nil {
		emailString = email.(string)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"URL":"http://127.0.0.1:3000","Email":emailString, "DEBUG": DEBUG})
}

func NewUserPage(c *gin.Context) {
	type userForm struct {
		FirstName string `json:"first_name" form:"first_name" binding:"required"`
		LastName string `json:"last_name" form:"last_name" binding:"required"`
	}
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "new_user.html", gin.H{})
	} else if c.Request.Method == http.MethodPost {
		var newUser userForm
		err := c.ShouldBind(&newUser)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Println(newUser.FirstName)
		log.Println(newUser.LastName)

		session := sessions.Default(c)
		email := session.Get("email")

		emailString := ""
		if email != nil {
			emailString = email.(string)
		}

		user := database.User{
			Email:     emailString,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Active:    true,
		}

		database.DBConn.Create(&user)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
