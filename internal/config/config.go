package config

import (
    "bufio"
    "os"
    "strings"
)

// ReadChannels reads the channel IDs from a specified text file.
// It automatically skips empty lines and comments (lines starting with #).
func ReadChannels(filePath string) ([]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err // Return the error if the file doesn't exist
    }
    defer file.Close()

    var channels []string
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        // Read the line and remove leading/trailing whitespaces
        line := strings.TrimSpace(scanner.Text())

        // Skip empty lines and lines that start with '#' (comments)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        channels = append(channels, line)
    }

    // Check if the scanner encountered any errors during reading
    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return channels, nil
}
