// go_packages/go_env/env.go
package goEnv

import (
	"log"

	"github.com/joho/godotenv"
)

// EnvVar holds all environment variables required by the application.
type EnvVar struct {
	Env          Env
	ShouldUseCdn ShouldUseCdn
}

// GlobalEnvVar is the global instance holding all environment variables.
var GlobalEnvVar EnvVar

// GetEnvVar loads environment variables from the .env file and environment,
// parses them, assigns them to GlobalEnvVar, and initializes Stripe.
func GetEnvVar() error {
	// Try to load the .env file. If it fails, log the error but do not exit.
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Parse enums
	envType, err := ParseEnv()
	if err != nil {
		return err
	}

	shouldUseCdn, err := ParseShouldUseCdn()
	if err != nil {
		return err
	}

	// Retrieve environment variables
	envVars := EnvVar{
		Env:          envType,
		ShouldUseCdn: shouldUseCdn,
	}

	GlobalEnvVar = envVars

	return nil
}
