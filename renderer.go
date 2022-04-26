package echo_jet_template_renderer

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/httpfs"
	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"reflect"
	"time"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templateFolder string
	jetviews       *jet.Set

	onRender func(*echo.Context, *map[string]interface{})
}

func NewTemplateRenderer(templateFolder string, fs http.FileSystem) TemplateRenderer {
	loader, _ := httpfs.NewLoader(fs)

	jetset := jet.NewSet(
		loader,
		jet.InDevelopmentMode(), // remove in production
	)

	jetset.AddGlobalFunc("humantime", func(a jet.Arguments) reflect.Value {
		a.RequireNumOfArguments("humantime", 1, 1)
		dt := a.Get(0).Interface().(time.Time)

		// Return a humanized version of the time object
		return reflect.ValueOf(humanize.Time(dt))
	})

	return TemplateRenderer{
		templateFolder: templateFolder,
		jetviews:       jetset,
	}
}

func (t *TemplateRenderer) SetRenderHook(onRender func(*echo.Context, *map[string]interface{})) {
	t.onRender = onRender
}

// Render renders a template document
func (t TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Ensure our data map is never nil
	var datamap map[string]interface{}
	if data == nil {
		datamap = make(map[string]interface{})
	} else {
		datamap = data.(map[string]interface{})
	}

	view, err := t.jetviews.GetTemplate(name)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Template rendering error: %s", err.Error()))
	}

	if t.onRender != nil {
		t.onRender(&c, &datamap)
	}

	return view.Execute(w, nil, datamap)
}
