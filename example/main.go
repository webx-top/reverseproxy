package main

import (
	"github.com/webx-top/echo"
	// "github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/reverseproxy"
)

func main() {
	e := echo.New()
	e.Use(mw.Log())
	proxyOptions := &reverseproxy.ProxyOptions{
		Hosts:      []string{"localhost:8080"},
		PathPrefix: "/api/",
	}
	e.Use(reverseproxy.Proxy(proxyOptions))
	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	}))
	e.Get("/v2", echo.HandlerFunc(func(c echo.Context) error {
		return c.String(200, "Echo v2")
	}))

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// Standard
	e.Run(standard.New(":4444"))
}
