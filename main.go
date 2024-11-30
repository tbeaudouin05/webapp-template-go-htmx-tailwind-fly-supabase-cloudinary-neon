package main

import (
	"context"
	"log"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/getNeonConnection"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/testNeonDatabase"
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

	ctx := context.Background()

	// Initialize Gin
	r := gin.Default()

	neonConnection := getNeonConnection.GetNeonConnection(ctx)
	defer neonConnection.Close(ctx)

	// Serve static files from the "static" directory
	r.Static("/"+goConstants.StaticFolder, "./"+goConstants.StaticFolder)

	// Serve the initial page
	r.GET("/test-neon-database", func(c *gin.Context) {
		testNeonDatabase.TestNeonDatabase(c.Request.Context(), neonConnection)
	})

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
