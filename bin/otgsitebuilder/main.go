package main

import (
    "log"
    "fmt"
    "flag"
    "os"
    "io/fs"
    "errors"

    "github.com/dotoscat/otgsitebuilder/src/manager"
    "github.com/dotoscat/otgsitebuilder/src/builder"
    )

const (
    MANAGER_MODE = "manager"
    BUILDER_MODE = "builder"
)

type FlagList struct {
    Mode string
    Content string
    Filename string
}

func manageDatabase(flagList FlagList) {
    content := manager.OpenContent(flagList.Content)
    fmt.Println("content: ", content)
    if isPost, err := content.CheckInPostsFolder(flagList.Filename); err != nil && !errors.Is(err, fs.ErrNotExist) {
        log.Fatalln("Is not 'ErrNotExist'", err)
    } else if isPost {
        post := content.GetPostFile(flagList.Filename)
        fmt.Println("post:", post)
    } else if isPage, err := content.CheckInPagesFolder(flagList.Filename); err != nil && !errors.Is(err, fs.ErrNotExist) {
        log.Fatalln("Is not 'ErrNotExist'", err)
    } else if isPage {
        page := content.GetPageFile(flagList.Filename)
        fmt.Println("page:", page)
    } else {
        fmt.Println(flagList.Filename, "does not exist in", flagList.Content)
    }
}

func main() {
    flagList := FlagList{}
    flag.StringVar(&flagList.Mode, "mode", "", "Set the mode of use of this tool")
    flag.StringVar(&flagList.Content, "content", "", "The content to work with (a valid directory path)")
    flag.StringVar(&flagList.Filename, "filename", "", "A filename from the content")
    flag.Parse()
    if len(flagList.Content) == 0 {
        log.Fatalln("'-content' path is empty")
    }
    if pathInfo, err := os.Stat(flagList.Content); err != nil {
        log.Fatalln(err)
    } else if !pathInfo.IsDir() {
        log.Fatalln(flagList.Content, "is not a valid dir")
    }
    switch flagList.Mode {
        case MANAGER_MODE:
            fmt.Println("Manager mode")
            if flagList.Filename != "" {
                manageDatabase(flagList)
            }
        case BUILDER_MODE:
            builder.Build(flagList.Content)
            fmt.Println("Builder mode")
        default:
            log.Fatalln("Specify '-mode' (manager or builder)")
    }
    flag.PrintDefaults()
    fmt.Println(flag.Arg(0), flag.Arg(1), flag.Arg(2))
}
