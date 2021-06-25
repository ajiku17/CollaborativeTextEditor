package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const appDir = "./web/app"

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Usage: server [port]")
		return
	}

	port := ":" + args[1]

	appFs := http.FileServer(http.Dir(appDir))

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		appFs.ServeHTTP(resp, req)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}