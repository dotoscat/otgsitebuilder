package manager

import (
    "path/filepath"
    "log"
    "fmt"
    "time"
    "database/sql"
    "os"
    "errors"
    "strings"
    _ "embed"

    _ "github.com/mattn/go-sqlite3"
)

//go:embed database-struct.sql
var databaseStruct string

const (
    METADATA_FILE = ".metadata.db"
)

var (
    ErrIsDir = errors.New("File is a directory.")
    ErrNotIndexed = errors.New("File is not indexed.")
    ErrNoValid = errors.New("This file is not valid")
)

type Filer interface { //Fil(l)er
    Fill(*sql.Row, string) error
    FillFromRows(*sql.Rows, string) error
    Name() string
    Path() string
    Header() string
}

type File struct {
    id int64
    name string
    path string
    db *sql.DB
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

func (f *File) setDB(db *sql.DB) {
    f.db = db
}

func (f File) IsValid() bool {
    return f.db != nil
}

func (f *File) SetPath(path string) {
    f.path = path
}

type Post struct {
    File
    date time.Time
}

func newPost(db *sql.DB) Post {
    post := Post{}
    post.setDB(db)
    return post
}

func (p Post) Date() time.Time {
    return p.date
}

func (p Post) SetDate(date string) error {
    if !p.IsValid() {
        return ErrNoValid
    }
    const QUERY = "UPDATE POST SET date = ? WHERE id = ?"
    result, err := p.db.Exec(QUERY, date, p.id);
    n, err2 := result.RowsAffected()
    fmt.Println("err1:", err)
    fmt.Println("err2:", n, err2)
    if err != nil {
        log.Fatalln(err)
    }
    return err
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

func (p *Post) FillFromRows(rows *sql.Rows, basePath string) error {
    err := rows.Scan(&p.id, &p.name, &p.date)
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
    result, err := p.db.Exec(QUERY, reference, p.id);
    n, _ := result.RowsAffected()
    if err != nil {
        log.Fatalln(err)
    }
    if n > 0 {
        p.reference = reference
    }
    return err
}

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

type Content struct {
    db *sql.DB
    postsPath string
    pagesPath string
}

func checkInFolder(path string) (bool, error) {
    if info, err := os.Stat(path); err != nil {
        return false, err
    } else if info.IsDir() {
        return false, ErrIsDir
    }
    return true, nil
}

func (c Content) CheckInPagesFolder(filename string) (bool, error) {
    pagesFilePath := filepath.Join(c.pagesPath, filename)
    return checkInFolder(pagesFilePath)
}

func (c Content) CheckInPostsFolder(filename string) (bool, error) {
    postsFilePath := filepath.Join(c.postsPath, filename)
    return checkInFolder(postsFilePath)
}

func (c Content) getMetadata(recipient Filer, query string, values ...interface{}) error {
    row := c.db.QueryRow(query, values...)
    err := recipient.Fill(row, c.postsPath)
    if err == sql.ErrNoRows {
        return ErrNotIndexed
    } else if err != nil {
        log.Fatalln(err)
    }
    return err
}

func (c Content) GetPageMetadata(filename string) (Page, error) {
    const QUERY = "SELECT id, name, reference FROM Page WHERE name = ?"
    page := newPage(c.db)
    err := c.getMetadata(&page, QUERY, filename)
    return page, err
}

func (c Content) GetPostMetadata(filename string) (Post, error) {
    const QUERY = "SELECT id, name, date FROM Post WHERE name = ?"
    post := newPost(c.db)
    err := c.getMetadata(&post, QUERY, filename)
    fmt.Println("post metadata:", post.db != nil, "content db:", c.db != nil)
    return post, err
}

func (c Content) createMetadata(query string, values ...interface{}) (int64, error) {
    result, err := c.db.Exec(query, values...)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    return id, err
}

func (c Content) CreatePostMetadata(filename string) (int64, error) {
    const QUERY = "INSERT INTO Post (name) VALUES (?)"
    return c.createMetadata(QUERY, filename)
}

func (c Content) CreatePageMetadata(filename string) (int64, error) {
    const QUERY = "INSERT INTO Page (name) VALUES (?)"
    return c.createMetadata(QUERY, filename)
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

func (c Content) GetPageFile(filename string) Page {
    if page, err := c.GetPageMetadata(filename); err != nil && err != ErrNotIndexed {
        log.Fatalln(err)
    } else if err == ErrNotIndexed {
        if _, err := c.CreatePageMetadata(filename); err != nil {
            log.Fatalln(err)
        } else {
            return c.GetPageFile(filename)
        }
    } else {
        return page
    }
    return Page{}
}

func (c Content) GetPosts() []Post {
    // Index all files if they are not indexed

    entries, err := os.ReadDir(c.postsPath)
    if err != nil {
        log.Fatalln(err)
    }
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        c.GetPostFile(entry.Name()) // This function indexes if the file is not indexed, ignore the return value
    }
    const QUERY = "SELECT id, name, date FROM Post ORDER BY date DESC"
    rows, err := c.db.Query(QUERY)
    defer rows.Close()
    files := make([]Post, 0)
    for rows.Next() {
        post := newPost(c.db)
        post.FillFromRows(rows, c.postsPath)
        files = append(files, post)
    }
    return files
}

func (c Content) GetPages() []Page {
    // Index all files if they are not indexed

    entries, err := os.ReadDir(c.pagesPath)
    if err != nil {
        log.Fatalln(err)
    }
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        c.GetPageFile(entry.Name()) // This function indexes if the file is not indexed, ignore the return value
    }
    const QUERY = "SELECT id, name, reference FROM Page"
    rows, err := c.db.Query(QUERY)
    defer rows.Close()
    files := make([]Page, 0)
    for rows.Next() {
        page := newPage(c.db)
        page.FillFromRows(rows, c.pagesPath)
        files = append(files, page)
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
