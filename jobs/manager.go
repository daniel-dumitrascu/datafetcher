package jobs

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"datafetcher/profiles"
	"datafetcher/utils"
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

	mappings := make(map[string]profiles.MapData, 0)
	log.Printf("Loading mappings from %s\n", path)
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

		if len(mapdata.MapID) == 0 {
			// generate a new ID and add it to the mapdata
			mapdata.MapID = utils.GenerateUUID()

			rawData, err := json.MarshalIndent(mapdata, "", "  ")
			if err != nil {
				log.Printf("[ERROR] Error marshaling JSON: %s", err)
				continue
			}

			// write it back to the file
			err = os.WriteFile(filepath.Join(path, filenames[i]), rawData, 0644)
			if err != nil {
				log.Printf("[ERROR] Error writing JSON on file %s: %s", filenames[i], err)
				continue
			}
		}

		log.Printf("Map %s loaded sucesfully\n", filenames[i])
		mappings[mapdata.MapID] = mapdata
	}

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
