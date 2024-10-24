package jobs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"datafetcher/profiles"
)

func StartJobs(path string) {

	//TODO create the ticker

	err := updateJobs(path)
	if err != nil {
		log.Printf("[ERROR] There was an issue when updating the jobs: %s\n", err)
		//TODO will retry at the next tick
	}
}

func updateJobs(path string) error {
	filenames, err := getFilesnames(path)
	if err != nil {
		return err
	}

	mappings := make([]profiles.MapData, 0)
	log.Println("Loading mappings...")
	for i := 0; i < len(filenames); i++ {
		rawjson, err := os.ReadFile(filepath.Join(path, filenames[i]))
		if err != nil {
			log.Printf("[ERROR] Could not read the map %s: %v\n", filenames[i], err)
			continue
		}

		var mapdata profiles.MapData
		err = json.Unmarshal(rawjson, &mapdata)
		if err != nil {
			log.Printf("[ERROR] Could not unmarshal the map %s: %v\n", filenames[i], err)
			continue
		}

		mappings = append(mappings, mapdata)
	}

	fmt.Println(mappings)

	return nil
}

func getFilesnames(dir string) ([]string, error) {
	root := os.DirFS(dir)

	jsonFiles, err := fs.Glob(root, "*.json")
	if err != nil {
		return nil, err
	}

	return jsonFiles, nil
}
