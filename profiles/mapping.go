package profiles

type MapData struct {
	ID             int
	Filename       string
	File           string `json:"file"`
	Datatype       string `json:"datatype"`
	Fetchinterval  string `json:"fetchinterval"`
	LastUpdateDate string `json:"lastUpdateDate"`
	Map            []struct {
		Metafield string `json:"metafield"`
		Mapto     string `json:"mapto"`
		Update    bool   `json:"update"`
	} `json:"map"`
}
