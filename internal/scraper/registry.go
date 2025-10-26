package scraper

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"mango-monopoly/internal/data"
)

// PropertyScraper defines the behaviour for extracting a Property by parcel id.
type PropertyScraper interface {
	ScrapeProperty(ctx context.Context, propID string) (*data.Property, error)
}

var (
	registryMu sync.RWMutex
	registry   = map[string]PropertyScraper{}
)

// Register makes a scraper available for lookup by county identifier.
func Register(county string, s PropertyScraper) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[strings.ToLower(county)] = s
}

// Get returns a registered scraper by county identifier.
func Get(county string) (PropertyScraper, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	s, ok := registry[strings.ToLower(county)]
	return s, ok
}

// MustRegister registers the scraper and panics if county already exists.
func MustRegister(county string, s PropertyScraper) {
	registryMu.Lock()
	defer registryMu.Unlock()

	key := strings.ToLower(county)
	if _, exists := registry[key]; exists {
		panic(fmt.Sprintf("scraper for county %q already registered", county))
	}
	registry[key] = s
}
