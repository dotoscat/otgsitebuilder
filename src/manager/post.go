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
    "fmt"
    "time"
    "log"
    "path/filepath"
)

type Post struct {
	File
	date time.Time
}

func newPost(db *sql.DB) Post {
	post := Post{}
	post.setDB(db)
	return post
}

func (p Post) Date() time.Time {
	return p.date
}

func (p Post) SetDate(date string) error {
	if !p.IsValid() {
		return ErrNoValid
	}
	const QUERY = "UPDATE POST SET date = ? WHERE id = ?"
	result, err := p.db.Exec(QUERY, date, p.id)
	n, err2 := result.RowsAffected()
	fmt.Println("err1:", err)
	fmt.Println("err2:", n, err2)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// Deprecated: This is not a builder task
func (p Post) Header() string {
	return fmt.Sprint(p.Date())
}

func (p *Post) Fill(row *sql.Row, basePath string) error {
	err := row.Scan(&p.id, &p.name, &p.date)
	if err != nil {
		return err
	}
	p.SetPath(filepath.Join(basePath, p.name))
	return err
}

func (p *Post) FillFromRows(rows *sql.Rows, basePath string) error {
	err := rows.Scan(&p.id, &p.name, &p.date)
	if err != nil {
		return err
	}
	p.SetPath(filepath.Join(basePath, p.name))
	return err
}
