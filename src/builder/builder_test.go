package builder

import (
    "embed"
    "text/template"
    "strings"
    "testing"

	"github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed templates/base.tmpl
var base embed.FS

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
    var baseTemplate *template.Template
    var err error
    if baseTemplate, err = template.ParseFS(base, "templates/*.tmpl"); err != nil {
        t.Fatal(err)
    }

    for batch := range content.GetPostsByCategory(manager.ALL, content.PostsPerPage()) {
        var postsPage PostsPage
        postsPage = NewPostsPage(website, batch, "/pagestest", true)
        t.Log(batch.Index(), "; previous: ", postsPage.PreviousUrl(), ", next:", postsPage.NextUrl())
        builder := &strings.Builder{}
        if err := baseTemplate.Execute(builder, postsPage); err != nil {
            t.Fatal(err)
        }
        t.Log(builder)
    }
}
