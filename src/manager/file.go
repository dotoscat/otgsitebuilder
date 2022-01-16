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

type Filer interface { //Fil(l)er
	Fill(*sql.Row, string) error
	FillFromRows(*sql.Rows, string) error
	Id() int64
	Name() string
	Path() string
	Header() string
}

// File is the base struct for Post and Page
type File struct {
	id   int64
	name string
	path string
	db   *sql.DB
}

func (f File) Name() string {
	return f.name
}

func (f File) Id() int64 {
	return f.id
}

func (f File) Path() string {
	return f.path
}

func (f *File) setDB(db *sql.DB) {
	f.db = db
}

func (f File) IsValid() bool {
	return f.db != nil
}

func (f *File) SetPath(path string) {
	f.path = path
}
