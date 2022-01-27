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
    "strings"
    "path/filepath"
    "log"
)

type Page struct {
	File
	reference string
}

func newPage(db *sql.DB) Page {
	page := Page{}
	page.setDB(db)
	return page
}

func (p *Page) SetReference(reference string) error {
	if !p.IsValid() {
		return ErrNoValid
	}
	const QUERY = "UPDATE Page SET reference = ? WHERE id = ?"
	result, err := p.db.Exec(QUERY, reference, p.id)
	n, _ := result.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	if n > 0 {
		p.reference = reference
	}
	return err
}

// Deprecated: Header must be implemented in builder
func (p Page) Header() string {
	if p.reference == "" {
		return strings.Replace(p.name, ".md", "", 1)
	}
	return p.reference
}

func (p *Page) Fill(row *sql.Row, basePath string) error {
	err := row.Scan(&p.id, &p.name, &p.reference)
	if err != nil {
		return err
	}
	p.path = filepath.Join(basePath, p.name)
	return err
}

func (p *Page) FillFromRows(rows *sql.Rows, basePath string) error {
	err := rows.Scan(&p.id, &p.name, &p.reference)
	if err != nil {
		return err
	}
	p.path = filepath.Join(basePath, p.name)
	return err
}
