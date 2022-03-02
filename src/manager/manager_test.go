// Copyright 2021 Oscar Triano Garcí­a

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"testing"
    "os"
    "io/fs"
)

func copyFile(src, dst string) error {
    content, err := os.ReadFile(src)
    if err != nil {
        return err
    }
    return os.WriteFile(dst, content, fs.ModePerm)
}

func restoreDatabase() {
    copyFile("testdata/content/.metadata.old.db", "testdata/content/.metadata.db")
}

// list posts from the database
func TestPosts(t *testing.T) {
	// const CORE_2 = 8
	content := OpenContent("testdata/content")
	t.Log("content", content)
	batchCh := content.GetPostsByCategory(ALL, 3)
	i := 0
	for batch := range batchCh {
		t.Log(i, batch)
		//go func(){
		for post := range batch.Posts() {
			t.Log("post:", post)
		}
		//}()
		i++
	}
	//<-done
	if err := content.Close(); err != nil {
		t.Fatal(err)
	}
}

// Retrieve a single post from content
func TestPost(t *testing.T) {
	content := OpenContent("testdata/content")
    if exists, err := content.IsPost("one.md"); err != nil {
        t.Fatal(err)
    } else if post, err := content.GetPost("one.md"); exists == true {
        t.Log("exists? ", exists, ", prove it: ", post)
    } else if err != nil {
        t.Fatal(err)
    }
    if exists, err := content.IsPost("foo.md"); err != nil {
        t.Fatal(err)
    } else if exists == true {
        t.Fatal("'foo.md' must not exists.")
    }

}

func TestIndex (t *testing.T) {
	content := OpenContent("testdata/content")
    if err := content.IndexPosts(); err != nil {
        t.Fatal(err)
    }
    t.Log("Final indexed files")
    for batch := range content.GetPostsByCategory(ALL, SINGLE_PAGE) {
        for post := range batch.Posts() {
            t.Log("Final indexed post: ", post)
        }
    }
    restoreDatabase()
}

func TestCategories(t *testing.T) {
	content := OpenContent("testdata/content")
    t.Log("Categories:")
    if categories, err := content.Categories(); err != nil {
        t.Fatal(err)
    } else {
        for _, category := range categories {
            t.Log(category)
        }
    }
}

func TestPostsCategory(t *testing.T) {
	content := OpenContent("testdata/content")
    page := content.GetPostsByCategory("Second Cat", SINGLE_PAGE)
    for page := range page {
    for post := range page.Posts() {
        t.Log(post)
    }}
}

func TestOptionPost(t *testing.T) {
	content := OpenContent("testdata/content")
    options := FileOption{}
    options.RemoveCategories = []string{"Second Cat", "Third Cat", "meow"}
    if err := content.ModifyPost("one.md", options); err != nil {
        t.Fatal(err)
    }
    restoreDatabase()
}

