package utils

import "os"

// Check if the environment variable exists, if not return the default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Populates the environment variables defined by a file.
// There's also a callback 'mapping' to apply business rules when needed.
func ExpandEnvWithDefault(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}
