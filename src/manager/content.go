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
	"database/sql"
	"os"
	"fmt"
	"log"
	"path/filepath"
)

type Content struct {
	db        *sql.DB
	postsPath string
	pagesPath string
}

func (c Content) Close() error {
	return c.db.Close()
}

func (c Content ) Categories() ([]string, error) {
    return []string{}, nil
}

func (c Content) index(table, sourcePath string) error {
    entries, err := os.ReadDir(sourcePath)
	if err != nil {
        return err
	}
	QUERY := fmt.Sprintf("SELECT EXISTS (SELECT name FROM %v WHERE name = ?)", table)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fmt.Println("Index: ", entry.Name())
        if exists, err := c.exists(QUERY, entry.Name()); err != nil {
            return err
        } else if exists == true {
            fmt.Println(entry.Name(), "indexed!")
            continue
        }
        QUERY := fmt.Sprintf("INSERT INTO %v (name) VALUES (?)", table)
        if result, err := c.db.Exec(QUERY, entry.Name()); err != nil {
            return err
        } else if _, err := result.LastInsertId(); err != nil {
            return err
        }
	}
	return nil
}

func (c Content) IndexPosts() error {
    return c.index("Post", c.postsPath)
}

func (c Content) IndexPages() error {
    return c.index("Page", c.pagesPath)
}

func (c Content) IndexFiles() (bool, error) {
    return false, nil
}

func (c Content) exists(query, name string) (bool, error) {
    row := c.db.QueryRow(query, name)
    if row.Err() != nil {
        return false, row.Err()
    }
    var exists int64
    if err := row.Scan(&exists); err != nil {
        return false, err
    }
    if exists == 1 {
        return true, nil
    }
    return false, nil
}

func (c Content) IsPost(name string) (bool, error) {
    const QUERY = "SELECT EXISTS (SELECT name FROM Post WHERE name = ?)"
    return c.exists(QUERY, name)
}

func (c Content) IsPage(name string) (bool, error) {
    const QUERY = "SELECT EXISTS (SELECT name FROM Page WHERE name = ?)"
    return c.exists(QUERY, name)
}

func (c Content) get(query, name string, dest ...interface{}) error {
    row := c.db.QueryRow(query, name)
    if row.Err() != nil {
        return row.Err()
    }
    if err := row.Scan(dest...); err != nil {
        return err
    }
    return nil
}

func (c Content) GetPost(name string) (Post, error) {
    const QUERY = "SELECT id, name, date FROM Post WHERE name = ?"
    post := Post{}
    err := c.get(QUERY, name, &post.file.id, &post.file.name, &post.date)
    fmt.Println("GetPost", post)
    return post, err
}

func (c Content) GetPage(name string) (Page, error) {
    const QUERY = "SELECT id, name, reference FROM Page WHERE name = ?"
    page := Page{}
    err := c.get(QUERY, name, &page.file.id, &page.file.name, &page.reference)
    return page, err
}

func (c Content) ModifyPost(name string, options FileOption) (bool, error) {
    return false, nil
}

func (c Content) ModifyPage(name string, options FileOption) (bool, error) {
    return false, nil
}

func (c Content) GetPages() []Page {
	const QUERY = "SELECT id, name, reference FROM Page"
	rows, err := c.db.Query(QUERY)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	files := make([]Page, 0)
	for rows.Next() {
		page := Page{}
		err := rows.Scan(&page.file.id, &page.file.name, &page.reference)
		if err != nil {
			log.Fatalln(err)
		}
		files = append(files, page)
	}
	return files
}

const ALL = ""

// GetPostsByCategory returns batch from
func (c Content) GetPostsByCategory(category string, postsPerPage int) <-chan Batch {
	if postsPerPage <= 0 {
		postsPerPage = 3
	}

	const QUERY_ALL = "SELECT id, name, date FROM Post LIMIT %v OFFSET %v"
	const QUERY_CATEGORY = `SELECT id, name, date FROM POST
JOIN Category_Post ON Category_Post.post_id = Post.id
JOIN Category ON Category_Post.category_id = Category.id
WHERE Category.name = ? LIMIT %v OFFSET %v`

	const QUERY_COUNT = "SELECT count(*) FROM Post"
	var total int

	row := c.db.QueryRow(QUERY_COUNT)
	if row.Err() != nil {
		log.Fatalln(row.Err())
	}
	if err := row.Scan(&total); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("total: ", total)
	nPages := total / postsPerPage
	if total%postsPerPage > 0 {
		nPages++
	}
	fmt.Println("total pages: ", nPages)

	var query string
	if category == ALL {
		query = QUERY_ALL
	} else {
		query = QUERY_CATEGORY
	}

	//done := make(chan bool)
	batchCh := make(chan Batch)
	//postsDone := make(chan int)

	go func() {

		for i := 0; i < nPages; i++ {
			finalQuery := fmt.Sprintf(query, postsPerPage, i*postsPerPage)
			var extra string
			if category == ALL {
				extra = ""
			} else {
				extra = category
			}
			batch := Batch{c.db, finalQuery, extra, i + 1} // total, postsDone
			batchCh <- batch
		}
		close(batchCh)
	}()

	return batchCh //, done
}

func OpenContent(base string) Content {
	// check metadata.db
	// if not exists then create it
	metadataDB := filepath.Join(base, METADATA_FILE)
	posts := filepath.Join(base, "posts")
	pages := filepath.Join(base, "pages")
	var db *sql.DB
	var err error
	if db, err = sql.Open("sqlite3", metadataDB); err != nil {
		log.Fatalln(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	if _, err := db.Exec(databaseStruct); err != nil {
		log.Fatalln(err)
	}
	return Content{db, posts, pages}
}
