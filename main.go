package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

type Project struct {
	Volume int `json:"volume,string"`
}

var (
	file *os.File
)

func main() {
	log.Println("starting...")
	defer func() {
		if file != nil {
			file.Close()
			log.Println("closing opened file")
		} else {
			log.Println("no file need to close")
		}
	}()

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

	app.Post("/open", func(ctx iris.Context) {
		var err error
		if file != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}
		file, err = os.OpenFile("project.json", os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
		} else {
			log.Println("opened file")
			ctx.JSON(iris.Map{
				"success": true,
			})
		}
	})

	app.Post("/save", func(ctx iris.Context) {
		if file == nil {
			ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("no file opened"))
			return
		}
		var req Project
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}

		bytes, _ := json.Marshal(req)

		file.Truncate(0)
		file.Seek(0, 0)

		w := bufio.NewWriter(file)
		w.Write(bytes)

		w.Flush()
		log.Println("saved file")

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
