package goPinotAPI

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

type Opt interface {
	apply(*cfg)
}

type cfg struct {
	controllerUrl  string
	authToken      string
	httpAuthWriter httpAuthWriter
	logger         *slog.Logger
}

type clientOpt struct{ fn func(*cfg) }

func (opt clientOpt) apply(cfg *cfg) { opt.fn(cfg) }

func ControllerUrl(pinotControllerUrl string) Opt {
	return clientOpt{fn: func(cfg *cfg) { cfg.controllerUrl = pinotControllerUrl }}
}

func AuthToken(token string) Opt {
	return clientOpt{fn: func(cfg *cfg) { cfg.authToken = token }}
}

func Logger(logger *slog.Logger) Opt {
	return clientOpt{fn: func(cfg *cfg) { cfg.logger = logger }}
}

func validateOpts(opts ...Opt) (*cfg, *url.URL, error) {

	// with default auth writer that does nothing
	optCfg := defaultCfg()
	for _, opt := range opts {
		opt.apply(optCfg)
	}

	// validate controller url
	pinotControllerUrl, err := url.Parse(optCfg.controllerUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("controller url is invalid: %w", err)
	}

	// if auth token passed, handle authenticated requests
	if optCfg.authToken != "" {
		optCfg.httpAuthWriter = func(req *http.Request) {
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", optCfg.authToken))
		}
	}

	return optCfg, pinotControllerUrl, nil

}

func defaultCfg() *cfg {
	return &cfg{
		httpAuthWriter: defaultAuthWriter(),
		logger:         defaultLogger(),
	}
}

func defaultAuthWriter() func(*http.Request) {
	return func(req *http.Request) {
		// do nothing
	}
}

func defaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
