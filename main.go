package main

import (
	"log"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/frontend/mainPage"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goConstants"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goEnv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := goEnv.GetEnvVar()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	r := gin.Default()

	// Serve static files from the "static" directory
	r.Static("/"+goConstants.StaticFolder, "./"+goConstants.StaticFolder)

	// Serve the initial page
	r.GET("/", func(c *gin.Context) {
		c.Status(200)
		mainPage.Page().Render(c.Request.Context(), c.Writer)
	})

	// Handle HTMX request to update the greeting
	r.GET("/update", func(c *gin.Context) {
		c.Status(200)
		mainPage.Hello("HTMX").Render(c.Request.Context(), c.Writer)
	})

	// must run on 0.0.0.0, not localhost or anything else so that it works on fly.io
	r.Run("0.0.0.0:8080")
}
