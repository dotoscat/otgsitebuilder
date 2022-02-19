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
    // "database/sql"
    "time"
)

type File struct {
    id  int64
    name string
}

func (f File) Id() int64 {
    return f.id
}

func (f File) Name() string {
    return f.name
}

type Post struct{
    file File
    date    time.Time
}

func (p Post) File() File {
    return p.file
}

func (p Post) Date() time.Time {
    return p.date
}

type Page struct{
    file File
    reference   string
}

func (p Page) File() File {
    return p.file
}

func (p Page) Reference() string {
    return p.reference
}
