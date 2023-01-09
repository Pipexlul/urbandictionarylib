package urbandictionarylib_test

import (
	"fmt"
	"testing"
	"time"

	udl "github.com/pipexlul/urbandictionarylib"
)

func TestInitial(t *testing.T) {
	resp, err := udl.SearchTerm("Chile")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Full unedited response:\n %+v\n\n", resp)

	resp.SortByThumbsUp()
	fmt.Printf("Sorted by thumbs up:\n %+v\n\n", resp)

	resp.SortByThumbsDown()
	fmt.Printf("Sorted by thumbs down:\n %+v\n\n", resp)

	resp.SortByThumbsUp()
	resp.FilterMaxNDefinitions(2)
	fmt.Printf("Sorted by thumbs up and filtered to 2 definitions:\n %+v\n\n", resp)

	if !resp.IsEmpty() {
		fmt.Printf("First definition:\n %+v\n\n", resp.List[0].Definition)
	}
}

func TestRateLimit(t *testing.T) {
	time.Sleep(1200 * time.Millisecond) // Wait at least 1.2 seconds to avoid rate limit

	udl.SearchTerm("Chile")
	_, err := udl.SearchTerm("LittleBigPlanet")

	if err == nil {
		fmt.Printf("Error: %v\n", err)
		t.Fatal("Expected rate limit error, got nil")
	}

	time.Sleep(1200 * time.Millisecond) // Wait at least 1.2 seconds to avoid rate limit

	resp, err := udl.SearchTerm("LittleBigPlanet")
	if err != nil {
		t.Fatal(err)
	}

	if !resp.IsEmpty() {
		fmt.Printf("Word: %v\n First definition:\n %+v\n\n", resp.List[0].Word, resp.List[0].Definition)
	}
}
