package main

import (
	"fmt"
	"os"

	"github.com/yanodincov/json-schema-detector/internal/root"
)

func main() {
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
