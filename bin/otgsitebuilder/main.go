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

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dotoscat/otgsitebuilder/src/builder"
	"github.com/dotoscat/otgsitebuilder/src/manager"
)

const (
	MANAGER_MODE = "manager"
	BUILDER_MODE = "builder"
	VERSION      = "0.2.0"
)

//DataValue is special flag to handle time
type DateValue struct {
	time      time.Time
	requested bool
}

//String returns the date as a string in year-month-day format
func (dv DateValue) String() string {
	year, month, day := dv.time.Date()
	return fmt.Sprintf("%v-%02v-%02v", year, int(month), day)
}

//Set is used internally by the standard flag package
func (dv *DateValue) Set(value string) error {
	if len(value) == 0 {
		dv.requested = false
		return nil
	}
	parts := strings.Split(value, "-")
	if len(parts) != 3 {
		return errors.New("A date has a year, a month and a day separated by '-' (YYYY-M-D) ")
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	dv.time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	dv.requested = true
	fmt.Println("requested:", dv.requested)
	return nil
}

//IsRequested returns a boolean if this flag is provided as a parameter
func (dv DateValue) IsRequested() bool {
	return dv.requested
}

//FlagList stores the flags to be passed to different functions
type FlagList struct {
	Mode            string
	Content         string
	Filename        string
	Date            DateValue
	Reference       string
	RemoveReference bool
	RemoveCategory  bool
	Version         bool
	Theme           string
	Title           string
	PostsPerPage    int
	Output          string
	Category        string
	License         string
}

//managePost is the main point entry to manage a post
func managePost(post manager.Post, flagList FlagList, content manager.Content) {
	if flagList.Date.IsRequested() {
		fmt.Println("date is request for post:", post, "===")
		if err := post.SetDate(flagList.Date.String()); err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("remove categories:", flagList.RemoveCategory)
	fmt.Println("category to remove:", flagList.Category)
	if flagList.Category != "" && flagList.RemoveCategory == false {
		if err := content.Categories().AddPostForElement(post, flagList.Category); err != nil {
			log.Fatalln(err)
		}
	} else if flagList.Category != "" && flagList.RemoveCategory == true {
		categories := content.Categories()
		if err := categories.RemovePostForElement(post, flagList.Category); err != nil {
			log.Fatalln(err)
		} else {
			categories.DeleteUnusedElements()
		}
	}
}

//managePost is the main point entry to manage a page
func managePage(page manager.Page, flagList FlagList) {
	if flagList.RemoveReference {
		page.SetReference("")
	} else if flagList.Reference != "-1" {
		fmt.Println("Set 'reference':", flagList.Reference)
		page.SetReference(flagList.Reference)
	}
}

//manageDatabase is the main point entry to manage the metadata (database)
func manageDatabase(flagList FlagList) {
	fmt.Println("Category in manager:", flagList.Category)
	content := manager.OpenContent(flagList.Content)
	fmt.Println("categories:", content.Categories().Elements())
	// fmt.Println("content: ", content)
	if flagList.PostsPerPage > 0 {
		fmt.Println("Set posts per page:", flagList.PostsPerPage)
		if err := content.SetPostsPerPage(flagList.PostsPerPage); err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println("Posts per page:", content.PostsPerPage())
	}
	if flagList.Title != "" {
		content.SetTitle(flagList.Title)
	} else {
		fmt.Println("Title:", content.Title())
	}
	if flagList.License != "" {
		content.SetLicense(flagList.License)
		fmt.Println("License:", content.License())
	}
	if flagList.Output != "" {
		if err := content.SetOutput(flagList.Output); err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println("Output:", content.Output())
	}
	if flagList.Filename == "" {
		return
	}
	if isPost, err := content.CheckInPostsFolder(flagList.Filename); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalln("Is not 'ErrNotExist'", err)
	} else if isPost {
		post := content.GetPostFile(flagList.Filename)
		managePost(post, flagList, content)
		fmt.Println("post:", post)
	} else if isPage, err := content.CheckInPagesFolder(flagList.Filename); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalln("Is not 'ErrNotExist'", err)
	} else if isPage {
		page := content.GetPageFile(flagList.Filename)
		managePage(page, flagList)
		fmt.Println("page:", page)
	} else {
		fmt.Println(flagList.Filename, "does not exist in", flagList.Content)
	}
}

