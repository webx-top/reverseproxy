package main

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/reverseproxy"
)

func main() {
	e := echo.New()
	e.Use(mw.Log())
	proxyOptions := &reverseproxy.ProxyOptions{
		Hosts:  []string{"https://localhost:8080/admin/"},
		Prefix: "/admin/",
		Engine: "fast",
	}
	e.Use(reverseproxy.Proxy(proxyOptions))
	e.Get("/", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Hello, World!")
	}))
	e.Get("/v2", echo.HandlerFunc(func(c echo.Context) error {
		return c.String("Echo v2")
	}))
	switch "fast" {
	case "fast":
		// FastHTTP
		e.Run(fasthttp.New(":4444"))
	default:
		// Standard
		e.Run(standard.New(":4444"))
	}

}
