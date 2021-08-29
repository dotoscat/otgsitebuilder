package manager

import (
    "path/filepath"
    "log"
    "fmt"
    "time"
    "database/sql"
    "os"
    "errors"
    _ "embed"

    _ "github.com/mattn/go-sqlite3"
)

//go:embed database-struct.sql
var databaseStruct string

const (
    METADATA_FILE = ".metadata.db"
)

type Filer interface { //Fil(l)er
    Fill(*sql.Row, string) error
    Name() string
    Path() string
    Header() string
}

type File struct {
    id int64
    name string
    path string
}

func (f File) Name() string {
    return f.name
}

func (f File) Id() int64{
    return f.id
}

func (f File) Path() string {
    return f.path
}

func (f *File) SetPath(path string) {
    f.path = path
}

type Post struct {
    File
    date time.Time
}

func (p Post) Date() time.Time {
    return p.date
}

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

type Page struct {
    File
    reference string
}

func (p Page) Header() string {
    return ""
}

func (p *Page) Fill(row *sql.Row, basePath string) error {
    err := row.Scan(&p.id, &p.name, &p.reference)
    if err != nil {
        return err
    }
    p.path = filepath.Join(basePath, p.name)
    return err
}

type Content struct {
    db *sql.DB
    postsPath string
    pagesPath string
}

// *get content (post | page)*
// look at the file system: posts, pages
// if not found, exit
// select or created indexed content
// return content

var (
    ErrIsDir = errors.New("File is a directory.")
    ErrNotIndexed = errors.New("File is not indexed.")
)

func (c Content) CheckInPostsFolder(filename string) (bool, error) {
    postsFilePath := filepath.Join(c.postsPath, filename)
    if info, err := os.Stat(postsFilePath); err != nil {
        return false, err
    } else if info.IsDir() {
        return false, ErrIsDir
    } else {
        fmt.Println("info:", info)
    }
    return true, nil
}

func (c Content) GetPostMetadata(filename string) (Post, error) {
    const QUERY = "SELECT id, name, date FROM Post WHERE name = ?"
    row := c.db.QueryRow(QUERY, filename)
    post := Post{}
    err := post.Fill(row, c.postsPath)
    if err == sql.ErrNoRows {
        return post, ErrNotIndexed
    } else if err != nil {
        log.Fatalln(err)
    }
    return post, err
}

func (c Content) CreatePostMetadata(filename string) (int64, error) {
    const QUERY = "INSERT INTO Post (name) VALUES (?)"
    result, err := c.db.Exec(QUERY, filename)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    return id, err
}

func (c Content) GetPostFile(filename string) Post {
    if post, err := c.GetPostMetadata(filename); err != nil && err != ErrNotIndexed {
        log.Fatalln(err)
    } else if err == ErrNotIndexed {
        if _, err := c.CreatePostMetadata(filename); err != nil {
            log.Fatalln(err)
        } else {
            return c.GetPostFile(filename)
        }
    } else {
        return post
    }
    return Post{}
}

func (c Content) GetFile(filename string) interface{} {
    if isPost, err := c.CheckInPostsFolder(filename); err != nil {
        log.Fatalln(err)
    } else if isPost {
        return c.GetPostFile(filename)
    }
    return nil
    /*
    postsFilePath := filepath.Join(c.postsPath, filename)
    // var inFileSystem bool
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
        result, err := c.db.Exec(QUERY_INDEX_FILE, filename)
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
    */
}

func (c Content) GetPosts() []Post {
    // Index all files if they are not indexed
    files := make([]Post, 0)
    entries, err := os.ReadDir(c.postsPath)
    if err != nil {
        log.Fatalln(err)
    }
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        post := c.GetPostFile(entry.Name())
        files = append(files, post)
    }
    return files
}

func OpenContent(base string) Content {
    // check metadata.db
    // if not exists then create it
    metadataDB := filepath.Join(base, METADATA_FILE)
    posts := filepath.Join(base, "posts")
    pages := filepath.Join(base, "pages")
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
    return Content{db, posts, pages}
}

func ManageDatabase(base, filename string) {
    content := OpenContent(base)
    fmt.Println("content: ", content)
    contentFile := content.GetFile(filename)
    fmt.Println("file:", contentFile)
}
