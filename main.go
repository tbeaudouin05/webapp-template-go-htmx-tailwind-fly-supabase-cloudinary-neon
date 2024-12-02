package main

import (
	"context"
	"log"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/middlewares/passageAuthMiddleware"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/getNeonConnection"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/backend/neonDatabase/testNeonDatabase"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/frontend/components/helloComponent"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/frontend/pages/mainPage"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/frontend/pages/passageAuthPage"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goConstants"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goEnv"

	"github.com/gin-gonic/gin"
	"github.com/passageidentity/passage-go"
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

	// Initialize Passage
	psg, err := passage.New(goEnv.GlobalEnvVar.PassageAppId, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Unprotected route
	r.GET("/login", func(c *gin.Context) {
		c.Status(200)
		passageAuthPage.PassageAuthPage().Render(c.Request.Context(), c.Writer)
	})

	// Unprotected route
	r.GET("/test-neon-database", func(c *gin.Context) {
		testNeonDatabase.TestNeonDatabase(c.Request.Context(), neonConnection)
	})

	// Unprotected route: Home page
	r.GET("/", func(c *gin.Context) {
		c.Status(200)
		mainPage.Page().Render(c.Request.Context(), c.Writer)
	})

	// Unprrotected route: Handle HTMX request to update the greeting
	r.GET("/update", func(c *gin.Context) {
		c.Status(200)
		helloComponent.Hello("HTMX").Render(c.Request.Context(), c.Writer)
	})

	// Protected routes
	protected := r.Group("/")
	protected.Use(passageAuthMiddleware.PassageAuthMiddleware(psg))

	// Protected route: Dashboard
	protected.GET("/dashboard", func(c *gin.Context) {
		c.Status(200)
		mainPage.Page().Render(c.Request.Context(), c.Writer)
	})

	// Must run on 0.0.0.0, not localhost or anything else so that it works on fly.io
	r.Run("0.0.0.0:8080")
}
