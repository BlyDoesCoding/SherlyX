package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var globalMap = make(map[string][]string)
var mapMutex sync.RWMutex

func dostuff(fileName string) {

	key, _ := calculateFileHash(fileName)

	addData(key, fileName)

}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for fileName := range jobs {
		dostuff(fileName)
	}
}

func start(paths []string, container *fyne.Container) {

	files := listFilesInFolders(paths) // Replace with your list of files

	numWorkers := runtime.NumCPU()

	// Set the number of concurrent workers

	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	for _, file := range files {
		jobs <- file
	}

	close(jobs)
	wg.Wait()
	removeKeysWithOneValue()

	mapMutex.RLock()
	defer mapMutex.RUnlock()

	groups := 0
	duplicateFiles := 0

	for key, files := range globalMap {
		string := "MD5 " + key + ": "
		groups = groups + 1
		fmt.Println("Associated Files:")
		for _, file := range files {
			string = string + file + " "
			duplicateFiles = duplicateFiles + 1
		}
		println(string)

	}
	result := "Duplicate Groups found: " + fmt.Sprint(groups)
	dups := "Duplicate Files found: " + fmt.Sprint(duplicateFiles)

	container.Add(widget.NewLabel(result))
	container.Add(widget.NewLabel(dups))

	container.Add(widget.NewButton("Clear Duplicates", deleteFilesExceptFirst))

}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Convert the MD5 hash to a hexadecimal string
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, err
}

func addData(key, value string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	// Retrieve the existing slice of strings for the key
	existingValues, exists := globalMap[key]

	// If the key doesn't exist, create a new entry with a slice containing the value
	if !exists {
		globalMap[key] = []string{value}
	} else {
		// If the key already exists, append the value to the existing slice
		globalMap[key] = append(existingValues, value)
	}
}

func getData(key string) []string {
	mapMutex.RLock()
	defer mapMutex.RUnlock()
	return globalMap[key]
}

func listFilesInFolders(folderPaths []string) []string {
	var files []string

	// Iterate through each folder path
	for _, folderPath := range folderPaths {
		_ = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
	}

	return files
}
func printMapContents() {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	fmt.Println("Global Map Contents:")
	for key, files := range globalMap {
		fmt.Printf("Key: %s\n", key)
		fmt.Println("Associated Files:")
		for _, file := range files {
			fmt.Printf("  %s\n", file)
		}
	}
}

func removeKeysWithOneValue() {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	for key, values := range globalMap {
		if len(values) == 1 {
			delete(globalMap, key)
		}
	}
}

func uriSliceToFilePaths(uris []fyne.URI) []string {
	filePaths := make([]string, 0, len(uris))

	for _, uri := range uris {
		uriString := uri.String()
		// Check if the URI starts with "file://"
		if strings.HasPrefix(uriString, "file://") {
			// Remove the "file://" prefix
			filePath := uriString[7:]
			filePaths = append(filePaths, filePath)
		} else {
			// Not a file URI, simply add the URI string
			filePaths = append(filePaths, uriString)
		}
	}

	return filePaths
}

func stringArrayToString(stringSlice []string) string {
	// Use strings.Join to concatenate the strings in the slice
	resultString := strings.Join(stringSlice, "")

	return resultString
}

func deleteFilesExceptFirst() {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	for key, filepaths := range globalMap {
		if len(filepaths) <= 1 {
			continue // Skip if there's only one or no file to delete
		}

		// Keep the first file and delete the rest
		firstFilePath := filepaths[0]
		for _, filePath := range filepaths[1:] {
			err := os.Remove(filePath)
			if err != nil {
			}
		}

		// Update the global map with only the first file path
		globalMap[key] = []string{firstFilePath}
	}
}
