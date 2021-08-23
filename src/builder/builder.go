package builder

import (
    "fmt"
    "embed"
    "path/filepath"
    "os"
    "log"
    "strings"
    "text/template"

    "github.com/gomarkdown/markdown"
    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed templates/base.tmpl
//go:embed templates/postspage.tmpl
var basicTemplates embed.FS

//go:embed templates/base.tmpl
//go:embed templates/writing.tmpl
var writingTemplate embed.FS

type WritingContext struct {
    Writing Writing
    Website Website
}

type Page struct {
    parent *Pages
    index int
    writings []Writing
    url string
}

func (p Page) HasLast() bool {
    return p.index - 1 >= 0
}

func (p Page) HasNext() bool {
    return p.index + 1 < len(*p.parent)
}

func (p Page) Last() Page {
    if p.HasLast() {
        return (*p.parent)[p.index-1]
    }
    return Page{}
}

func (p Page) Next() Page {
    if p.HasNext() {
        return (*p.parent)[p.index+1]
    }
    return Page{}
}

func (p Page) Url() string {
    return p.url
}

func (p Page) Empty() bool {
    return p.parent == nil
}

func (p Page) Writings() []Writing {
    return p.writings
}

func (p *Page) addWriting(writing Writing) Writing {
    p.writings = append(p.writings, writing)
    return writing
}

type PageContext struct {
    CurrentPage Page
    Website Website
}

type Pages []Page

type Website struct {
    pages Pages
}

func (w Website) Pages() Pages {
    return w.pages
}

func NewWebsite(postsPerPage int, posts []manager.File) Website {
    nPages := len(posts) / postsPerPage
    postsExtraPage := len(posts) % postsPerPage
    extraPage := postsExtraPage > 0
    if extraPage {
        nPages++
    }
    var url string
    iPosts := 0
    pages := make(Pages, nPages)
    for iPage := 0; iPage < nPages; iPage++ {
        var totalPosts int
        if iPage == nPages-1 && extraPage {
            totalPosts = postsExtraPage
        } else {
            totalPosts = postsPerPage
        }
        if iPage == 0 {
            url = "/index.html"
        } else {
            url = fmt.Sprint("/index", iPage, ".html")
        }
        newPage := Page{parent: &pages, index: iPage, url: url}
        pages[iPage] = newPage
        for i := 0; i < totalPosts; i++ {
            writing := NewWriting(posts[iPosts], "posts")
            pages[iPage].addWriting(writing)
            iPosts++
        }
    }
    fmt.Println("nPages:", nPages, ";extraPage:", extraPage)
    return Website{pages}
}

func Build(base string) {
    //to output
    outputDirPath := "output"
    postsOutputDirPath := filepath.Join(outputDirPath, "posts")
    if outputDirInfo, err := os.Stat(outputDirPath); os.IsNotExist(err) {
        fmt.Println("Create", outputDirPath)
        if err := os.MkdirAll(outputDirPath, os.ModeDir); err != nil {
            log.Fatalln(err)
        }
    } else if !outputDirInfo.IsDir() {
        log.Fatalln(outputDirPath, "is not a dir!")
    }
    if err := os.MkdirAll(postsOutputDirPath, os.ModeDir); err != nil {
        log.Fatalln(err)
    }
    content := manager.OpenContent(base)
    fmt.Println(content)
    posts := content.GetPosts()
    fmt.Println(posts)
    // distribute posts (files) in pages
    const postsPerPage = 3
    website := NewWebsite(postsPerPage, posts)
    fmt.Println("website pages:", website.Pages())
    postTemplate, err := template.ParseFS(basicTemplates, "templates/*.tmpl")
    if err != nil {
        log.Fatalln(err)
    }
    writingTemplate, err := template.ParseFS(writingTemplate, "templates/*.tmpl")
    if err != nil {
        log.Fatalln(err)
    }
    for i, page := range website.Pages() {
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
        if err := postTemplate.Execute(outputFile, PageContext{page, website}); err != nil {
            log.Fatalln(err)
        }
        for _, writing := range page.Writings() {
            outputWritingPath := filepath.Join(outputDirPath, writing.Url())
            fmt.Println("output writing:", outputWritingPath)
            outputWriting, err := os.Create(outputWritingPath)
            defer outputWriting.Close()
            if err != nil {
                log.Fatalln(err)
            }
            if err := writingTemplate.Execute(outputWriting, WritingContext{writing, website}); err != nil {
                log.Fatalln(err)
            }
            fmt.Println(writing)
        }
    }
}
