package main

import (
	"fmt"
	"net/url"
)
import "net/http"

func main() {
	resp, err := http.PostForm("http://portalmeteo.pl/index/login",
		url.Values{"username": {"usr"}, "password": {"pass"}})
	fmt.Println(resp)
	fmt.Println(err)
}
