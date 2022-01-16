package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed assets/*
var rootFS embed.FS

//go:embed templates/*
var tplFS embed.FS

func runServer(w *WordStore) {
	fs, err := fs.Sub(rootFS, "assets")
	if err != nil {
		panic(err)
	}

	s := &WordleServer{
		w:        w,
		assetsFS: http.FS(fs),
	}
	r := gin.Default()

	templ := template.Must(template.New("").ParseFS(tplFS, "templates/*.html"))
	r.SetHTMLTemplate(templ)

	r.GET("/", s.indexHandler)
	r.POST("/query", s.queryHandler)

	r.StaticFS("/assets", s.assetsFS)

	r.Run(":3000")
}

type WordleServer struct {
	w        *WordStore
	assetsFS http.FileSystem
}

func (w *WordleServer) indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
func (w *WordleServer) queryHandler(c *gin.Context) {
	type request struct {
		Present []struct {
			Char     string `json:"char"`
			Position int    `json:"pos"`
		} `json:"present"`

		NotPresent []string `json:"notPresent"`
	}
	var req request

	c.BindJSON(&req)

	present := []posChar{}
	for _, pc := range req.Present {
		if len(pc.Char) > 0 {
			present = append(present, posChar{
				rune(pc.Char[0]),
				pc.Position - 1,
			})
		}
	}

	notPresent := []posChar{}
	for _, pc := range req.NotPresent {
		if len(pc) > 0 {
			notPresent = append(present, posChar{
				rune(pc[0]),
				-1,
			})
		}
	}
	query := &wordleQuery{
		present:    present,
		notPresent: notPresent,
	}
	data, _ := json.MarshalIndent(query, "", " ")
	fmt.Println(string(data))
	result := w.w.Execute(query)

	type response struct {
		Words []string `json:"words"`
	}
	c.JSON(http.StatusOK, response{
		Words: result,
	})
}
