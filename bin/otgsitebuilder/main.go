package main

import (
    "log"
    "fmt"
    "flag"
    "os"
    "io/fs"
    "errors"
    "strings"
    "time"
    "strconv"

    "github.com/dotoscat/otgsitebuilder/src/manager"
    "github.com/dotoscat/otgsitebuilder/src/builder"
    )

const (
    MANAGER_MODE = "manager"
    BUILDER_MODE = "builder"
)

type DateValue struct {
    time time.Time
    requested bool
}

func (dv DateValue) String() string {
    year, month, day := dv.time.Date()
    return fmt.Sprintf("%v-%02v-%02v", year, int(month), day)
}

func (dv *DateValue) Set(value string) error {
    if len(value) == 0 {
        dv.requested = false
        return nil
    }
    parts := strings.Split(value, "-")
    if len(parts) != 3 {
        return errors.New("A date has a year, a month and a day separated by '-' (YYYY-M-D) ")
    }
    year, err := strconv.Atoi(parts[0])
    if err != nil {
        return err
    }
    month, err := strconv.Atoi(parts[1])
    if err != nil {
        return err
    }
    day, err := strconv.Atoi(parts[2])
    if err != nil {
        return err
    }
    dv.time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
    dv.requested = true
    fmt.Println("requested:", dv.requested)
    return nil
}

func (dv DateValue) IsRequested() bool {
    return dv.requested
}

type FlagList struct {
    Mode string
    Content string
    Filename string
    Date DateValue
}

func managePost(post manager.Post, flagList FlagList) {
    if flagList.Date.IsRequested(){
        fmt.Println("date is request for post:", post, "===")
        if err := post.SetDate(flagList.Date.String()); err != nil {
            log.Fatalln(err)
        }
    }
}

func manageDatabase(flagList FlagList) {
    content := manager.OpenContent(flagList.Content)
    fmt.Println("content: ", content)
    if isPost, err := content.CheckInPostsFolder(flagList.Filename); err != nil && !errors.Is(err, fs.ErrNotExist) {
        log.Fatalln("Is not 'ErrNotExist'", err)
    } else if isPost {
        post := content.GetPostFile(flagList.Filename)
        managePost(post, flagList)
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
    flag.Var(&flagList.Date, "date", "Set a date, in YYYY-M-D format, for a post")
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
