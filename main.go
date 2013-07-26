// Email project main.go
package main

import (
	"github.com/coocood/jas"
	"log"
	"net/http"
	"runtime"
)

var MaxProcess = runtime.NumCPU() * 2

func main() {

	runtime.GOMAXPROCS(MaxProcess)
	defer DBSession.Close()

	init_DB()

	router := jas.NewRouter(new(Newsletter))
	router.BasePath = "/v1/"
	router.DisableAutoUnmarshal = true

	log.Println(router.HandledPaths(true))

	http.Handle(router.BasePath, router)
	http.ListenAndServe(":8080", nil)

}
