package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	e := gee.NewEngine()
	e.Get("/", func(c *gee.Context) {
		c.Json(http.StatusOK, "<h1>hello gee</h1>")
	})
	e.Get("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %v", c.Query("name"))
	})
	e.Post("/login", func(c *gee.Context) {
		c.Json(http.StatusOK, map[string]interface{}{
			"frame": "gee",
			"name":  c.PostForm("name"),
			"pwd":   c.PostForm("pwd"),
		})
	})

	e.Get("/hello/*filepath", func(c *gee.Context) {
		c.Json(http.StatusOK, gee.H{"name": c.GetParams("filepath")})
	})
	e.Get("/hello/:name", func(c *gee.Context) {
		c.Json(http.StatusOK, gee.H{"name": c.GetParams("name")})
	})

	e.Get("/assets/*filepath", func(c *gee.Context) {
		c.Json(http.StatusOK, gee.H{"filepath": c.GetParams("filepath")})
	})

	v1 := e.Group("/v1")
	v1.Use(gee.Logger())
	{
		v1.Get("/hello", func(c *gee.Context) {
			c.Json(http.StatusOK, "this is /v1/hello")
		})
		v1.Get("/hi", func(c *gee.Context) {
			c.Json(http.StatusOK, "this is /v1/hi")
		})
	}

	e.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	e.LoadHTMLGlob("templates/*")
	e.Static("/assets", "./static")

	stu1 := &student{Name: "geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	e.Get("/", func(c *gee.Context) {
		c.Html(http.StatusOK, "css.tmpl", nil)
	})
	e.Get("/students", func(c *gee.Context) {
		c.Html(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	e.Get("/date", func(c *gee.Context) {
		c.Html(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	e.Get("/panic", func(c *gee.Context) {
		name := []string{"1"}
		c.String(http.StatusOK, name[2])
	})

	log.Fatal(e.Run(":9999"))
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
