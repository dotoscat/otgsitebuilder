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
	// "fmt"

	"github.com/dotoscat/otgsitebuilder/src/manager"
)

// Set pages
// ie url /category/testing1.html

type ElementPage struct {
	element manager.Element
	url     string
	posts   []Writing
}

func newElementPage(element manager.Element, url string, writings []Writing) ElementPage {
	elementPage := ElementPage{element: element, url: url}
	for _, writing := range writings {
		file := writing.File()
		post := file.(*manager.Post)
		if element.PostIn(*post) {
			elementPage.posts = append(elementPage.posts, writing)
		}
	}
	return elementPage
}

func (s ElementPage) Name() string {
	return s.element.Name()
}

func (s ElementPage) Url() string {
	return s.url
}

//Website represents a website with its posts PostsPage and PostsPage.
type Website struct {
	postsPages PostsPages
	posts      []Writing
	pages      []Writing
	categories []ElementPage
	title      string
	style      string
}

func (w Website) Categories() []ElementPage {
	return w.categories
}

//PostsPages returns its PostsPages.
func (w Website) PostsPages() PostsPages {
	return w.postsPages
}

//Pages returns its Pages.
func (w Website) Pages() []Writing {
	return w.pages
}

//Posts returns all the posts
func (w Website) Posts() []Writing {
	return w.posts
}

func (w Website) Title() string {
	return w.title
}

func (w *Website) SetStyle(style string) {
	w.style = style
}

func (w Website) Style() string {
	return w.style
}

func (w Website) HasStyle() bool {
	return w.style != ""
}

//NewWebsite returns info about the website.
func NewWebsite(title string, postsPerPage int, posts []manager.Post, pages []manager.Page, content manager.Content) Website {
	postsWritings := make([]Writing, len(posts))
	for i, post := range posts {
		postsWritings[i] = NewWriting(&post, "posts")
	}
	postsPages := NewPostsPages(postsPerPage, postsWritings, "index")
	nWebsitePages := len(pages) // no posts pages
	websitePages := make([]Writing, nWebsitePages)
	for i := 0; i < nWebsitePages; i++ {
		websitePages[i] = NewWriting(&(pages[i]), "/pages")
	}
	categories := make([]ElementPage, 0)
	for _, element := range content.Categories().Elements() {
		url := "/categories/" + element.Name() + ".html"
		elementPage := newElementPage(element, url, postsWritings)
		categories = append(categories, elementPage)
	}
	return Website{postsPages, postsWritings, websitePages, categories, title, ""}
}
