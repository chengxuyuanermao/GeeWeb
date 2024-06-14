package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	handlers map[string]HandleFunc
	roots    map[string]*node
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandleFunc),
		roots:    make(map[string]*node),
	}
}

func parsePattern(pattern string) []string {
	res := make([]string, 0)
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if part != "" {
			res = append(res, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return res
}

func (r *Router) AddRoute(method string, pattern string, handleFunc HandleFunc) {
	parts := parsePattern(pattern)

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(parts, pattern, 0)

	key := method + "-" + pattern
	r.handlers[key] = handleFunc
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root := r.roots[method]

	if root == nil {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			} else if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Router) Handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.String(http.StatusNotFound, "404 not found, method:%s, path:%s", c.Method, c.Path)
	}

	c.Next()
}
