package builder

import (
    "fmt"
    "embed"
    "path/filepath"
    "os"
    "log"
    "text/template"

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

type PostsPageContext struct {
    CurrentPage PostsPage
    Website Website
}

func mkdir(base, ext string) {
    baseExtPath := filepath.Join(base, ext)
    if baseExtPathInfo, err := os.Stat(baseExtPath); os.IsNotExist(err) {
        fmt.Println("Create", baseExtPath)
        if err := os.MkdirAll(baseExtPath, os.ModeDir); err != nil {
            log.Fatalln(err)
        }
    } else if !baseExtPathInfo.IsDir() {
        log.Fatalln(base, "is not a dir!")
    }
}

func writeWriting(website Website, writing Writing, outputPath string, template *template.Template) {
    outputWritingPath := filepath.Join(outputPath, writing.Url())
    fmt.Println("output writing:", outputWritingPath)
    outputWriting, err := os.Create(outputWritingPath)
    defer outputWriting.Close()
    if err != nil {
        log.Fatalln(err)
    }
    if err := template.Execute(outputWriting, WritingContext{writing, website}); err != nil {
        log.Fatalln(err)
    }
    fmt.Println(writing)
}

func Build(base string) {
    //to output
    outputDirPath := "output"
    mkdir(outputDirPath, "posts")
    mkdir(outputDirPath, "pages")
    content := manager.OpenContent(base)
    fmt.Println(content)
    posts := content.GetPosts()
    pages := content.GetPages()
    fmt.Println(posts)
    // distribute posts (files) in pages
    const postsPerPage = 3
    website := NewWebsite("MySite", postsPerPage, posts, pages)
    fmt.Println("website pages:", website.PostsPages())
    postTemplate, err := template.ParseFS(basicTemplates, "templates/*.tmpl")
    if err != nil {
        log.Fatalln(err)
    }
    writingTemplate, err := template.ParseFS(writingTemplate, "templates/*.tmpl")
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
        if err := postTemplate.Execute(outputFile, PostsPageContext{page, website}); err != nil {
            log.Fatalln(err)
        }
        for _, writing := range page.Writings() {
            writeWriting(website, writing, outputDirPath, writingTemplate)
        }
    }
}
