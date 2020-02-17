package engine

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/gocolly/colly/v2"
)

// Props : The scraping engine Properties and description about the engine (e.g NetNaijaEngine)
type Props struct {
	Name        string
	BaseURL     *url.URL // The Base URL for the engine
	SearchURL   *url.URL // URL for searching
	ListURL     *url.URL // URL to return movie lists
	Description string
}

// Engine : interface for all engines
type Engine interface {
	Search(query string) SearchResult
	Scrape(mode string) ([]Movie, error)
	List(page int) SearchResult
	String() string
}

// Movie : the structure of all downloadable movies
type Movie struct {
	Index          int
	Title          string
	CoverPhotoLink string
	Description    string
	Size           string
	DownloadLink   *url.URL
	Year           int
	IsSeries       bool
	SDownloadLink  []*url.URL // Other links for downloads if movies is series
	UploadDate     string
}

func (m *Movie) String() string {
	return fmt.Sprintf("%s (%s)", m.Title, m.Year)
}

// SearchResult : the results of search from engine
type SearchResult struct {
	Query  string
	Movies []Movie
}

// Titles : Get a slice of the titles of movies
func (s *SearchResult) Titles() []string {
	var titles []string
	for _, movie := range s.Movies {
		titles = append(titles, movie.Title)
	}
	return titles
}

// GetMovieByTitle : Return a movie object from title passed
func (s *SearchResult) GetMovieByTitle(title string) (Movie, error) {
	for _, movie := range s.Movies {
		if movie.Title == title {
			return movie, nil
		}
	}
	return Movie{}, errors.New("Movie not Found")
}

// GetEngines : Returns all the usable engines in the application
func GetEngines() map[string]Engine {
	engines := make(map[string]Engine)
	engines["NetNaija"] = NewNetNaijaEngine()
	return engines
}

// GetEngine : Return an engine
func GetEngine(engine string) Engine {
	return GetEngines()[engine]
}

// Get the movie index context stored in Request
func getMovieIndexFromCtx(r *colly.Request) int {
	movieIndex, err := strconv.Atoi(r.Ctx.Get("movieIndex"))
	if err != nil {
		log.Fatal(err)
	}
	return movieIndex
}
