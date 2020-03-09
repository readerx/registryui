package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
)

type Registry struct {
	s  string
	r  client.Registry
	rt http.RoundTripper
}

func NewRegistry(server string) (*Registry, error) {
	rt := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	r, err := client.NewRegistry(server, rt)
	if err != nil {
		return nil, err
	}

	return &Registry{s: server, r: r, rt: rt}, nil
}

func (r *Registry) ListRepositories() ([]string, error) {
	repos := make([]string, 0)

	last := ""
	for end := false; !end; {
		entries := make([]string, 1024)

		n, err := r.r.Repositories(context.TODO(), entries, last)
		if err != nil {
			if err != io.EOF {
				return nil, err
			} else {
				end = true
			}
		}

		for i := 0; i < n; i++ {
			repos = append(repos, entries[i])
			last = entries[i]
		}
	}

	return repos, nil
}

func (r *Registry) ListTags(repo string) ([]string, error) {
	n, err := parseNamed(repo)
	if err != nil {
		return nil, err
	}

	repository, err := client.NewRepository(n, r.s, r.rt)
	if err != nil {
		return nil, err
	}

	tags, err := repository.Tags(context.TODO()).All(context.TODO())

	return tags, err
}

func (r *Registry) GetManifest(repo, tag string) (string, error) {
	n, err := parseNamed(repo)
	if err != nil {
		return "", err
	}

	repository, err := client.NewRepository(n, r.s, r.rt)
	if err != nil {
		return "", err
	}

	manifestService, err := repository.Manifests(context.TODO())
	if err != nil {
		return "", err
	}

	types := distribution.WithManifestMediaTypes([]string{schema2.MediaTypeManifest})
	manifest, err := manifestService.Get(context.TODO(), "", types, distribution.WithTag(tag))
	if err != nil {
		return "", err
	}

	_, data, err := manifest.Payload()
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err := json.Indent(&buff, data, "", "    "); err != nil {
		log.Printf("format manifest json error, %v\n", err)
		return string(data), err
	}

	return buff.String(), nil
}

func parseNamed(repo string) (reference.Named, error) {
	r, err := reference.Parse(repo)
	if err != nil {
		return nil, err
	}

	n, ok := r.(reference.Named)
	if !ok {
		return nil, fmt.Errorf("%s reference type %T\n", repo, r)
	}

	return n, nil
}
