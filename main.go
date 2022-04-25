package main

import (
	"log"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

type Project struct {
	Volume int `json:"volume,string"`
}

func main() {
	log.Println("starting...")

	iris.RegisterOnInterrupt(func() {
		// TODO
	})

	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app.Use(crs)

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.Writef("404 not found here")
	})

	app.Post("/save", func(ctx iris.Context) {
		var req Project
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}
		log.Println(req.Volume)

		ctx.JSON(iris.Map{
			"success": true,
		})
	})

	app.HandleDir("/", iris.Dir("./public"))

	err := app.Run(
		// Start the web server at localhost:8080
		iris.Addr(":5000"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)

	if err != nil {
		log.Println(err.Error())
	}
}
