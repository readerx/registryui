package view

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"registryui/common"

	"github.com/CloudyKit/jet"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

// Template Jet template.
type Template struct {
	View *jet.Set
}

// Render render template.
func (r *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	t, err := r.View.GetTemplate(name)
	if err != nil {
		panic(fmt.Errorf("fatal error template file: %s", err))
	}
	vars, ok := data.(jet.VarMap)
	if !ok {
		vars = jet.VarMap{}
	}
	err = t.Execute(w, vars, nil)
	if err != nil {
		panic(fmt.Errorf("error rendering template %s: %s", name, err))
	}
	return nil
}

// SetupRender 初始化模板引擎
func SetupRender(debug bool, registryHost, basePath string) *Template {
	view := jet.NewHTMLSet("resources/templates")
	view.SetDevelopmentMode(debug)

	view.AddGlobal("version", common.Version)
	view.AddGlobal("base_path", basePath)
	view.AddGlobal("registry_host", registryHost)
	view.AddGlobal("pretty_size", func(size interface{}) string {
		var value float64
		switch i := size.(type) {
		case gjson.Result:
			value = float64(i.Int())
		case int64:
			value = float64(i)
		}
		return common.PrettySize(value)
	})
	view.AddGlobal("pretty_time", func(datetime interface{}) string {
		d := strings.Replace(datetime.(string), "T", " ", 1)
		d = strings.Replace(d, "Z", "", 1)
		return strings.Split(d, ".")[0]
	})
	view.AddGlobal("url_decode", func(m interface{}) string {
		res, err := url.PathUnescape(m.(string))
		if err != nil {
			return m.(string)
		}
		return res
	})
	view.AddGlobal("url_encode", func(m interface{}) string {
		res := url.PathEscape(m.(string))
		return res
	})

	return &Template{View: view}
}
