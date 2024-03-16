package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

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
	Id    int
}

var id = 0

func newContact(name string, email string) Contact {
	id++
	return Contact{
		Name:  name,
		Email: email,
		Id:    id,
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

func (pd *PageData) indexOf(id int) int {
	for i, contact := range pd.Contacts {
		if contact.Id == id {
			return i
		}
	}
	return -1
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
	e.Static("/images", "images")

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

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid id")
		}

		index := page.PageData.indexOf(id)
		if index == -1 {
			return c.String(http.StatusNotFound, "Contact not found")
		}

		page.PageData.Contacts = append(page.PageData.Contacts[:index], page.PageData.Contacts[index+1:]...)

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start("localhost:42069"))
}
