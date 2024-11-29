package main

import (
	utilscmd "webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/cmd/utils_cmd"
)

func main() {

	utilscmd.RunCommands([]utilscmd.Command{
		{
			Name: "Generate tailwind_mini.css",
			Args: []string{"npx", "tailwindcss", "-i", "tailwind.css", "-o", "static/tailwind_mini.css", "--minify"},
		},
		{
			Name: "Run Go app",
			Args: []string{"go", "run", "."},
		},
	})
}
