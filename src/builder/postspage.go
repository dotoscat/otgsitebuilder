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

//PostsPage stores info about the PostsPages which it belongs, index,
//a list of writings related to this PostsPage and its url.
type PostsPage struct {
	parent   *PostsPages
	index    int
	writings []Writing
	url      string
}

//PostsPages is defined type for a slice of PostsPage.
type PostsPages []PostsPage

//HasLast returns true if this is not the last PostsPage.
func (p PostsPage) HasLast() bool {
	return p.index-1 >= 0
}

//HasNext returns true if there are more PostsPages.
func (p PostsPage) HasNext() bool {
	return p.index+1 < len(*p.parent)
}

//Last returns the last PostsPage on this PostsPage.
//Returns an empty PostsPage if not.
func (p PostsPage) Last() PostsPage {
	if p.HasLast() {
		return (*p.parent)[p.index-1]
	}
	return PostsPage{}
}

//Next returns the next PostsPage on this PostsPage.
//Returns an empty PostsPage if is not possible.
func (p PostsPage) Next() PostsPage {
	if p.HasNext() {
		return (*p.parent)[p.index+1]
	}
	return PostsPage{}
}

//Url returns the final url of this PostsPage.
func (p PostsPage) Url() string {
	return p.url
}

//Empty tells if this PostsPage is empty, does not belong to PostsPages.
func (p PostsPage) Empty() bool {
	return p.parent == nil
}

//Writinngs returns the writings stored in this PostsPage
func (p PostsPage) Writings() []Writing {
	return p.writings
}

//addWriting is a convenient method to add writings to this PostsPage.
//The return value is the same writing passed as argument.
func (p *PostsPage) addWriting(writing Writing) Writing {
	p.writings = append(p.writings, writing)
	return writing
}
