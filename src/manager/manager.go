package manager

import (
    "path/filepath"
    "log"
    "fmt"
    "time"
    "database/sql"
    "os"
    _ "embed"

    _ "github.com/mattn/go-sqlite3"
)

//go:embed database-struct.sql
var databaseStruct string

const (
    METADATA_FILE = ".metadata.db"
)

const (
    TYPE_POST = iota + 1
    TYPE_PAGE
)

type File struct {
    id int64
    file string
    ftype int64
    date time.Time
    path string
}

func (f *File) Fill(row *sql.Row, basePath string) error {
    err := row.Scan(&f.id, &f.file, &f.date, &f.ftype)
    if err != nil {
        return err
    }
    f.path = filepath.Join(basePath, f.file)
    return err
}

func (f *File) Path() string {
    return f.path
}

type Content struct {
    db *sql.DB
    postsPath string
}

func (c Content) indexFile(filename string, ftype int64) {
    const QUERY_INDEX_FILE = "INSERT INTO CONTENT (file, contenttype_id) VALUES (?, ?)"
    c.db.Exec(QUERY_INDEX_FILE, filename, ftype)
}

func (c Content) GetFile(filename string) File {
    postsFilePath := filepath.Join(c.postsPath, filename)
    // var inFileSystem bool
    fileType := TYPE_POST
    if info, err := os.Stat(postsFilePath); err != nil {
        log.Fatalln(err)
    } else if info.IsDir() {
        log.Fatalln(filename, "is a directory.")
    } else {
        fmt.Println("info:", info)
    }
    // check if indexed
    fmt.Println("file:", filename)
    file := File{}
    const QUERY_FILE = "SELECT id, file, date, contenttype_id FROM Content WHERE file = ?"
    row := c.db.QueryRow(QUERY_FILE, filename)
    err := file.Fill(row, c.postsPath)
    fmt.Println("query:", QUERY_FILE, ";filename:", filename)
    fmt.Println("First fill error:", err)
    if err == sql.ErrNoRows {
        const QUERY_INDEX_FILE = "INSERT INTO CONTENT (file, contenttype_id) VALUES (?, ?)"
        result, err := c.db.Exec(QUERY_INDEX_FILE, filename, fileType)
        fmt.Println("result:", result, ";err:", err)
        if err != nil {
            log.Fatalln(err)
        } else {
            id, err := result.LastInsertId()
            fmt.Println("id:", id)
            if err != nil {
                log.Fatalln(err)
            }
            const QUERY_INDEX_FILE_ID = "SELECT id, file, date, contenttype_id FROM Content WHERE id = ?"
            row := c.db.QueryRow(QUERY_INDEX_FILE_ID, id)
            if err := file.Fill(row, c.postsPath); err != nil { // I hope not
                log.Fatalln(err)
            }
        }
    } else if err != nil {
        log.Fatalln(err)
    }
    return file
}

func (c Content) GetPosts() []File {
    // Index all files if they are not indexed
    files := make([]File, 0)
    entries, err := os.ReadDir(c.postsPath)
    if err != nil {
        log.Fatalln(err)
    }
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        file := c.GetFile(entry.Name())
        files = append(files, file)
    }
    return files
}

func OpenContent(base string) Content {
    // check metadata.db
    // if not exists then create it
    metadataDB := filepath.Join(base, METADATA_FILE)
    posts := filepath.Join(base, "posts")
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
    return Content{db, posts}
}

func ManageDatabase(base, filename string) {
    content := OpenContent(base)
    fmt.Println("content: ", content)
    contentFile := content.GetFile(filename)
    fmt.Println("file:", contentFile)
}
