package main

import (
	"database/sql"
	"errors"
	"fmt"
	"mango-monopoly/internal/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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
	data.Form = propertyCreateForm{
		FieldErrors: make(map[string]string),
	}
	app.render(w, r, http.StatusOK, "createProperty.tmpl", data)
}

type propertyCreateForm struct {
	Address         string
	City            string
	State           string
	Zip             string
	County          string
	ParcelID        string
	PropertyType    string
	LandValue       string
	BuildingValue   string
	FairMarketValue string
	LotSize         string
	SquareFootage   string
	Bedrooms        string
	Bathrooms       string
	YearBuilt       string
	TaxURL          string
	ZillowURL       string
	FieldErrors     map[string]string
}

func (app *application) propertyCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := propertyCreateForm{
		Address:         r.PostForm.Get("address"),
		City:            r.PostForm.Get("city"),
		State:           r.PostForm.Get("state"),
		Zip:             r.PostForm.Get("zip_code"),
		County:          r.PostForm.Get("county_id"),
		ParcelID:        r.PostForm.Get("parcel_id"),
		PropertyType:    r.PostForm.Get("property_type"),
		LandValue:       r.PostForm.Get("land_value"),
		BuildingValue:   r.PostForm.Get("building_value"),
		FairMarketValue: r.PostForm.Get("fair_market_value"),
		LotSize:         r.PostForm.Get("lot_size"),
		SquareFootage:   r.PostForm.Get("square_footage"),
		Bedrooms:        r.PostForm.Get("bedrooms"),
		Bathrooms:       r.PostForm.Get("bathrooms"),
		YearBuilt:       r.PostForm.Get("year_built"),
		TaxURL:          r.PostForm.Get("tax_assessor_url"),
		ZillowURL:       r.PostForm.Get("zillow_url"),
		FieldErrors:     make(map[string]string),
	}

	prop := models.Property{}

	// Address
	prop.Address = sql.NullString{String: r.PostForm.Get("address"), Valid: true}
	if prop.Address.String == "" {
		form.FieldErrors["address"] = "This field cannot be blank"
	}

	// City
	prop.City = sql.NullString{String: r.PostForm.Get("city"), Valid: true}
	if prop.City.String == "" {
		form.FieldErrors["city"] = "This field cannot be blank"
	}

	// State
	prop.State = sql.NullString{String: r.PostForm.Get("state"), Valid: true}
	if prop.State.String == "" {
		form.FieldErrors["state"] = "This field cannot be blank"
	}

	// Zip Code
	prop.Zip = sql.NullString{String: r.PostForm.Get("zip_code"), Valid: true}
	if utf8.RuneCountInString(prop.Zip.String) > 7 {
		form.FieldErrors["zip_code"] = "This field cannot be more than 7 characters long"
	} else if prop.State.String == "" {
		form.FieldErrors["zip_code"] = "This field cannot be blank"
	}

	// County
	prop.County = sql.NullString{String: r.PostForm.Get("county_id"), Valid: true}

	// Parcel ID
	prop.ParcelID = sql.NullString{String: r.PostForm.Get("parcel_id"), Valid: true}
	if strings.TrimSpace(prop.ParcelID.String) == "" {
		form.FieldErrors["parcel_id"] = "This field cannot be blank"
	} else if _, err := app.properties.GetByParcel(prop.ParcelID.String); err != nil {
		form.FieldErrors["parcel_id"] = "This property already exists"
	}

	// Property Type
	prop.PropertyType = sql.NullString{String: r.PostForm.Get("property_type"), Valid: true}

	// Numeric fields
	if landValue, err := strconv.ParseInt(r.PostForm.Get("land_value"), 10, 64); err == nil {
		if landValue >= 0 {
			prop.LandValue = sql.NullInt64{Int64: landValue, Valid: true}
		} else {
			form.FieldErrors["land_value"] = "Land value must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["land_value"] = "Invalid value for land value"
	}

	if buildingValue, err := strconv.ParseInt(r.PostForm.Get("building_value"), 10, 64); err == nil {
		if buildingValue >= 0 {
			prop.BuildingValue = sql.NullInt64{Int64: buildingValue, Valid: true}
		} else {
			form.FieldErrors["building_value"] = "Building value must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["building_value"] = "Invalid value for building value"
	}

	if fairMarketValue, err := strconv.ParseInt(r.PostForm.Get("fair_market_value"), 10, 64); err == nil {
		if fairMarketValue >= 0 {
			prop.FairMarketValue = sql.NullInt64{Int64: fairMarketValue, Valid: true}
		} else {
			form.FieldErrors["fair_market_value"] = "Fair market value must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["fair_market_value"] = "Invalid value for fair market value"
	}

	if lotSize, err := strconv.ParseFloat(r.PostForm.Get("lot_size"), 64); err == nil {
		if lotSize >= 0 {
			prop.LotSize = sql.NullFloat64{Float64: lotSize, Valid: true}
		} else {
			form.FieldErrors["lot_size"] = "Lot size must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["lot_size"] = "Invalid value for lot size"
	}

	if squareFootage, err := strconv.ParseInt(r.PostForm.Get("square_footage"), 10, 64); err == nil {
		if squareFootage >= 0 {
			prop.SquareFt = sql.NullInt64{Int64: squareFootage, Valid: true}
		} else {
			form.FieldErrors["square_footage"] = "Square footage must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["square_footage"] = "Invalid value for square footage"
	}

	if bedrooms, err := strconv.ParseInt(r.PostForm.Get("bedrooms"), 10, 16); err == nil {
		if bedrooms >= 0 {
			prop.Bedrooms = sql.NullInt16{Int16: int16(bedrooms), Valid: true}
		} else {
			form.FieldErrors["bedrooms"] = "Bedrooms must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["bedrooms"] = "Invalid value for bedrooms"
	}

	if bathrooms, err := strconv.ParseFloat(r.PostForm.Get("bathrooms"), 64); err == nil {
		if bathrooms >= 0 {
			prop.Bathrooms = sql.NullFloat64{Float64: bathrooms, Valid: true}
		} else {
			form.FieldErrors["bathrooms"] = "Bathrooms must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["bathrooms"] = "Invalid value for bathrooms"
	}

	if yearBuilt, err := strconv.ParseInt(r.PostForm.Get("year_built"), 10, 16); err == nil {
		if yearBuilt >= 0 {
			prop.YearBuilt = sql.NullInt16{Int16: int16(yearBuilt), Valid: true}
		} else {
			form.FieldErrors["year_built"] = "Year built must be greater than or equal to 0"
		}
	} else {
		form.FieldErrors["year_built"] = "Invalid value for year built"
	}

	// URLs
	prop.TaxURL = sql.NullString{String: r.PostForm.Get("tax_assessor_url"), Valid: true}
	prop.ZillowURL = sql.NullString{String: r.PostForm.Get("zillow_url"), Valid: true}

	// Check for field errors
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "createProperty.tmpl", data)
		return
	}

	// Insert property
	prop.ID, err = app.properties.Insert(&prop)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Success
	app.sessionManager.Put(r.Context(), "flash", "Property successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/property/%d", prop.ID), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
