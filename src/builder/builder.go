// Copyright 2021 Oscar Triano García

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builder

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
    "path"
	"path/filepath"
	"strings"
	//"text/template"

	"github.com/dotoscat/otgsitebuilder/src/manager"
)

//go:embed templates/base.tmpl
//go:embed templates/postspage.tmpl
var BasicTemplates embed.FS

//go:embed templates/base.tmpl
//go:embed templates/writing.tmpl
var WritingTemplates embed.FS

//go:embed templates/base.tmpl
//go:embed templates/setpage.tmpl
var SetTemplates embed.FS

type CategoryUrl struct {
    Category string
    Url string
}

type PageUrl struct {
    Reference string
    Url string
}

type Website struct {
    Title string
    Style string
    License string
    CategoriesUrl []CategoryUrl
    PagesUrl []PageUrl
    Output string
}

func NewWebsite(content manager.Content) Website {
    title := content.Title()
    license := content.License()
    categories := make([]CategoryUrl, 0)
    contentCategories, err := content.Categories()
    if err != nil {
        log.Fatalln(err)
    }
    for _, category := range contentCategories {
        url := "/" + strings.Join([]string{strings.ReplaceAll(category, " ", "-"), "page1.html"}, "/")
        categories = append(categories, CategoryUrl{category, url})
    }

    pages := make([]PageUrl, 0)
    for _, page := range content.GetPages() {
        var reference string
        if page.Reference() != "" {
            reference = page.Reference()
        } else {
            reference = strings.TrimSuffix(page.File().Name(), ".md")
        }
        url := "/pages/" + strings.ReplaceAll(reference, " ", "-") + ".html"
        pageUrl := PageUrl{reference, url}
        pages = append(pages, pageUrl)
    }

    output := content.Output()

    return Website{
        Title: title,
        Style: "",
        License: license,
        CategoriesUrl: categories,
        PagesUrl: pages,
        Output: output,
    }
}

func (w Website) HasStyle() bool {
    return w.Style != ""
}

func (w Website) Render(destination string) string {
    root := filepath.Join(destination, w.Output)
    // index
    return root
}

type Webpage struct {
    Website Website
    Url string
}

type Page struct {
    Webpage Webpage
    Page manager.Page
}

type Post struct {
    Webpage Webpage
    Post manager.Post
}

type PostsPage struct {
    Webpage Webpage
    batch manager.Batch
    pathPrefix string
    firstPageAsIndex bool
}

func NewPostsPage(website Website, batch manager.Batch,
                  pathPrefix string,
                  firstPageAsIndex bool) PostsPage {
    var url string
    if firstPageAsIndex == true && batch.Index() == 1 {
        url = "/index.html"
    } else {
        url = filepath.Join(pathPrefix, "page" + fmt.Sprint(batch.Index()) + ".html")
    }
    postsPage := PostsPage{
        Webpage: Webpage {
            Website: website,
            Url: url,
        },
        batch: batch,
        pathPrefix: pathPrefix,
        firstPageAsIndex: firstPageAsIndex,
    }
    return postsPage
}

func (p PostsPage) Posts() <-chan manager.Post {
    return p.batch.Posts()
}

func (p PostsPage) HasNext() bool {
    return p.batch.Index() < p.batch.TotalPages()
}

func (p PostsPage) NextUrl() string {
    if p.HasNext() == false {
        return ""
    }
    return path.Join("/", p.pathPrefix, "page" + fmt.Sprint(p.batch.Index() + 1) + ".html")
}

func (p PostsPage) HasPrevious() bool {
    return p.batch.Index() > 1
}

func (p PostsPage) PreviousUrl() string {
    index := p.batch.Index()
    if p.HasPrevious() == false {
        return ""
    }
    if index - 1 == 1 && p.firstPageAsIndex == true {
        return "/index.html"
    }
    return path.Join("/", p.pathPrefix, "page" + fmt.Sprint(p.batch.Index() - 1) + ".html")
}

func Mkdir(base, ext string) {
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

func CopyFile(src, dst string) {
	content, err := os.ReadFile(src)
	if err != nil {
		log.Fatalln("(read file)", err)
	}
	if err := os.WriteFile(dst, content, os.ModePerm); err != nil {
		log.Fatalln("(write file)", err)
	}
}

func CopyDir(src, dst string) {

	_dirWalker := func(path string, d fs.DirEntry, err error) error {
		if path == src {
			return nil
		}
		cleanedPath := strings.TrimPrefix(path, src)
		dstPath := filepath.Join(dst, cleanedPath)

		fmt.Println("path:", path, err, "; to dst:", cleanedPath, dstPath)
		if d.IsDir() {
			if err := os.MkdirAll(dstPath, os.ModeDir); err != nil {
				log.Fatalln(err)
			}
		} else {
			CopyFile(path, dstPath)
		}
		return nil
	}

	err := filepath.WalkDir(src, _dirWalker)
	fmt.Println("End walking", src, err)

}

/*

func WriteWriting(website Website, writing Writinger, outputPath string, template *template.Template) {
	outputWritingPath := filepath.Join(outputPath, writing.Url())
	fmt.Println("output writing path:", outputWritingPath)
	fmt.Println("output writing:", writing)
	outputWriting, err := os.Create(outputWritingPath)
	defer outputWriting.Close()
	if err != nil {
		log.Fatalln("output writing err:", err)
	}
	fmt.Println("Template execute:", *template, writing, *outputWriting)
	var context interface{}
	var writingInterface interface{} = writing
	switch writingInterface.(type) {
	case Writing:
		context = WritingContext{writing.(Writing), website}
	case PostWriting:
		//if err := template.Execute(outputWriting, PostWritingContext{writing.(PostWriting), website}); err != nil {
		//    log.Fatalln("output writing template err:", err)
		//}
		postData := writing.(PostWriting)
		context = PostWritingContext{postData, website}
		//return
	default:
		log.Fatalln("Default!")
	}
	if err := template.Execute(outputWriting, context); err != nil {
		log.Fatalln("output writing template err:", err)
	}
	fmt.Println("done writing:", writing)
}

*/
