package views

import (
	"net/http"
	"path/filepath"
	"text/template"
)

type View struct {
	Folders Page
	Tasks   Page
}

type Page struct {
	Template *template.Template
	Layout   string
}

func NewView(layoutDir string) View {
	return View{
		Folders: NewPage(layoutDir, "Folders"),
		Tasks:   NewPage(layoutDir, "Tasks"),
	}
}

func NewPage(layoutDir string, layout string) Page {
	files, err := filepath.Glob(layoutDir + "/*.html")
	if err != nil {
		panic(err)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return Page{
		Template: t,
		Layout:   layout,
	}
}

func (p *Page) Render(w http.ResponseWriter, data any) error {
	return p.Template.ExecuteTemplate(w, p.Layout, data)
}
