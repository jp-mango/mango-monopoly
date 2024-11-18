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
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) viewAllProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := app.properties.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Properties = properties

	app.render(w, r, http.StatusOK, "properties.tmpl", data)
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

	data := app.newTemplateData(r)
	data.Property = *property

	app.render(w, r, http.StatusOK, "property.tmpl", data)
}

func (app *application) createPropertyPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "createProperty.tmpl", data)
}

func (app *application) propertyCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	prop := models.Property{}

	prop.Address = sql.NullString{String: r.PostForm.Get("situs"), Valid: true}

	prop.City = sql.NullString{String: r.PostForm.Get("city"), Valid: true}

	prop.State = sql.NullString{String: r.PostForm.Get("state"), Valid: true}

	prop.Zip = sql.NullString{String: r.PostForm.Get("zip_code"), Valid: true}

	prop.County = sql.NullString{String: r.PostForm.Get("county_id"), Valid: true}

	prop.ParcelID = sql.NullString{String: r.PostForm.Get("parcel_id"), Valid: true}

	prop.PropertyType = sql.NullString{String: r.PostForm.Get("property_type"), Valid: true}

	landValue, err := strconv.ParseInt(r.PostForm.Get("land_value"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.LandValue = sql.NullInt64{Int64: landValue, Valid: true}

	buildingValue, err := strconv.ParseInt(r.PostForm.Get("building_value"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.BuildingValue = sql.NullInt64{Int64: buildingValue, Valid: true}

	fairMarketValue, err := strconv.ParseInt(r.PostForm.Get("fair_market_value"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.FairMarketValue = sql.NullInt64{Int64: fairMarketValue, Valid: true}

	lotSize, err := strconv.ParseFloat(r.PostForm.Get("lot_size"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.LotSize = sql.NullFloat64{Float64: lotSize, Valid: true}

	squareFootage, err := strconv.ParseInt(r.PostForm.Get("square_footage"), 10, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.SquareFt = sql.NullInt64{Int64: squareFootage, Valid: true}

	bedrooms, err := strconv.ParseInt(r.PostForm.Get("bedrooms"), 10, 16)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.Bedrooms = sql.NullInt16{Int16: int16(bedrooms), Valid: true}

	bathrooms, err := strconv.ParseFloat(r.PostForm.Get("bathrooms"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.Bathrooms = sql.NullFloat64{Float64: bathrooms, Valid: true}

	yearBuilt, err := strconv.ParseInt(r.PostForm.Get("year_built"), 10, 16)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	prop.YearBuilt = sql.NullInt16{Int16: int16(yearBuilt), Valid: true}

	prop.TaxURL = sql.NullString{String: r.PostForm.Get("tax_assessor_url"), Valid: true}
	prop.ZillowURL = sql.NullString{String: r.PostForm.Get("zillow_url"), Valid: true}

	prop.ID, err = app.properties.Insert(&prop)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Property successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/property/%d", prop.ID), http.StatusSeeOther)
}
