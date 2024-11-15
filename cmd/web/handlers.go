package main

import (
	"database/sql"
	"errors"
	"fmt"
	"mango-monopoly/internal/models"
	"net/http"
	"strconv"
)

// home handler with a byte slice as the response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl", templateData{})
}

func (app *application) viewAllProperties(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, http.StatusOK, "properties.tmpl", templateData{})
}

func (app *application) propertyView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	property, err := app.properties.Get(int64(id))
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.render(w, r, http.StatusOK, "property.tmpl", templateData{
		Property:      *property,
		PropertyModel: &models.PropertyModel{},
	})
}

func (app *application) createPropertyPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "createProperty.tmpl", templateData{})
}

func (app *application) propertyCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	prop := models.Property{}

	prop.Address = sql.NullString{String: r.PostForm.Get("situs")}

	prop.City = sql.NullString{String: r.PostForm.Get("city")}

	prop.State = sql.NullString{String: r.PostForm.Get("state")}

	prop.Zip = sql.NullString{String: r.PostForm.Get("zip_code")}

	prop.County = sql.NullString{String: r.PostForm.Get("county_id")}

	prop.ParcelID = sql.NullString{String: r.PostForm.Get("parcel_id")}

	prop.PropertyType = sql.NullString{String: r.PostForm.Get("property_type")}

	landValue, err := strconv.ParseInt(r.PostForm.Get("land_value"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.LandValue = sql.NullInt64{Int64: landValue}

	buildingValue, err := strconv.ParseInt(r.PostForm.Get("building_value"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.BuildingValue = sql.NullInt64{Int64: buildingValue}

	fairMarketValue, err := strconv.ParseInt(r.PostForm.Get("fair_market_value"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.FairMarketValue = sql.NullInt64{Int64: fairMarketValue}

	lotSize, err := strconv.ParseFloat(r.PostForm.Get("lot_size"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.LotSize = sql.NullFloat64{Float64: lotSize}

	squareFootage, err := strconv.ParseInt(r.PostForm.Get("square_footage"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.SquareFt = sql.NullInt64{Int64: squareFootage}

	bedrooms, err := strconv.ParseInt(r.PostForm.Get("bedrooms"), 10, 16)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.Bedrooms = sql.NullInt16{Int16: int16(bedrooms)}

	bathrooms, err := strconv.ParseFloat(r.PostForm.Get("bathrooms"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.Bathrooms = sql.NullFloat64{Float64: bathrooms}

	yearBuilt, err := strconv.ParseInt(r.PostForm.Get("year_built"), 10, 16)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.YearBuilt = sql.NullInt16{Int16: int16(yearBuilt)}

	prop.TaxURL = sql.NullString{String: r.PostForm.Get("tax_assessor_url")}
	prop.ZillowURL = sql.NullString{String: r.PostForm.Get("zillow_url")}

	id, err := app.properties.Insert(&prop)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/property/view/%d", id), http.StatusSeeOther)
}
