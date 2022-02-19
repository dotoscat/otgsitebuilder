// Copyright 2022 Oscar Triano García

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
//"database/sql"
//"fmt"
//"log"
)

/*

// Element represents a member, or row, of the set name which this belongs.
type Element struct {
	set  Set
	id   int
	name string
	db   *sql.DB
}

// Set returns its parent
func (e Element) Set() Set {
	return e.set
}

func (e Element) Id() int {
	return e.id
}

func (e Element) Name() string {
	return e.name
}

// PostIn tells if a specific post is in this element as true or false
func (e Element) PostIn(post Post) bool {
	query := fmt.Sprintf("SELECT count(*) AS ISIN FROM %v_Post WHERE category_id = ? AND post_id = ?", e.set.Name())
	isIn := 0
	//fmt.Println("PostIn:", query, e.id, post.Id())
	if row := e.db.QueryRow(query, e.id, post.Id()); row.Err() != nil {
		log.Fatalln(row.Err())
	} else {
		row.Scan(&isIn)
	}
	return isIn >= 1 // Should be only 1
}

func (e Element) Posts() []Post {
	const QUERY = `SELECT id, name, date
FROM Post
JOIN Category_Post ON Category_Post.post_id = Post.id
JOIN Category ON Category.name = ? AND Category.id = Category_Post.category_id`
	posts := make([]Post, 0)
	rows, err := e.db.Query(QUERY, e.name)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		post := newPost(e.db)
		post.FillFromRows(rows, "")
		posts = append(posts, post)
	}
	return posts
}

*/
