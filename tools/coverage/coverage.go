package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

func iteratePaths(swagger *spec.Swagger, f []byte) {
	// Convert the file content to a string
	fileContent := string(f)

	var foundPathCount int
	var totalPathCount int

	totalPathCount = len(swagger.Paths.Paths)

	// Iterate through swagger.Paths to get all API paths and operations
	for path, _ := range swagger.Paths.Paths {
		re := regexp.MustCompile(`\{.*?\}`)
		// Replace all matches with "%s"
		path = re.ReplaceAllString(path, "%s")
		// Check if the file content contains the path
		if strings.Contains(fileContent, path) {
			foundPathCount++
		} else {
			// fmt.Printf("Path not found in go-pinot-api.go: %s\n", path)
		}
	}

	// fmt.Printf("Number of paths found in go-pinot-api.go: %d\n", foundPathCount)
	fmt.Printf("Coverage: %f\n", float64(foundPathCount)/float64(totalPathCount))

}

func readGoPinotAPIs() []byte {

	filePath := filepath.Join("../../", "go-pinot-api.go")

	f, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return f

}

func main() {
	filePath := filepath.Join(".", "swagger.json")

	doc, err := loads.Spec(filePath)
	if err != nil {
		panic(err) // Handle error
	}
	swagger := doc.Spec()

	f := readGoPinotAPIs()

	iteratePaths(swagger, f)

}
