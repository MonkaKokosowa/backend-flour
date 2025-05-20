package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog/log"
)

// NewProxy creates a new reverse proxy for a target URL.
func NewProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse target URL")
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Override the director if you need to modify the request further.
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// You can modify req here.
		// it's already in req.URL.Path.
	}

	return proxy
}
