package main

import (
	utilscmd "webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/cmd/utils_cmd"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goConstants"
)

func main() {

	utilscmd.RunCommands([]utilscmd.Command{
		{
			Name: "Generate tailwind_mini.css",
			Args: []string{"npx", "tailwindcss", "-i", "tailwind.css", "-o", goConstants.StaticFolder + "/tailwind_mini.css", "--minify"},
		},
		{
			Name: "Run Go app",
			Args: []string{"go", "run", "."},
		},
	})
}
