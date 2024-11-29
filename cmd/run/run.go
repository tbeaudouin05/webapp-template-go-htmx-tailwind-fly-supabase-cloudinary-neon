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

	// Ensure port 8080 is available
	err = utilscmd.KillProcessOnPort("8080")
	if err != nil {
		log.Fatalf("Failed to free up port 8080: %v", err)
	}

	utilscmd.RunCommands([]utilscmd.Command{
		{Name: "Generate Templ code", Args: []string{"templ", "generate", "--watch", "--proxy=http://localhost:8080", "--cmd=go run ./cmd/watch_run/watch_run.go"}},
	})

}
