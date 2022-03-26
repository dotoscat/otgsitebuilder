package builder

import (
    "testing"

	"github.com/dotoscat/otgsitebuilder/src/manager"
)

func TestNewWebsite(t *testing.T) {
    content := manager.OpenContent("testdata/content")
    website := NewWebsite(content)
    t.Log(website)
}

func TestNewPostsPage (t *testing.T) {
    content := manager.OpenContent("testdata/content")
    website := NewWebsite(content)

    for batch := range content.GetPostsByCategory(manager.ALL, content.PostsPerPage()) {
        postsPage := NewPostsPage(website, batch, "/pagestest", false)
        t.Log(postsPage)
    }

    for batch := range content.GetPostsByCategory(manager.ALL, content.PostsPerPage()) {
        var postsPage PostsPage
        postsPage = NewPostsPage(website, batch, "/pagestest", true)
        t.Log(batch.Index(), "; previous: ", postsPage.PreviousUrl(), ", next:", postsPage.NextUrl())
    }
}
