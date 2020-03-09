package view

import (
	"log"
	"net/http"
	"net/url"

	"registryui/client"

	"github.com/CloudyKit/jet"
	"github.com/labstack/echo/v4"
)

type model struct {
	b string
	r *client.Registry
}

func NewModel(basepath string, r *client.Registry) *model {
	return &model{b: basepath, r: r}
}

func (m *model) Repositories(c echo.Context) error {
	repos, _ := m.r.ListRepositories()
	data := jet.VarMap{}
	data.Set("repos", repos)
	return c.Render(http.StatusOK, "repositories.html", data)
}

func (m *model) Tags(c echo.Context) error {
	r := c.Param("repo")

	repo, err := url.PathUnescape(r)
	if err != nil {
		return err
	}

	tags, err := m.r.ListTags(repo)
	if err != nil {
		return err
	}

	data := jet.VarMap{}
	data.Set("repo", repo)
	data.Set("tags", tags)

	return c.Render(http.StatusOK, "tags.html", data)
}

func (m *model) Manifest(c echo.Context) error {
	r := c.Param("repo")
	t := c.Param("tag")

	repo, err := url.PathUnescape(r)
	if err != nil {
		return err
	}

	tag, err := url.PathUnescape(t)
	if err != nil {
		return err
	}

	manifest, err := m.r.GetManifest(repo, tag)
	if err != nil {
		log.Printf("get manifest error %v\n", err)
		return err
	}

	data := jet.VarMap{}
	data.Set("repo", repo)
	data.Set("tag", tag)
	data.Set("manifest", manifest)

	return c.Render(http.StatusOK, "manifest.html", data)
}

func (m *model) Redrict(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, m.b)
}
