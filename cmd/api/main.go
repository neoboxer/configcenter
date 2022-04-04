package main

import (
	"github.com/neoboxer/configcenter/internal/core"
	"net/http"
)

func main() {
	engine := core.NewEngine()
	engine.GET("/", func(ctx *core.Context) {
		ctx.JSON(http.StatusOK, core.M{})
	})

	engine.GET("/hello", func(ctx *core.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	engine.GET("/hello/*a/b", func(ctx *core.Context) {
		ctx.JSON(http.StatusOK, core.M{
			"path": ctx.Param("a"),
		})
	})
	engine.Run(":8888")
}
