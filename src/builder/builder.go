package builder

import (
    "fmt"
    "embed"
    "path/filepath"
    "os"
    "log"
    "html/template"

    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed templates/base.tmpl
//go:embed templates/postspage.tmpl
var basicTemplates embed.FS

type Writing struct{
    manager.File
}

func NewWriting(file manager.File) Writing {
    return Writing{file}
}

type Page struct {
    writings []Writing
    // Last
    // Next
}

func (p *Page) addWriting(writing Writing) Writing {
    p.writings = append(p.writings, writing)
    return writing
}

func (p *Page) Writings() []Writing {
    return p.writings
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
    iPosts := 0
    pages := make(Pages, nPages)
    for iPage := 0; iPage < nPages; iPage++ {
        var totalPosts int
        if iPage == nPages-1 && extraPage {
            totalPosts = postsExtraPage
        } else {
            totalPosts = postsPerPage
        }
        pages[iPage] = Page{}
        for i := 0; i < totalPosts; i++ {
            pages[iPage].addWriting(NewWriting(posts[iPosts]))
            iPosts++
        }
    }
    fmt.Println("nPages:", nPages, ";extraPage:", extraPage)
    return Website{pages}
}

func Build(base string) {
    //to output
    outputDirPath := "output"
    if outputDirInfo, err := os.Stat(outputDirPath); os.IsNotExist(err) {
        fmt.Println("Create", outputDirPath)
        if err := os.MkdirAll(outputDirPath, os.ModeDir); err != nil {
            log.Fatalln(err)
        }
    } else if !outputDirInfo.IsDir() {
        log.Fatalln(outputDirPath, "is not a dir!")
    }
    content := manager.OpenContent(base)
    fmt.Println(content)
    posts := content.GetPosts()
    fmt.Println(posts)
    // distribute posts (files) in pages
    const postsPerPage = 3
    website := NewWebsite(postsPerPage, posts)
    fmt.Println(website)
    postTemplate, err := template.ParseFS(basicTemplates, "templates/*.tmpl")
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
        if err := postTemplate.Execute(outputFile, page); err != nil {
            log.Fatalln(err)
        }
    }
}
