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
    "log"
    "database/sql"
    //_ "github.com/mattn/go-sqlite3"
)

func (c Content) generateDefaultValues() {
    const QUERY = "INSERT INTO Option DEFAULT VALUES"
    if _, err := c.db.Exec(QUERY); err != nil {
        log.Fatalln(err)
    }
}

//SetPostsPerPage changes the posts por page when building a site
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
func (c Content) PostsPerPage() int64 {
    const QUERY = "SELECT posts_per_page FROM Option"
    row := c.db.QueryRow(QUERY)
    var ppp int64
    if err := row.Scan(&ppp); err == sql.ErrNoRows {
        c.generateDefaultValues()
        return c.PostsPerPage()
    } else if err != nil {
            log.Fatalln(err)
    }
    return ppp
}
