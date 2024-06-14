package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandleFunc func(c *Context)

func NewEngine() *Engine {
	engine := &Engine{
		router: NewRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.Use(Logger(), Recovery())
	return engine
}

type Engine struct {
	*RouterGroup
	router        *Router
	groups        []*RouterGroup     // 存储所有的分组
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc
	parent      *RouterGroup
	engine      *Engine // 指向所属的engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		parent: r,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (r *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	pattern := r.prefix + comp
	r.engine.router.AddRoute(method, pattern, handler)
}

func (r *RouterGroup) Get(pattern string, handler HandleFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *RouterGroup) Post(pattern string, handler HandleFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *RouterGroup) Use(middlewares ...HandleFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.GetParams("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.Get(urlPattern, handler)
}

// -----

func (e *Engine) Get(route string, handleFunc HandleFunc) {
	e.router.AddRoute("GET", route, handleFunc)
}

func (e *Engine) Post(route string, handleFunc HandleFunc) {
	e.router.AddRoute("POST", route, handleFunc)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	matchHandlers := make([]HandleFunc, 0)
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			matchHandlers = append(matchHandlers, group.middlewares...)
		}
	}

	c := NewContext(w, req)
	c.handlers = matchHandlers
	c.engine = e
	e.router.Handle(c)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
