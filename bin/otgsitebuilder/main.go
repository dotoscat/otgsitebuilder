package main

import (
    "log"
    "fmt"
    "flag"
    )

const (
    MANAGER_MODE = "manager"
    BUILDER_MODE = "builder"
)

func main() {
    var mode string
    var content string
    flag.StringVar(&mode, "mode", "", "Set the mode of use of this tool")
    flag.StringVar(&content, "content", "", "The content to work with (a valid directory path)")
    flag.Parse()
    switch (mode) {
        case MANAGER_MODE:
            fmt.Println("Manager mode")
        case BUILDER_MODE:
            fmt.Println("Builder mode")
        default:
            log.Fatalln("Specify '-mode' (manager or builder)")
    }
    if len(content) == 0 {
        log.Fatalln("'-content' path is empty")
    }
    flag.PrintDefaults()
    fmt.Println(flag.Arg(0), flag.Arg(1), flag.Arg(2))
}
