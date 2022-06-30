package server

import (
    "net/http"
    // "log"
    "embed"
    "mime"
    //"github.com/julienschmidt/httprouter"
)

//go:embed app/*
var appFS embed.FS
var appFileSystem http.FileSystem

//func postHandler(c *gin.Context) {
//    c.String(http.StatusOK, "OK")
//}

func Start(addr string) error {
    mime.AddExtensionType(".js", "text/javascript")
    // router := gin.Default()

    // router.StaticFS("/public", http.FS(appFS))

    //router.GET("/", func(c *gin.Context) {
    //    c.HTML(http.StatusOK, "index.html", nil)
    //})

    appFileSystem = http.FS(appFS)

    // router.Run(addr)
    //router := httprouter.New()
    //router.ServeFiles("/*filepath", appFileSystem)
    //http.ListenAndServe(addr, router)

    //http.HandleFunc("/",)
    http.Handle("/", http.FileServer(appFileSystem))
    http.ListenAndServe(addr, nil)

    return nil
}
