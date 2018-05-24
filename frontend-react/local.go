package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	var proxyUrl *url.URL
	if len(os.Args) == 2 && os.Args[1] == "useLocal" {
		log.Println("using local backend")
		proxyUrl, _ = url.Parse("http://127.0.0.1:5000/")
	} else {
		proxyUrl, _ = url.Parse("http://localhost:18000/")
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	http.Handle("/api/", proxy)

	localProxyUrl, _ := url.Parse("http://localhost:3000/")
	localProxy := httputil.NewSingleHostReverseProxy(localProxyUrl)
	http.Handle("/", localProxy)

	log.Println("Serving on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
