package main

import (
    "log"
    "fmt"
    "flag"
    "os"

    "github.com/dotoscat/otgsitebuilder/src/manager"
    )

const (
    MANAGER_MODE = "manager"
    BUILDER_MODE = "builder"
)

func main() {
    var mode string
    var content string
    var filename string
    flag.StringVar(&mode, "mode", "", "Set the mode of use of this tool")
    flag.StringVar(&content, "content", "", "The content to work with (a valid directory path)")
    flag.StringVar(&filename, "filename", "", "A filename from the content")
    flag.Parse()
    if len(content) == 0 {
        log.Fatalln("'-content' path is empty")
    }
    if pathInfo, err := os.Stat(content); err != nil {
        log.Fatalln(err)
    } else if !pathInfo.IsDir() {
        log.Fatalln(content, "is not a valid dir")
    }
    switch mode {
        case MANAGER_MODE:
            fmt.Println("Manager mode")
            if filename != "" {
                manager.ManageDatabase(content, filename)
            }
        case BUILDER_MODE:
            fmt.Println("Builder mode")
        default:
            log.Fatalln("Specify '-mode' (manager or builder)")
    }
    flag.PrintDefaults()
    fmt.Println(flag.Arg(0), flag.Arg(1), flag.Arg(2))
}
