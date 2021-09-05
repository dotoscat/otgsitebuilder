package builder

import (
    "fmt"
    "strings"
    "os"
    "log"

    "github.com/gomarkdown/markdown"
    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//Writing stores a copy of the manager.File and a final url of the post
type Writing struct{
    manager.Filer
    url string
}

//NewWriting constructs a Writing value with a baseUrl to be used along with the the manager.File Name
func NewWriting(file manager.Filer, baseUrl string) Writing {
    fmt.Println("base url:", baseUrl)
    url := fmt.Sprint(baseUrl, "/", strings.Replace(file.Name(), ".md", ".html", -1))
    return Writing{file, url}
}

//RenderHeader returns info about this file stored in the database to be used as a header
func (w Writing) RenderHeader() string {
   return w.Header()
}

//RenderContent returns HTML from a markdown format writing
func (w Writing) RenderContent() string{
    var content string
    if source, err := os.ReadFile(w.Path()); err != nil {
        log.Fatalln(err)
    } else {
        content = string(markdown.ToHTML(source, nil, nil))
    }
    return content
}

//Url returns final writing URL
func (w Writing) Url() string {
    return w.url
}

// RenderPartialContent returns up to 'n' characters from the markdown file
func (w Writing) RenderPartialContent(n int) string {
    content := w.RenderContent()
    if max := len(content); max < n || n <= 0{
        return content[:max]
    }
    return content[:n]
}
