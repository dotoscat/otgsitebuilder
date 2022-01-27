// Copyright 2021 Oscar Triano GarcÃ­a

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
	_ "embed"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed database-struct.sql
var databaseStruct string

const (
	METADATA_FILE = ".metadata.db"
)

var (
	ErrIsDir      = errors.New("File is a directory.")
	ErrNotIndexed = errors.New("File is not indexed.")
	ErrNoValid    = errors.New("This file is not valid")
)

func checkInFolder(path string) (bool, error) {
	if info, err := os.Stat(path); err != nil {
		return false, err
	} else if info.IsDir() {
		return false, ErrIsDir
	}
	return true, nil
}
