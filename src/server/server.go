package server

import (
    "net/http"
    "net/url"
    "os"
    "encoding/json"
    "log"
    "embed"
    "mime"
    "fmt"
    "path/filepath"

    "github.com/julienschmidt/httprouter"
    // "github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed app/*
var appFS embed.FS
var appFileSystem http.FileSystem

func NewWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "New website!")
}

func LoadWebsite(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
    path := ps.ByName("path")
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    fmt.Fprint(w, "Load website!", path)
}

func SaveWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Save website!")
}

func BuildWebsite(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Build website!")
}

type DirEntry struct {
    PathUrl string
    Name string
    Ftype string
}

type DirList struct {
    Parent string
    List []DirEntry
}

func listPath(path string) (dirList DirList) {
    // parent
    // list: [{url, name, type}]

    if paths, err := os.ReadDir(path); err != nil {
        return
    } else {
        dirList.Parent = url.PathEscape(filepath.Dir(path))
        for _, file := range paths {
            dirEntry := DirEntry{}

            dirEntry.PathUrl = url.PathEscape(filepath.Join(path, file.Name()))

            dirEntry.Name = file.Name()

            if file.Type().IsDir() == true {
                dirEntry.Ftype = "d"
            } else if file.Type().IsRegular() == true {
                dirEntry.Ftype = "f"
            }

            dirList.List = append(dirList.List, dirEntry)
        }
    }
    return
}

func outputPath(w http.ResponseWriter, path string) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    dirList := listPath(path)
    if output, err := json.Marshal(dirList); err != nil {
        fmt.Fprintln(w, err)
    } else {
        fmt.Fprintln(w, string(output))
    }
}

func HomeContent(w http.ResponseWriter) {
    if homeDir, err := os.UserHomeDir(); err != nil {
        fmt.Println(err)
    } else {
        outputPath(w, homeDir)
    }
}

func PathHandler(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
    requestedPath := ps.ByName("path")
    if requestedPath == "home" {
        HomeContent(w)
    } else {
        outputPath(w, requestedPath)
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

    router.GET("/path/:path", PathHandler)

    log.Println("Ready to go")
    err := http.ListenAndServe(addr, router)

    return err
}