// main entry for building a website
func build(flags FlagList) {
	base := flags.Content
	content := manager.OpenContent(base)
	fmt.Println(content)
	//to output
	outputDirPath := content.Output()
	staticDirPath := filepath.Join(outputDirPath, "static")
	builder.Mkdir(outputDirPath, "posts")
	builder.Mkdir(outputDirPath, "pages")
	builder.Mkdir(outputDirPath, "static")
	builder.Mkdir(outputDirPath, "categories")

	if flags.Theme != "" {
		switch {
		case fs.ValidPath(flags.Theme) == false:
			log.Fatalln(flags.Theme, "is not a valid path!")
		case strings.HasSuffix(flags.Theme, ".css") == false:
			log.Fatalln(flags.Theme, "is not a valid css")
		}
		builder.CopyFile(flags.Theme, filepath.Join(outputDirPath, filepath.Base(flags.Theme)))
	}

	posts := content.GetPosts()
	pages := content.GetPages()
	fmt.Println(posts)
	website := builder.NewWebsite(content.Title(), content.PostsPerPage(), posts, pages, content)
	if flags.Theme != "" {
		website.SetStyle(filepath.Join("/", filepath.Base(flags.Theme)))
	}
	fmt.Println("website pages:", website.PostsPages())
	postTemplate, err := template.ParseFS(builder.BasicTemplates, "templates/*.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	writingTemplate, err := template.ParseFS(builder.WritingTemplates, "templates/*.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	setpageTemplate, err := template.ParseFS(builder.SetTemplates, "templates/*.tmpl")
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
		if err := postTemplate.Execute(outputFile, builder.PostsPageContext{page, website}); err != nil {
			log.Fatalln(err)
		}
	}
	for _, writing := range website.Posts() {
		builder.WriteWriting(website, writing, outputDirPath, writingTemplate)
	}
	fmt.Println("RENDER PAGES")
	fmt.Println("pages:", website.Pages())
	for _, writing := range website.Pages() {
		fmt.Println("Render page url:", writing.Url(), outputDirPath)
		fmt.Println("Page content:", writing.RenderContent())
		builder.WriteWriting(website, writing, outputDirPath, writingTemplate)
	}
	for _, element := range website.Categories() {
		outputFilePath := filepath.Join(outputDirPath, element.Url())
		outputFile, err := os.Create(outputFilePath)
		defer outputFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
		context := builder.SetPageContext{element, website}
		if err := setpageTemplate.Execute(outputFile, context); err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("DONE")
	// fmt.Println("categories", website.Categories())
	builder.CopyDir(filepath.Join(base, "static"), staticDirPath)
	// render user pages, no posts pages
}

func main() {
	flagList := FlagList{}
	flag.StringVar(&flagList.Mode, "mode", "", "Set the mode of use of this tool")
	flag.StringVar(&flagList.Content, "content", "", "The content to work with (a valid directory path)")
	flag.StringVar(&flagList.Filename, "filename", "", "A filename from the content")
	flag.StringVar(&flagList.Reference, "reference", "-1", "Set a reference for a page instead its name")
	flag.StringVar(&flagList.Theme, "theme", "", "Set the theme (a style sheet) to use for building the site")
	flag.StringVar(&flagList.Title, "title", "", "Set the title to use for building the site")
	flag.StringVar(&flagList.Output, "output", "", "Set the output of the build process")
	flag.StringVar(&flagList.Category, "category", "", "Set the category for a post")
	flag.StringVar(&flagList.License, "license", "", "Set copyright for the website")
	flag.IntVar(&flagList.PostsPerPage, "posts-per-page", -1, "Set the posts per page for building the site")
	flag.BoolVar(&flagList.RemoveReference, "remove-reference", false, "Remove reference")
	flag.BoolVar(&flagList.RemoveCategory, "remove-category", false, "Remove category for post")
	flag.BoolVar(&flagList.Version, "version", false, "Show current version")
	flag.Var(&flagList.Date, "date", "Set a date, in YYYY-M-D format, for a post")
	flag.Parse()
	if flagList.Version == true {
		fmt.Println("version:", VERSION)
		return
	}
	if len(flagList.Content) == 0 {
		log.Fatalln("'-content' path is empty")
	}
	if pathInfo, err := os.Stat(flagList.Content); err != nil {
		log.Fatalln(err)
	} else if !pathInfo.IsDir() {
		log.Fatalln(flagList.Content, "is not a valid dir")
	}
	switch flagList.Mode {
	case MANAGER_MODE:
		fmt.Println("Manager mode")
		manageDatabase(flagList)
	case BUILDER_MODE:
		build(flagList)
		fmt.Println("Builder mode")
	default:
		log.Fatalln("Specify '-mode' (manager or builder)")
	}
	flag.PrintDefaults()
}
