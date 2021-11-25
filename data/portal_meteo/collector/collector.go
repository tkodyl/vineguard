package collector

import (
	"encoding/json"
	"errors"
	"github.com/tkodyl/vineguard/configuration"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type FilePath struct {
	Filename string
}

type Collector struct {
	client        http.Client
	sessionCookie http.Cookie
	config        *configuration.Config
}

func NewCollector(config *configuration.Config) Collector {
	client := http.Client{}
	sessionCookie := http.Cookie{Name: "PHPSESSID", Value: "msl9tkd9p7cifcsnumsnit9bnd"}
	return Collector{client: client, sessionCookie: sessionCookie, config: config}
}

func (coll *Collector) GetDataFromPortalMeteo() (string, error) {
	log.Println("Attempt to login to", coll.config.Server.Url)
	_, err := coll.loginRequest()
	if err != nil {
		log.Println("Error during login, message:", err.Error())
		return "", err
	}

	log.Println("Attempt to export data to csv file")
	filePath, err := coll.dataExportRequest()
	if err != nil {
		log.Println("Error during data creation, message:", err.Error())
		return "", err
	}
	log.Println("Retrieval data from path:", filePath)
	fileContent, err := coll.retrieveData(filePath)
	return fileContent, err
}

func (coll *Collector) loginRequest() (string, error) {
	form := url.Values{}
	form.Add("username", coll.config.Server.Credentials.Username)
	form.Add("password", coll.config.Server.Credentials.Password)

	req, err := http.NewRequest("POST", coll.config.Server.Url+"/index/login", strings.NewReader(form.Encode()))
	req.AddCookie(&coll.sessionCookie)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "www.portalmeteo.pl")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip,deflate,br")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	resp, err := coll.client.Do(req)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

func (coll *Collector) dataExportRequest() (string, error) {
	form := createFormForWeekDataExport()
	req, err := http.NewRequest("POST", coll.config.Server.Url+"/ajax/new-export", strings.NewReader(form.Encode()))
	if err != nil {
		return "", errors.New("Request is not valid at all")
	}
	req.AddCookie(&coll.sessionCookie)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, _ := coll.client.Do(req)
	if response.StatusCode != 200 {
		return "", errors.New("data export request failed")
	}
	var filePath FilePath
	decodeErr := json.NewDecoder(response.Body).Decode(&filePath)
	if decodeErr != nil {
		return "", errors.New("error on unmarshalling file path")
	}
	return filePath.Filename, nil
}

func (coll *Collector) retrieveData(filePath string) (string, error) {
	req, err := http.NewRequest("GET", coll.config.Server.Url+filePath, nil)
	if err != nil {
		return "", errors.New("Request is not valid at all")
	}
	req.AddCookie(&coll.sessionCookie)
	resp, _ := coll.client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString, nil
}

func createFormForWeekDataExport() url.Values {
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
	return form
}
