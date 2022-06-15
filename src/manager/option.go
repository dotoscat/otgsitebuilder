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

package manager

import (
	"database/sql"
	"io/fs"
	"log"
)

func (c Content) generateDefaultValues() {
	const QUERY = "INSERT INTO Option DEFAULT VALUES"
	if _, err := c.db.Exec(QUERY); err != nil {
		log.Fatalln(err)
	}
}

//SetPostsPerPage changes the posts for page when building a site
func (c Content) SetPostsPerPage(ppp int) error {
	const QUERY = "UPDATE Option SET posts_per_page = ?"
	if result, err := c.db.Exec(QUERY, ppp); err != nil {
		return err
	} else if affected, _ := result.RowsAffected(); affected == 0 {
		c.generateDefaultValues()
		return c.SetPostsPerPage(ppp)
	}
	return nil
}

//PostsPerPage returns the number of posts por pages stored in metadata
func (c Content) PostsPerPage() int {
	const QUERY = "SELECT posts_per_page FROM Option"
	row := c.db.QueryRow(QUERY)
	var ppp int
	if err := row.Scan(&ppp); err == sql.ErrNoRows {
		c.generateDefaultValues()
		return c.PostsPerPage()
	} else if err != nil {
		log.Fatalln(err)
	}
	return ppp
}

//SetTitle changes the title of the website
func (c Content) SetTitle(title string) error {
	const QUERY = "UPDATE Option SET title = ?"
	if result, err := c.db.Exec(QUERY, title); err != nil {
		return err
	} else if affected, _ := result.RowsAffected(); affected == 0 {
		c.generateDefaultValues()
		return c.SetTitle(title)
	}
	return nil
}

//Title returns the title stored
func (c Content) Title() string {
	const QUERY = "SELECT title FROM Option"
	row := c.db.QueryRow(QUERY)
	var title string
	if err := row.Scan(&title); err == sql.ErrNoRows {
		c.generateDefaultValues()
		return c.Title()
	} else if err != nil {
		log.Fatalln(err)
	}
	return title
}

//SetLicense sets the general license of the website contents
func (c Content) SetLicense(license string) error {
	const QUERY = "UPDATE Option SET license = ?"
	if result, err := c.db.Exec(QUERY, license); err != nil {
		return err
	} else if affected, _ := result.RowsAffected(); affected == 0 {
		c.generateDefaultValues()
		return c.SetLicense(license)
	}
	return nil
}

//License returns the website license
func (c Content) License() string {
	const QUERY = "SELECT license FROM Option"
	row := c.db.QueryRow(QUERY)
	var license string
	if err := row.Scan(&license); err == sql.ErrNoRows {
		c.generateDefaultValues()
		return c.License()
	} else if err != nil {
		log.Fatalln(err)
	}
	return license
}

//SetOutput changes the title of the website
func (c Content) SetOutput(output string) error {
	if fs.ValidPath(output) == false {
		return ErrNoValid
	}
	const QUERY = "UPDATE Option SET output = ?"
	if result, err := c.db.Exec(QUERY, output); err != nil {
		return err
	} else if affected, _ := result.RowsAffected(); affected == 0 {
		c.generateDefaultValues()
		return c.SetOutput(output)
	}
	return nil
}

//Output returns the output path where the website is builded
func (c Content) Output() string {
	const QUERY = "SELECT output FROM Option"
	row := c.db.QueryRow(QUERY)
	if err := row.Scan(&output); err == sql.ErrNoRows {
		c.generateDefaultValues()
		return c.Output()
	} else if err != nil {
		log.Fatalln(err)
	}
	return ""
}
