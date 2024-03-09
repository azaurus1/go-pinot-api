package goPinotAPI

type Opt interface {
	apply(*cfg)
}

type cfg struct {
	controllerUrl  string
	authToken      string
	httpAuthWriter httpAuthWriter
}

type controllerUrl struct {
	url string
}

type authToken struct {
	token string
}

func ControllerUrl(pinotControllerUrl string) Opt {
	return controllerUrl{url: pinotControllerUrl}
}

func AuthToken(token string) Opt {
	return authToken{token: token}
}

func (c controllerUrl) apply(cfg *cfg) {
	cfg.controllerUrl = c.url
}

func (a authToken) apply(cfg *cfg) {
	cfg.authToken = a.token
}
