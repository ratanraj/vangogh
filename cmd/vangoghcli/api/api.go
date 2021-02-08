package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ratanraj/vangogh/database"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type ClientConfig struct {
	Token string `json:"token"`
}

type API struct {
	token string
	Host  string
}

func NewAPI(host string) *API {
	fp, err := os.Open("/home/ratanraj/.vangogh")
	if err != nil {
		panic(err)
	}
	d := json.NewDecoder(fp)
	var conf ClientConfig
	err = d.Decode(&conf)
	if err != nil {
		panic(err)
	}

	return &API{token: conf.Token, Host: host}
}

func (a API) Request(method, urlPath string, body io.Reader, headers *map[string]string) (*http.Response, error) {
	client := &http.Client{}
	u, err := url.Parse(a.Host)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, urlPath)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))
	if headers != nil {
		for key, val := range *headers {
			req.Header.Add(key, val)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return res, err
	}
	if (res.StatusCode / 100) != 2 {
		return res, fmt.Errorf(res.Status)
	}
	return res, err
}

func (a API) get(url string) error {
	_, err := a.Request("GET", url, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (a API) DoLogin(username, password string) error {
	client := http.Client{}

	b, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	})

	requestBody := strings.NewReader(string(b))
	resp, err := client.Post("http://127.0.0.1:8080/auth/login", "application/json", requestBody)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Invalid login"))
	}

	fp, err := os.Create("/home/ratanraj/.vangogh")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(fp, resp.Body)
	if err != nil {
		panic(err)
	}
	return nil
}

func (a API) ListAlbums() error {
	res, err := a.Request("GET", "/api/album", nil, nil)
	if err != nil {
		return err
	}

	var jsonResponse struct {
		Albums []database.Album `json:"albums"`
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&jsonResponse)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title"})
	for _, album := range jsonResponse.Albums {
		t.AppendRows([]table.Row{{album.ID, album.Title}})
		t.AppendSeparator()
		//fmt.Printf("%d\t%s\n",album.ID,album.Title)
	}
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	return nil
}

func (a API) CreateAlbum(albumName string) error {
	b, err := json.Marshal(struct {
		AlbumTitle string `json:"album_title"`
	}{
		AlbumTitle: albumName,
	})
	if err != nil {
		return err
	}

	requestBody := strings.NewReader(string(b))

	res, err := a.Request(http.MethodPut, "/api/album", requestBody, nil)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (a API) DeleteAlbum(albumID uint) error {
	_, err := a.Request(http.MethodDelete, fmt.Sprintf("/api/album/%d", albumID), nil, nil)
	return err
}

func (a API) UploadPhoto(albumID uint, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	_, err = a.Request(http.MethodPut,
		fmt.Sprintf("/api/album/%d/photo", albumID),
		body,
		&map[string]string{
			"Content-Type": writer.FormDataContentType(),
		})

	if err != nil {
		return err
	}

	return nil
}

func (a API) ListPhotos(albumID uint) error {
	res, err := a.Request("GET", fmt.Sprintf("/api/album/%d/photo", albumID), nil, nil)
	if err != nil {
		return err
	}

	var photoResponse struct {
		Photos []database.Photo `json:"photos"`
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&photoResponse)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "FileName", "Album"})
	for i := range photoResponse.Photos {
		_photo := photoResponse.Photos[i]
		t.AppendRows([]table.Row{{_photo.ID, _photo.FileName, _photo.Album.Title}})
		t.AppendSeparator()
	}
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	return nil
}
