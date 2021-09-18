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
	"fmt"
	"log"
)

type Set struct {
	name string
	db   *sql.DB
}

func newSet(name string, db *sql.DB) Set {
	return Set{name, db}
}

func (s Set) AddElement(element string) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %v (name) VALUES (?)", s.name)
	result, err := s.db.Exec(query, element)
	id, _ := result.LastInsertId()
	return id, err
}

func (s Set) AddPostForElement(post Post, element string) error {
	// 1. Look if element exists in the table
	// 2. If not, create it
	// 3. Get element id
	// 4. Insert in many to many table element id and post id
	var elementId int64
	queryElement := fmt.Sprintf("SELECT id FROM %v WHERE name = ?", s.name)
	fmt.Println("Query: ", queryElement, "; element:", element)
	row := s.db.QueryRow(queryElement, element)
	if err := row.Scan(&elementId); err == sql.ErrNoRows {
		fmt.Println("No rows")
		if elementId, err = s.AddElement(element); err != nil {
			log.Fatalln(err)
		}
	} else if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("element: ", element, "; id: ", elementId)
	INSERT_NAME_POST := fmt.Sprintf("INSERT INTO %v_Post (%v_id, post_id) VALUES (?, ?)", s.name, s.name)
	_, err := s.db.Exec(INSERT_NAME_POST, elementId, post.Id())
	return err
	// insert as many to many
	//query := fmt.Sprintf("")
}

func (s Set) RemovePostForElement(post Post, element string) error {
	QUERY := fmt.Sprintf("DELETE FROM %v_POST WHERE post_id = ? AND %v_id = (SELECT id FROM %v WHERE name = ?)", s.name, s.name, s.name)
	_, err := s.db.Exec(QUERY, post.Id(), element)
	return err
}
