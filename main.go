package urbandictionarylib

var UrbanDictionaryEndpoints = map[string]string{
	"defineTerm": "http://api.urbandictionary.com/v0/define?term=",
	"defineId":   "http://api.urbandictionary.com/v0/define?defid=",
	"random":     "http://api.urbandictionary.com/v0/random",
}

type UrbanDictionaryDefinition struct {
	Word        string `json:"word"`
	Definition  string `json:"definition"`
	Example     string `json:"example"`
	Permalink   string `json:"permalink"`
	Author      string `json:"author"`
	WrittenOn   string `json:"written_on"`
	ThumbsUp    int    `json:"thumbs_up"`
	ThumbsDown  int    `json:"thumbs_down"`
	DefId       int    `json:"defid"`
	CurrentVote string `json:"current_vote"` // No idea what this is
}

type UrbanDictionaryResponse struct {
	List []UrbanDictionaryDefinition `json:"list"`
}
