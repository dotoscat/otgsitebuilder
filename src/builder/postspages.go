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

/*

import (
	"fmt"
)

//PostsPages is defined type for a slice of PostsPage.
type PostsPages []PostsPage

func NewPostsPages(postsPerPage int, posts []PostWriting, base string) PostsPages {
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
			url = fmt.Sprintf("/%v.html", base)
		} else {
			url = fmt.Sprintf("/%v%v.html", base, iPage)
		}
		newPage := PostsPage{parent: &postsPages, index: iPage, url: url}
		postsPages[iPage] = newPage
		for i := 0; i < totalPosts; i++ {
			postsPages[iPage].addWriting(posts[iPosts])
			iPosts++
		}
	}
	fmt.Println("nPages:", nPages, ";extraPage:", extraPage)
	return postsPages
}

*/
