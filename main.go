package urbandictionarylib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

var lastRequestTime *time.Time

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

func callUrbanDictionaryAPI(url string) (*UrbanDictionaryResponse, error) {
	shouldCallAPI := true

	if lastRequestTime == nil {
		lastRequestTime = new(time.Time)
		*lastRequestTime = time.Now()
	} else if time.Since(*lastRequestTime) < 1*time.Second {
		shouldCallAPI = false
	} else {
		*lastRequestTime = time.Now()
	}

	if !shouldCallAPI {
		return nil, fmt.Errorf("Don't call the API too fast! Wait at least 1 second between requests.")
	}

	resp, err := http.Get(url)
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

	// Handle 4xx and 5xx codes, allowing to add custom error handling for specific codes in the future
	switch sc := resp.StatusCode; {
	case sc >= 500:
		switch sc {
		default:
			return nil, fmt.Errorf("Server error: %d", sc)
		}
	case sc >= 400:
		switch sc {
		default:
			return nil, fmt.Errorf("Client error: %d", sc)
		}
	}

	var res UrbanDictionaryResponse
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func SearchTerm(term string) (*UrbanDictionaryResponse, error) {
	return callUrbanDictionaryAPI(UrbanDictionaryEndpoints["defineTerm"] + term)
}

func SearchTermId(id int) (*UrbanDictionaryResponse, error) {
	return callUrbanDictionaryAPI(UrbanDictionaryEndpoints["defineId"] + fmt.Sprint(id))
}

func SearchRandom() (*UrbanDictionaryResponse, error) {
	return callUrbanDictionaryAPI(UrbanDictionaryEndpoints["random"])
}
