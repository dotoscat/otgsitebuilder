package manager

import (
    "path/filepath"
    "log"
    "fmt"
    "database/sql"
    _ "embed"

    _ "github.com/mattn/go-sqlite3"
)

//go:embed database-struct.sql
var databaseStruct string

const (
    METADATA_FILE = ".metadata.db"
)

type Content struct {
    db *sql.DB
}

func OpenContent(base string) Content {
    // check metadata.db
    // if not exists then create it
    metadataDB := filepath.Join(base, METADATA_FILE)
    // posts := filepath.Join(base, "posts")
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
    return Content{db}
}

func ManageDatabase(base string) {
    content := OpenContent(base)
    fmt.Println(content)
}
