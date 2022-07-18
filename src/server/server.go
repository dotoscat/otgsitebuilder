package server

import (
    "net/http"
    // "log"
    "embed"
    "mime"
    "github.com/julienschmidt/httprouter"
    "fmt"

    // "github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed app/*
var appFS embed.FS
var appFileSystem http.FileSystem

func NewWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "New website!")
}

func LoadWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Load website!")
}

func SaveWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Save website!")
}

func BuildWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Build website!")
}

func PathContent(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
    path := ps.ByName("path")
    if path == "home" {
        fmt.Fprintf(w, "Default directory")
    } else {
        fmt.Fprintf(w, "Something else")
    }
}

func Start(addr string) error {
    mime.AddExtensionType(".js", "text/javascript")
    appFileSystem = http.FS(appFS)

    //http.HandleFunc("/",)
    // http.Handle("/", http.FileServer(appFileSystem))

    router := httprouter.New()

    // Specify the path to the website content
    router.POST("/website", NewWebsite)
    router.GET("/website/:path", LoadWebsite) //add to url the path
    router.PUT("/website", SaveWebsite)
    router.POST("/website/build", BuildWebsite)

    router.GET("/path/:path", PathContent)

    err := http.ListenAndServe(addr, router)

    return err
}
