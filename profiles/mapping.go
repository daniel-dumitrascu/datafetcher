package profiles

type MapData struct {
	Filename       string `json:"-"`
	FileID         string `json:"fileid"`
	Datatype       string `json:"datatype"`
	Fetchinterval  string `json:"fetchinterval"`
	LastUpdateDate string `json:"lastUpdateDate"`
	Map            []struct {
		Metafield string `json:"metafield"`
		Mapto     string `json:"mapto"`
		Update    bool   `json:"update"`
	} `json:"map"`
}
