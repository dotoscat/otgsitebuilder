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

type PageContext struct {
    CurrentPage Page
    Website Website
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
