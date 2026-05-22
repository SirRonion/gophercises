package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type PathToUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

type JSON struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := pathsToUrls[r.URL.Path]

		if url != "" {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYaml(yml []byte) ([]PathToUrl, error) {
	var parsed []PathToUrl
	err := yaml.Unmarshal(yml, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func buildMap(paths []PathToUrl) map[string]string {
	m := make(map[string]string)

	for _, path := range paths {
		m[path.Path] = path.Url
	}
	return m
}

func parseJson(data []byte) ([]JSON, error) {
	var paresedJson []JSON
	err := json.Unm
}
