package main

import (
	"log"
	"os"
	utilscmd "webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/cmd/utils_cmd"
)

func main() {
	// Change the working directory to the project root
	err := os.Chdir("../..") // Adjust the path if your structure is different
	if err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	utilscmd.RunCommands([]utilscmd.Command{
		{
			Name: "Generate Templ code",
			Args: []string{"templ", "generate"},
		},
		{
			Name: "Generate tailwind_mini.css",
			Args: []string{"npx", "tailwindcss", "-i", "tailwind.css", "-o", "static/tailwind_mini.css", "--minify"},
		},
		{
			Name: "Deploy app to Fly.io",
			Args: []string{"fly", "deploy"},
		},
	})

}
