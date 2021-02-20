package web

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ratanraj/vangogh/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}

	randomState = "random"
)

type CallbackData struct {
	State    string `form:"state"`
	Code     string `form:"code"`
	Scope    string `form:"scope"`
	AuthUser uint   `form:"authuser"`
	Prompt   string `form:"prompt"`
}

type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func LoginController(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func CallbackController(c *gin.Context) {
	var callbackData CallbackData
	err := c.ShouldBind(&callbackData)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	log.Println("******************")
	log.Println(callbackData)
	log.Println("******************")
	if callbackData.State != randomState {
		log.Println("state is not valid")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, callbackData.Code)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	var userInfo OAuthUserInfo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&userInfo)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.Set("email", userInfo.Email)

	var user database.User
	database.DBConn.Where("email = ?", userInfo.Email).Find(&user)

	session := sessions.Default(c)
	session.Set("email", userInfo.Email)
	session.Set("picture", userInfo.Picture)
	session.Set("uid", user.ID)
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
	//c.JSON(http.StatusOK, gin.H{"Response": userInfo})
}

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("email")
		u := session.Get("uid")

		if uid, ok := u.(uint); ok {
			c.Set("uid", uid)
		}

		if emailString, ok := email.(string); ok {
			c.Set("email", emailString)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}