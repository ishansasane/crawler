package storage

import (
	"os"
)

func SaveToFile(filename string, data []string) {
    file, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    for _, line := range data {
        file.WriteString(line + "\n")
    }
}
