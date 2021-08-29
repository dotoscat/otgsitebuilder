package builder

import (
    "fmt"

    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//Website represents a website with its posts pages and pages.
type Website struct {
    pages Pages
}

//Pages returns its pages.
func (w Website) Pages() Pages {
    return w.pages
}

//NewWebsite returns info about the website.
func NewWebsite(postsPerPage int, posts []manager.Post) Website {
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
            writing := NewWriting(&(posts[iPosts]), "posts")
            pages[iPage].addWriting(writing)
            iPosts++
        }
    }
    fmt.Println("nPages:", nPages, ";extraPage:", extraPage)
    return Website{pages}
}
