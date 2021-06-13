package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type yamlPath struct {
	path string `yaml:"path"`
	url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, 307)
			return
		} else {
			fallback.ServeHTTP(w, r)
			return
		}
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
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var ymlP []yamlPath
	err := yaml.Unmarshal(yml, &ymlP)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(ymlP)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(ymlP []yamlPath) map[string]string {
	ret := make(map[string]string)
	for _, uP := range ymlP {
		ret[uP.path] = uP.url
	}
	return ret
}
