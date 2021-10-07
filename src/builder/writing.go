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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dotoscat/otgsitebuilder/src/manager"
	"github.com/gomarkdown/markdown"
)

//Writing stores a copy of the manager.File and a final url of the post
type Writing struct {
	file manager.Filer
	url  string
}

//NewWriting constructs a Writing value with a baseUrl to be used along with the the manager.File Name
func NewWriting(file manager.Filer, baseUrl string) Writing {
	fmt.Println("base url:", baseUrl)
	url := fmt.Sprint(baseUrl, "/", strings.Replace(file.Name(), ".md", ".html", -1))
	return Writing{file, url}
}

func (w Writing) File() manager.Filer {
	return w.file
}

//RenderHeader returns info about this file stored in the database to be used as a header
func (w Writing) RenderHeader() string {
	return w.file.Header()
}

//RenderContent returns HTML from a markdown format writing
func (w Writing) RenderContent() string {
	var content string
	if source, err := os.ReadFile(w.file.Path()); err != nil {
		log.Fatalln(err)
	} else {
		content = string(markdown.ToHTML(source, nil, nil))
	}
	return content
}

//Url returns final writing URL
func (w Writing) Url() string {
	return w.url
}

// RenderPartialContent returns up to 'n' characters from the markdown file
func (w Writing) RenderPartialContent(n int) string {
	content := w.RenderContent()
	if max := len(content); max < n || n <= 0 {
		return content[:max]
	}
	return content[:n]
}
