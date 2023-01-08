package urbandictionarylib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

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

func (res *UrbanDictionaryResponse) IsEmpty() bool {
	return len(res.List) == 0
}

func (res *UrbanDictionaryResponse) SortCustom(f func(a, b *UrbanDictionaryDefinition) bool) {
	sort.Slice(res.List, func(i, j int) bool {
		return f(&res.List[i], &res.List[j])
	})
}

func (res *UrbanDictionaryResponse) SortByThumbsUp() {
	res.SortCustom(func(a, b *UrbanDictionaryDefinition) bool {
		return a.ThumbsUp > b.ThumbsUp
	})
}

func (res *UrbanDictionaryResponse) SortByThumbsDown() {
	res.SortCustom(func(a, b *UrbanDictionaryDefinition) bool {
		return a.ThumbsDown > b.ThumbsDown
	})
}

func (res *UrbanDictionaryResponse) FilterByAuthor(author string) {
	sanitized := strings.ToLower(strings.TrimSpace(author))

	var filtered []UrbanDictionaryDefinition
	for _, def := range res.List {
		if strings.ToLower(strings.TrimSpace(def.Author)) == sanitized {
			filtered = append(filtered, def)
		}
	}

	res.List = filtered
}

func (res *UrbanDictionaryResponse) FilterMaxNDefinitions(n int) {
	if len(res.List) < n {
		return
	}

	res.List = res.List[:n]
}

func SearchTerm(term string) (*UrbanDictionaryResponse, error) {
	resp, err := http.Get(UrbanDictionaryEndpoints["defineTerm"] + term)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			fmt.Println("Error while trying to close response body")
			fmt.Println(closeErr)
		}
	}()

	var res UrbanDictionaryResponse
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
