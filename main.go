package main

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title       string
	Description string
	ActiveNav   string
	ContentTpl  string
}

// TemplateExec is a helper to execute templates by name
type TemplateExec struct {
	tmpl *template.Template
}

func (t *TemplateExec) Execute(name string, data interface{}) (template.HTML, error) {
	var buf bytes.Buffer
	err := t.tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func main() {
	r := gin.Default()

	// Must create FuncMap BEFORE parsing templates
	var templ *template.Template
	templ = template.New("").Funcs(template.FuncMap{
		"templateExec": func(name string, data interface{}) (template.HTML, error) {
			var buf bytes.Buffer
			err := templ.ExecuteTemplate(&buf, name, data)
			if err != nil {
				return "", err
			}
			return template.HTML(buf.String()), nil
		},
	})
	templ = template.Must(templ.ParseGlob("templates/*.html"))

	r.SetHTMLTemplate(templ)
	r.Static("/static", "./tmp")

	r.GET("/", handleHome)
	r.GET("/about", handleAbout)
	r.GET("/contact", handleContact)
	r.POST("/contact", handleContactSubmit)

	r.Run(":8080")
}

func handleHome(c *gin.Context) {
	data := PageData{
		Title:       "Home",
		Description: "Welcome to elegance in darkness",
		ActiveNav:   "home",
		ContentTpl:  "home_content",
	}
	renderPage(c, data)
}

func handleAbout(c *gin.Context) {
	data := PageData{
		Title:       "About",
		Description: "Architect of shadows",
		ActiveNav:   "about",
		ContentTpl:  "about_content",
	}
	renderPage(c, data)
}

func handleContact(c *gin.Context) {
	data := PageData{
		Title:       "Contact",
		Description: "Whispers in the dark",
		ActiveNav:   "contact",
		ContentTpl:  "contact_content",
	}
	renderPage(c, data)
}

func renderPage(c *gin.Context, data PageData) {
	if c.GetHeader("HX-Request") == "true" {
		// For HTMX, just return the content template
		c.HTML(http.StatusOK, data.ContentTpl, data)
		return
	}
	// For full page, render base which will include the content
	c.HTML(http.StatusOK, "base", data)
}

func handleContactSubmit(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	message := c.PostForm("message")

	if name == "" || email == "" || message == "" {
		c.HTML(http.StatusOK, "contact_form", gin.H{
			"Error": "All fields must be filled",
			"Name":  name,
			"Email": email,
			"Msg":   message,
		})
		return
	}

	c.HTML(http.StatusOK, "contact_success", gin.H{
		"Name": name,
	})
}
