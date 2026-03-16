package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type PageData struct {
	Title       string
	Description string
	ActiveNav   string
	ContentTpl  string
}

var tmpl *template.Template

func init() {
	var err error
	tmpl = template.New("").Funcs(template.FuncMap{
		"templateExec": func(name string, data interface{}) (template.HTML, error) {
			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
				return "", err
			}
			return template.HTML(buf.String()), nil
		},
	})
	tmpl, err = tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
}

// compileSCSS compiles static/css/style.scss into static/css/style.css using
// the `sass` binary. If the output file already exists it skips compilation
// and returns immediately, preserving the cached artifact across restarts.
func compileSCSS() {
	const (
		src = "static/css/style.scss"
		dst = "static/css/style.css"
	)

	if _, err := os.Stat(dst); err == nil {
		log.Printf("scss: %s already exists, skipping compilation", dst)
		return
	}

	cmd := exec.Command("sass", "--no-source-map", src, dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("scss: compilation failed: %v", err)
	}

	log.Printf("scss: compiled %s -> %s", src, dst)
}

func main() {
	compileSCSS()

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("GET /{$}", handleHome)
	mux.HandleFunc("GET /about", handleAbout)
	mux.HandleFunc("GET /contact", handleContact)
	mux.HandleFunc("POST /contact", handleContactSubmit)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func renderPage(w http.ResponseWriter, r *http.Request, data PageData) {
	tplName := "base"
	if r.Header.Get("HX-Request") == "true" {
		tplName = data.ContentTpl
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, tplName, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderPage(w, r, PageData{
		Title:       "Home",
		Description: "Welcome to elegance in darkness",
		ActiveNav:   "home",
		ContentTpl:  "home_content",
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, PageData{
		Title:       "About",
		Description: "Architect of shadows",
		ActiveNav:   "about",
		ContentTpl:  "about_content",
	})
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, PageData{
		Title:       "Contact",
		Description: "Whispers in the dark",
		ActiveNav:   "contact",
		ContentTpl:  "contact_content",
	})
}

func handleContactSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	message := r.FormValue("message")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if name == "" || email == "" || message == "" {
		tmpl.ExecuteTemplate(w, "contact_form", map[string]string{
			"Error": "All fields must be filled",
			"Name":  name,
			"Email": email,
			"Msg":   message,
		})
		return
	}

	tmpl.ExecuteTemplate(w, "contact_success", map[string]string{
		"Name": name,
	})
}
