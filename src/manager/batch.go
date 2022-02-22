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
	"log"
	"time"
)

// Batch is a batch, page of posts, a posts ready to be read
type Batch struct {
	db         *sql.DB
	query      string
	queryExtra string
	i          int
	// postsDone chan int
	// totalPosts int
	// done chan bool
}

func (b Batch) Index() int {
	return b.i
}

func (b Batch) Posts() <-chan Post {
	postCh := make(chan Post)
	go func() {
		var (
			rows *sql.Rows
			err  error
		)

		rows, err = b.db.Query(b.query, b.queryExtra)

		defer rows.Close()

		if err != nil {
			log.Fatalln(err)
		}

		for rows.Next() {
			var (
				id   int64
				name string
				date time.Time
			)
			if err := rows.Scan(&id, &name, &date); err != nil {
				log.Fatalln(err)
			}
			post := Post{File{id, name}, date}
			postCh <- post
			//i := <-b.postsDone
			//i++
			//b.postsDone <- i
			//if i == b.totalPosts {
			//    b.done<- true
			//close(b.postsDone)
			//}
		}
		close(postCh)
	}()
	return postCh

}
