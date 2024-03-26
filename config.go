package goPinotAPI

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

const controllerUrl = "controllerUrl"
const authToken = "authToken"
const authType = "authType"
const logger = "logger"

type Opt interface {
	apply(*cfg)
	Type() string
}

type cfg struct {
	controllerUrl  string
	authToken      string
	authType       string
	httpAuthWriter httpAuthWriter
	logger         *slog.Logger
}

type clientOpt struct{ fn func(*cfg) }

// controllerUrlOpt is an option to set the controller url for the client
type controllerUrlOpt struct {
	controllerUrl string
}

func (o *controllerUrlOpt) apply(c *cfg) {
	c.controllerUrl = o.controllerUrl
}

func (o *controllerUrlOpt) Type() string {
	return controllerUrl
}

// authTokenOpt is an option to set the auth token for the client
type authTokenOpt struct {
	authToken string
}

func (o *authTokenOpt) apply(c *cfg) {
	c.authToken = o.authToken
}

func (o *authTokenOpt) Type() string {
	return authToken
}

// authTypeOpt is an option to set the auth type for the client
type authTypeOpt struct {
	authType string
}

func (o *authTypeOpt) apply(c *cfg) {
	c.authType = o.authType
}

func (o *authTypeOpt) Type() string {
	return authType
}

// loggerOpt is an option to set the logger for the client
type loggerOpt struct {
	logger *slog.Logger
}

func (o *loggerOpt) apply(c *cfg) {
	c.logger = o.logger
}

func (o *loggerOpt) Type() string {
	return logger
}

func (opt clientOpt) apply(cfg *cfg) { opt.fn(cfg) }

func ControllerUrl(pinotControllerUrl string) Opt {
	return &controllerUrlOpt{controllerUrl: pinotControllerUrl}
}

func AuthToken(token string) Opt {
	return &authTokenOpt{authToken: token}
}

func Logger(logger *slog.Logger) Opt {
	return &loggerOpt{logger: logger}
}

func AuthType(authType string) Opt {
	return &authTypeOpt{authType: authType}
}

func validateOpts(opts ...Opt) (*cfg, *url.URL, error) {

	// with default auth writer that does nothing
	optCfg := defaultCfg()
	optCounts := make(map[string]int)
	for _, opt := range opts {

		switch opt.(type) {
		case *authTypeOpt:
			optCounts[authType]++
		case *authTokenOpt:
			optCounts[authToken]++
		case *controllerUrlOpt:
			optCounts[controllerUrl]++
		case *loggerOpt:
			optCounts[logger]++
		default:
			optCounts[opt.Type()]++
		}

		opt.apply(optCfg)

		if optCounts[authType] > 1 {
			return nil, nil, fmt.Errorf("multiple auth types provided")
		}
	}

	// validate controller url
	pinotControllerUrl, err := url.Parse(optCfg.controllerUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("controller url is invalid: %w", err)
	}
	// TODO: remove the redundant check
	// Currently this is designed to avoid a breaking change
	if optCfg.authType != "" && optCfg.authToken == "" {
		return nil, nil, fmt.Errorf("auth token is required when auth type is set")
	}
	// if auth token passed, handle authenticated requests
	if optCfg.authToken != "" {
		switch optCfg.authType {
		case "Bearer":
			optCfg.httpAuthWriter = func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", optCfg.authToken))
			}
		default:
			optCfg.httpAuthWriter = func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Basic %s", optCfg.authToken))
			}
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
