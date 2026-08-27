package config

import (
    "bufio"
    "os"
    "strings"
)

// ReadLines reads lines from a specified text file.
// It automatically skips empty lines and comments (lines starting with #).
func ReadLines(filePath string) ([]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err // Return the error if the file doesn't exist
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        // Read the line and remove leading/trailing whitespaces
        line := strings.TrimSpace(scanner.Text())

        // Skip empty lines and lines that start with '#' (comments)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        lines = append(lines, line)
    }

    // Check if the scanner encountered any errors during reading
    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return lines, nil
}

// LoadEnv loads environment variables from a specified .env file into the process.
// It skips non-existent files gracefully, ignores comments (#) and empty lines,
// strips surrounding quotes, and does not overwrite existing environment variables.
func LoadEnv(filePath string) {
    lines, err := ReadLines(filePath)
    if err != nil {
        return // Missing file is silently ignored
    }

    for _, line := range lines {
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        val := strings.TrimSpace(parts[1])

        // Strip matching surrounding quotes if present
        if len(val) >= 2 {
            if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
                val = val[1 : len(val)-1]
            }
        }

        if key != "" && os.Getenv(key) == "" {
            _ = os.Setenv(key, val)
        }
    }
}
