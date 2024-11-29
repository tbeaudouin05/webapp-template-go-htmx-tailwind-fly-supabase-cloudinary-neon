package goEnv

import (
	"fmt"
	"os"
)

type ShouldUseCdn string

const (
	YesShouldUseCdn   ShouldUseCdn = "yes"
	NoShouldNotUseCdn ShouldUseCdn = "no"
)

// ParseShouldUseCdn validates if the environment variable is valid and returns the ShouldUseCdn value
func ParseShouldUseCdn() (ShouldUseCdn, error) {
	cdnVar := os.Getenv("ShouldUseCdn")
	switch ShouldUseCdn(cdnVar) {
	case YesShouldUseCdn, NoShouldNotUseCdn:
		return ShouldUseCdn(cdnVar), nil
	default:
		if cdnVar == "" {
			return "", fmt.Errorf("SHOULD_USE_CDN environment variable is not set")
		}
		return "", fmt.Errorf("invalid SHOULD_USE_CDN value: %s", cdnVar)
	}
}

type Env string

const (
	Production  Env = "production"
	Development Env = "development"
	Local       Env = "local"
)

// ParseEnv validates if the environment variable is valid and returns the Env value
func ParseEnv() (Env, error) {
	envVar := os.Getenv("Env")
	switch Env(envVar) {
	case Production, Development, Local:
		return Env(envVar), nil
	default:
		if envVar == "" {
			return "", fmt.Errorf("ENV environment variable is not set")
		}
		return "", fmt.Errorf("invalid ENV value: %s", envVar)
	}
}
