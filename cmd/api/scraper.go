package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"mango-monopoly/internal/scraper"
)

func (app *application) scrapeHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	county := params.ByName("county")
	if county == "" {
		http.Error(w, "missing county", http.StatusBadRequest)
		return
	}

	s, ok := scraper.Get(county)
	if !ok {
		http.Error(w, "unsupported county", http.StatusBadRequest)
		return
	}

	propID := params.ByName("propID")
	if propID == "" {
		http.Error(w, "missing property id", http.StatusBadRequest)
		return
	}

	property, err := s.ScrapeProperty(r.Context(), propID)
	if err != nil {
		app.logger.Error("scraper failed", "error", err, "county", county, "propID", propID)
		status := http.StatusInternalServerError
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusRequestTimeout
		}
		http.Error(w, "scrape failed", status)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, property, nil); err != nil {
		app.logger.Error("writing scrape response failed", "error", err, "county", county, "propID", propID)
	}
}
