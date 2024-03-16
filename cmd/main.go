package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/* -------------------------------- Templates ------------------------------- */
type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

/* --------------------------------- Domain --------------------------------- */
type Contact struct {
	Name  string
	Email string
}

func newContact(name string, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Contacts = []Contact

type PageData struct {
	Contacts Contacts
}

func newPageData() PageData {
	return PageData{
		Contacts: Contacts{
			newContact("John", "jd@email.com"),
			newContact("Cane", "ck@email.com"),
			newContact("test", "test"),
		},
	}
}

func (pd *PageData) hasEmail(email string) bool {
	for _, contact := range pd.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

type ContactFormData struct {
	Values map[string]string
	Errors map[string]string
}

func newContactFormData() ContactFormData {
	return ContactFormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	PageData        PageData
	ContactFormData ContactFormData
}

func newPage() Page {
	return Page{
		PageData:        newPageData(),
		ContactFormData: newContactFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()
	e.Static("/css", "css")

	// Variables
	page := newPage()

	// Endpoints
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusAccepted, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.PageData.hasEmail(email) {
			formData := newContactFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			return c.Render(http.StatusUnprocessableEntity, "create-contact-form", formData)
		}

		contact := newContact(name, email)
		page.PageData.Contacts = append(page.PageData.Contacts, contact)

		c.Render(http.StatusAccepted, "create-contact-form", newContactFormData())
		return c.Render(http.StatusAccepted, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start("localhost:42069"))
}
