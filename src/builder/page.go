package builder

//Page stores info about the pages which it belongs, index,
//a list of writings related to this page and its url.
type Page struct {
    parent *Pages
    index int
    writings []Writing
    url string
}

//Pages is defined type for a slice of Page.
type Pages []Page

//HasLast returns true if this is not the last page.
func (p Page) HasLast() bool {
    return p.index - 1 >= 0
}

//HasNext returns true if there are more pages.
func (p Page) HasNext() bool {
    return p.index + 1 < len(*p.parent)
}

//Last returns the last page on this page.
//Returns an empty Page if not.
func (p Page) Last() Page {
    if p.HasLast() {
        return (*p.parent)[p.index-1]
    }
    return Page{}
}

//Next returns the next page on this page.
//Returns an empty Page if is not possible.
func (p Page) Next() Page {
    if p.HasNext() {
        return (*p.parent)[p.index+1]
    }
    return Page{}
}

//Url returns the final url of this page.
func (p Page) Url() string {
    return p.url
}

//Empty tells if this Page is empty, does not belong to pages.
func (p Page) Empty() bool {
    return p.parent == nil
}

//Writinngs returns the writings stored in this page
func (p Page) Writings() []Writing {
    return p.writings
}

//addWriting is a convenient method to add writings to this page.
//The return value is the same writing passed as argument.
func (p *Page) addWriting(writing Writing) Writing {
    p.writings = append(p.writings, writing)
    return writing
}
