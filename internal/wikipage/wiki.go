package wikipage

import (
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	f := p.Title + ".page"
	return ioutil.WriteFile(f, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	f := title + ".page"
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: data}, nil
}
