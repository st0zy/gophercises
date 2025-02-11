package main

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/st0zy/gophercises/transform/primitive"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html> <body>
		<form action=/upload method=post enctype="multipart/form-data">
		<input type="file" name="image">
			<button type="submit"> Upload image </button>
		</form>
		</body>
		</html>
		`
		fmt.Fprint(w, html)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer file.Close()
		img, err := genImage(file, filepath.Ext(header.Filename), 55, primitive.Beziers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		tpl := template.Must(template.New("").Parse(`<html><body>
		{{range .}}
			<img src="/{{.}}">
		{{end}}
		</body></html>`))
		tpl.Execute(w, []string{img})
	})

	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img", fs))
	http.ListenAndServe(":8000", mux)
}

func genImage(file multipart.File, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(file, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}
	img, _ := os.CreateTemp("./img", fmt.Sprintf("out_*%s", ext))
	defer img.Close()
	io.Copy(img, out)
	return img.Name(), nil
}
