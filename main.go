package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tkodyl/vineguard/configuration"
)

type FilePath struct {
	Filename string
}

func main() {
	client := http.Client{}
	phpSessionCookie := http.Cookie{Name: "PHPSESSID", Value: "msl9tkd9p7cifcsnumsnit9bn3"}

	config := configuration.GetConfig()

	statusCode, err := loginRequest(&client, phpSessionCookie, config)
	if err != nil {
		fmt.Println("Error during login, message:", err.Error())
		return
	}
	fmt.Println(statusCode)

	filePath, err := dataCreationRequest(&client, phpSessionCookie, config.Server.Url)
	if err != nil {
		fmt.Println("Error during data creation, message:", err.Error())
		return
	}
	fileContent, err := retrieveData(&client, phpSessionCookie, filePath, config.Server.Url)
	fmt.Println(fileContent)
}

func loginRequest(client *http.Client, cookie http.Cookie, config configuration.Config) (string, error) {
	form := url.Values{}
	form.Add("username", config.Server.Credentails.Username)
	form.Add("password", config.Server.Credentails.Password)

	req, err := http.NewRequest("POST", config.Server.Url+"/index/login", strings.NewReader(form.Encode()))
	req.AddCookie(&cookie)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "www.portalmeteo.pl")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip,deflate,br")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

func dataCreationRequest(client *http.Client, cookie http.Cookie, serverAddress string) (string, error) {
	currentTime := time.Now()
	todaysDate := currentTime.Format("2006-01-02")
	sevenDaysBefore := currentTime.AddDate(0, 0, -7).Format("2006-01-02")

	form := url.Values{}
	form.Add("type", "rimpro")
	form.Add("station_id", "483")
	form.Add("limit", "7")
	form.Add("page", "1")
	form.Add("date_start", sevenDaysBefore)
	form.Add("date_end", todaysDate)

	req, err := http.NewRequest("POST", serverAddress+"/ajax/new-export", strings.NewReader(form.Encode()))

	if err != nil {
		return "", errors.New("Request is not valid at all")
	}

	req.AddCookie(&cookie)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp2, _ := client.Do(req)
	var filePath FilePath
	decodeErr := json.NewDecoder(resp2.Body).Decode(&filePath)
	if decodeErr != nil {
		return "", errors.New("Error on unmarshalling file path")
	}
	return filePath.Filename, nil
}

func retrieveData(client *http.Client, cookie http.Cookie, filePath string, serverAddress string) (string, error) {
	req, err := http.NewRequest("GET", serverAddress+filePath, nil)
	if err != nil {
		return "", errors.New("Request is not valid at all")
	}
	req.AddCookie(&cookie)
	resp, _ := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString, nil
}
