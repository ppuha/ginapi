package ginapi

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	g *gin.Engine
}

func init() {
	caddy.RegisterModule(Handler{})
	httpcaddyfile.RegisterHandlerDirective("ginapi", parseCaddyfile)
}

func (h Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.ginapi",
		New: func() caddy.Module { return new(Handler) },
	}
}

func (h *Handler) Provision(ctx caddy.Context) error {
	router := gin.New()
	router.Use(logger())
	router.GET("/", handleRoot)
	h.g = router
	return nil
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	h.g.ServeHTTP(w, r)
	return nil
}

func (h *Handler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next()
	return nil
}

func parseCaddyfile(helper httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var h Handler
	err := h.UnmarshalCaddyfile(helper.Dispenser)
	return h, err
}

func logger() gin.HandlerFunc {
	logger := caddy.Log()

	return func(c *gin.Context) {
		logger.Info(
			"http request: ",
			zap.Object("request", caddyhttp.LoggableHTTPRequest{Request: c.Request}),
			zap.Object("headers", caddyhttp.LoggableHTTPHeader{Header: c.Writer.Header()}),
		)
	}
}
