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

type Content struct {
	db        *sql.DB
	postsPath string
	pagesPath string
}

func (c Content) Categories() Set {
	return newSet("Category", c.db)
}

func (c Content) CheckInPagesFolder(filename string) (bool, error) {
	pagesFilePath := filepath.Join(c.pagesPath, filename)
	return checkInFolder(pagesFilePath)
}

func (c Content) CheckInPostsFolder(filename string) (bool, error) {
	postsFilePath := filepath.Join(c.postsPath, filename)
	return checkInFolder(postsFilePath)
}

func (c Content) getMetadata(recipient Filer, query string, values ...interface{}) error {
	row := c.db.QueryRow(query, values...)
	err := recipient.Fill(row, c.postsPath)
	if err == sql.ErrNoRows {
		return ErrNotIndexed
	} else if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (c Content) GetPageMetadata(filename string) (Page, error) {
	const QUERY = "SELECT id, name, reference FROM Page WHERE name = ?"
	page := newPage(c.db)
	err := c.getMetadata(&page, QUERY, filename)
	return page, err
}

func (c Content) GetPostMetadata(filename string) (Post, error) {
	const QUERY = "SELECT id, name, date FROM Post WHERE name = ?"
	post := newPost(c.db)
	err := c.getMetadata(&post, QUERY, filename)
	fmt.Println("post metadata:", post.db != nil, "content db:", c.db != nil)
	return post, err
}

func (c Content) createMetadata(query string, values ...interface{}) (int64, error) {
	result, err := c.db.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

func (c Content) CreatePostMetadata(filename string) (int64, error) {
	const QUERY = "INSERT INTO Post (name) VALUES (?)"
	return c.createMetadata(QUERY, filename)
}

func (c Content) CreatePageMetadata(filename string) (int64, error) {
	const QUERY = "INSERT INTO Page (name) VALUES (?)"
	return c.createMetadata(QUERY, filename)
}

func (c Content) GetPostFile(filename string) Post {
	if post, err := c.GetPostMetadata(filename); err != nil && err != ErrNotIndexed {
		log.Fatalln(err)
	} else if err == ErrNotIndexed {
		if _, err := c.CreatePostMetadata(filename); err != nil {
			log.Fatalln(err)
		} else {
			return c.GetPostFile(filename)
		}
	} else {
		return post
	}
	return Post{}
}

func (c Content) GetPageFile(filename string) Page {
	if page, err := c.GetPageMetadata(filename); err != nil && err != ErrNotIndexed {
		log.Fatalln(err)
	} else if err == ErrNotIndexed {
		if _, err := c.CreatePageMetadata(filename); err != nil {
			log.Fatalln(err)
		} else {
			return c.GetPageFile(filename)
		}
	} else {
		return page
	}
	return Page{}
}

func (c Content) GetPosts() []Post {
	// Index all files if they are not indexed

	entries, err := os.ReadDir(c.postsPath)
	if err != nil {
		log.Fatalln(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		c.GetPostFile(entry.Name()) // This function indexes if the file is not indexed, ignore the return value
	}
	const QUERY = "SELECT id, name, date FROM Post ORDER BY date DESC"
	rows, err := c.db.Query(QUERY)
	defer rows.Close()
	files := make([]Post, 0)
	for rows.Next() {
		post := newPost(c.db)
		post.FillFromRows(rows, c.postsPath)
		files = append(files, post)
	}
	return files
}

func (c Content) GetPages() []Page {
	// Index all files if they are not indexed

	entries, err := os.ReadDir(c.pagesPath)
	if err != nil {
		log.Fatalln(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		c.GetPageFile(entry.Name()) // This function indexes if the file is not indexed, ignore the return value
	}
	const QUERY = "SELECT id, name, reference FROM Page"
	rows, err := c.db.Query(QUERY)
	defer rows.Close()
	files := make([]Page, 0)
	for rows.Next() {
		page := newPage(c.db)
		page.FillFromRows(rows, c.pagesPath)
		files = append(files, page)
	}
	return files
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
