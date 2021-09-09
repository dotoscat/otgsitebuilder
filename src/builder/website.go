package builder

import (
    "fmt"

    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//Website represents a website with its posts PostsPage and PostsPage.
type Website struct {
    postsPages PostsPages
    pages []Writing
    title string
    style string
}

//PostsPages returns its PostsPages.
func (w Website) PostsPages() PostsPages {
    return w.postsPages
}

//Pages returns its Pages.
func (w Website) Pages() []Writing {
    return w.pages
}

func (w Website) Title() string {
    return w.title
}

func (w Website) Style() string {
    return w.style
}

//NewWebsite returns info about the website.
func NewWebsite(title string, postsPerPage int, posts []manager.Post, pages []manager.Page) Website {
    nPages := len(posts) / postsPerPage
    postsExtraPage := len(posts) % postsPerPage
    extraPage := postsExtraPage > 0
    if extraPage {
        nPages++
    }
    var url string
    iPosts := 0
    postsPages := make(PostsPages, nPages)
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
        newPage := PostsPage{parent: &postsPages, index: iPage, url: url}
        postsPages[iPage] = newPage
        for i := 0; i < totalPosts; i++ {
            writing := NewWriting(&(posts[iPosts]), "posts")
            postsPages[iPage].addWriting(writing)
            iPosts++
        }
    }
    fmt.Println("nPages:", nPages, ";extraPage:", extraPage)
    nUserPages := len(pages) // no posts pages
    userPages := make([]Writing, nUserPages)
    for i := 0; i < nUserPages; i++ {
        userPages[i] = NewWriting(&(pages[i]), "/pages")
    }
    return Website{postsPages, userPages, title, ""}
}
