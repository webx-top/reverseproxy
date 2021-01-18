package main

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/reverseproxy"
)

func main() {
	e := echo.New()
	e.Use(mw.Log())
	proxyOptions := &reverseproxy.ProxyOptions{
		Hosts:  []string{"http://127.0.0.1:8084"},
		Prefix: "/",
		Engine: "fast",
	}
	e.Use(reverseproxy.Proxy(proxyOptions))
	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Hello, World!")
	}))
	e.Get("/v2", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Echo v2")
	}))

	c := &engine.Config{
		Address:     ":8084",
		TLSAuto:     false,
		TLSCertFile: "192.168.15.105.pem",
		TLSKeyFile:  "192.168.15.105-key.pem",
	}
	switch proxyOptions.Engine {
	case "FastHTTP", "fast":
		// FastHTTP
		e.Run(fasthttp.NewWithConfig(c))
	default:
		// Standard
		e.Run(standard.NewWithConfig(c))
	}

}
