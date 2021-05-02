package mux

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type Mux struct {
	Routes []route
}

type route struct {
	Spec    string
	Pattern *regexp.Regexp
	Fn      func(w http.ResponseWriter, req *http.Request, params map[string]string)
	Doc     template.HTML
}

func New() *Mux {
	return &Mux{
		Routes: []route{},
	}
}

func (mux *Mux) HandleFunc(spec string, fn func(w http.ResponseWriter, req *http.Request, params map[string]string), doc string) {
	mux.Routes = append(mux.Routes, route{
		Spec:    spec,
		Pattern: regexp.MustCompile("^" + spec + "$"),
		Fn:      fn,
		Doc:     template.HTML(doc),
	})
}

func (mux Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Serving %s %s", req.Method, req.URL.String())

	for _, route := range mux.Routes {
		match := route.Pattern.FindStringSubmatch(req.URL.Path)
		if match != nil {
			details := make(map[string]string)
			names := route.Pattern.SubexpNames()
			for i, name := range names {
				if name != "" {
					details[name] = match[i]
				}
			}
			route.Fn(w, req, details)
			return
		}
	}

	http.NotFound(w, req)
}
