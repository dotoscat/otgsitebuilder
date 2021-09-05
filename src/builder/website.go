package builder

import (
    "fmt"

    "github.com/dotoscat/otgsitebuilder/src/manager"
)

//Website represents a website with its posts PostsPage and PostsPage.
type Website struct {
    postsPages PostsPages

}

//PostsPages returns its PostsPages.
func (w Website) PostsPages() PostsPages {
    return w.postsPages
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
    return Website{postsPages}
}
