package profiles

type MapData struct {
	Filename       string `json:"-"`
	MapID          string `json:"mapid"`
	ExternFile     string `json:"externfile"`
	Datatype       string `json:"datatype"`
	Fetchinterval  string `json:"fetchinterval"`
	LastUpdateDate string `json:"lastUpdateDate"`
	Map            []struct {
		Metafield string `json:"metafield"`
		Mapto     string `json:"mapto"`
		Update    bool   `json:"update"`
	} `json:"map"`
}
