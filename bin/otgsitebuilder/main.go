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
    "path/filepath"
    "text/template"

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
    Reference string
    RemoveReference bool
    Theme string
}

func managePost(post manager.Post, flagList FlagList) {
    if flagList.Date.IsRequested(){
        fmt.Println("date is request for post:", post, "===")
        if err := post.SetDate(flagList.Date.String()); err != nil {
            log.Fatalln(err)
        }
    }
}

func managePage(page manager.Page, flagList FlagList) {
    if flagList.RemoveReference {
        page.SetReference("")
    } else if flagList.Reference != "-1" {
        fmt.Println("Set 'reference':", flagList.Reference)
        page.SetReference(flagList.Reference)
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
        managePage(page, flagList)
        fmt.Println("page:", page)
    } else {
        fmt.Println(flagList.Filename, "does not exist in", flagList.Content)
    }
}

func build(base string, flags FlagList) {
    //to output
    outputDirPath := "output"
    staticDirPath := filepath.Join(outputDirPath, "static")
    builder.Mkdir(outputDirPath, "posts")
    builder.Mkdir(outputDirPath, "pages")
    builder.Mkdir(outputDirPath, "static")
    content := manager.OpenContent(base)
    fmt.Println(content)
    posts := content.GetPosts()
    pages := content.GetPages()
    fmt.Println(posts)
    // distribute posts (files) in pages
    const postsPerPage = 3
    website := builder.NewWebsite("MySite", postsPerPage, posts, pages)
    fmt.Println("website pages:", website.PostsPages())
    postTemplate, err := template.ParseFS(builder.BasicTemplates, "templates/*.tmpl")
    if err != nil {
        log.Fatalln(err)
    }
    writingTemplate, err := template.ParseFS(builder.WritingTemplates, "templates/*.tmpl")
    if err != nil {
        log.Fatalln(err)
    }
    for i, page := range website.PostsPages() {
        var outputFilePath string
        if i == 0 {
            outputFilePath = filepath.Join(outputDirPath, "index.html")
        } else {
            outputFilePath = filepath.Join(outputDirPath, fmt.Sprint("index", i, ".html"))
        }
        outputFile, err := os.Create(outputFilePath)
        defer outputFile.Close()
        if err != nil {
            log.Fatalln(err)
        }
        if err := postTemplate.Execute(outputFile, builder.PostsPageContext{page, website}); err != nil {
            log.Fatalln(err)
        }
        for _, writing := range page.Writings() {
            builder.WriteWriting(website, writing, outputDirPath, writingTemplate)
        }
    }
    fmt.Println("RENDER PAGES")
    fmt.Println("pages:", website.Pages())
    for _, writing := range website.Pages() {
        fmt.Println("Render page url:", writing.Url(), outputDirPath)
        builder.WriteWriting(website, writing, outputDirPath, writingTemplate)
    }
    fmt.Println("DONE")
    builder.CopyDir(filepath.Join(base, "static"), staticDirPath)
    // render user pages, no posts pages
}

func main() {
    flagList := FlagList{}
    flag.StringVar(&flagList.Mode, "mode", "", "Set the mode of use of this tool")
    flag.StringVar(&flagList.Content, "content", "", "The content to work with (a valid directory path)")
    flag.StringVar(&flagList.Filename, "filename", "", "A filename from the content")
    flag.StringVar(&flagList.Reference, "reference", "-1", "Set a reference for a page instead its name")
    flag.StringVar(&flagList.Theme, "theme", "", "Set the theme (a style sheet) to use for building the site")
    flag.BoolVar(&flagList.RemoveReference, "remove-reference", false, "Remove reference")
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
            build(flagList.Content, flagList)
            fmt.Println("Builder mode")
        default:
            log.Fatalln("Specify '-mode' (manager or builder)")
    }
    flag.PrintDefaults()
    fmt.Println(flag.Arg(0), flag.Arg(1), flag.Arg(2))
}
