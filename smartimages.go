package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/pat"
    "github.com/prasmussen/smartimages/config"
    "github.com/prasmussen/smartimages/handler"
    "github.com/prasmussen/smartimages/image"
    "github.com/prasmussen/smartimages/log"
)

func main() {
    // Load config file
    cfg, err := config.Load()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Instantiate logger 
    logger, err := log.New(cfg.LogFile)
    if err != nil {
        fmt.Println(err)
        return
    }

    pool := image.NewImagePool(cfg.ImageDir)
    handlers := handler.New(pool, logger)

    router := pat.New()
    router.Get("/images/{uuid}/file", handlers.GetImageFile())
    router.Get("/images/{uuid}", handlers.GetImage())
    router.Get("/images", handlers.ListImages())
    router.Delete("/images/{uuid}", handlers.DeleteImage())
    router.Post("/images/{uuid}", handlers.ImageAction())
    router.Post("/images", handlers.CreateImage())
    router.Put("/images/{uuid}/file", handlers.AddImageFile())
    router.Get("/ping", handlers.Ping())
    http.Handle("/", router)

    fmt.Printf("Listening for http connections on %s\n", cfg.Listen)
    if err := http.ListenAndServe(cfg.Listen, nil); err != nil {
       fmt.Println(err) 
    }
}
