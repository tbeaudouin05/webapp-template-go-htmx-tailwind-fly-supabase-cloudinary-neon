package utilscmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Name string   // Descriptive name for the command
	Args []string // Arguments for the command
}

// RunCommands takes a slice of Command structures and executes them sequentially.
func RunCommands(commands []Command) {
	for _, command := range commands {
		if len(command.Args) == 0 {
			continue
		}

		// Log the command being executed.
		log.Printf("%s: Executing command %s with args %v", command.Name, command.Args[0], command.Args[1:])

		// Create the command with the given arguments.
		cmd := exec.Command(command.Args[0], command.Args[1:]...)

		// Set the command to use the current environment's PATH.
		cmd.Env = os.Environ()

		// Redirect stdout and stderr to the standard output of this program.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Start the command.
		err := cmd.Start()
		if err != nil {
			log.Fatalf("%s: Error starting command: %v", command.Name, err)
		}

		// Wait for the command to finish.
		err = cmd.Wait()
		if err != nil {
			log.Fatalf("%s: Command finished with error: %v", command.Name, err)
		}
	}
}

func KillProcessOnPort(port string) error {
	// Find the PID using the port
	cmd := exec.Command("lsof", "-ti", "tcp:"+port)
	output, err := cmd.CombinedOutput()

	// Check if the error is because no process was found on the port
	if err != nil {
		if len(output) == 0 { // If there's no output, it means no process is using the port
			return nil // This is not an actual error for our purpose
		}
		return fmt.Errorf("failed to find process using port %s: %v", port, err)
	}

	// If output is empty, no process is using the port
	if len(output) == 0 {
		return nil
	}

	// Kill the process using the port
	pid := strings.TrimSpace(string(output))
	killCmd := exec.Command("kill", "-9", pid)
	if err := killCmd.Run(); err != nil {
		return fmt.Errorf("failed to kill process %s using port %s: %v", pid, port, err)
	}

	return nil
}
