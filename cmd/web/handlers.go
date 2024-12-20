package main

import (
	"database/sql"
	"errors"
	"fmt"
	"mango-monopoly/internal/models"
	"mango-monopoly/internal/validator"
	"mango-monopoly/scraper"
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
	data.Form = propertyCreateForm{}
	app.render(w, r, http.StatusOK, "createProperty.tmpl", data)
}

type propertyCreateForm struct {
	Address             string `form:"address"`
	City                string `form:"city"`
	State               string `form:"state"`
	Zip                 string `form:"zip_code"`
	County              string `form:"county_id"`
	ParcelID            string `form:"parcel_id"`
	PropertyType        string `form:"property_type"`
	LandValue           string `form:"land_value"`
	BuildingValue       string `form:"building_value"`
	AppraisalValue      string `form:"appraisal_value"`
	LotSize             string `form:"lot_size"`
	SquareFootage       string `form:"square_footage"`
	Bedrooms            string `form:"bedrooms"`
	Bathrooms           string `form:"bathrooms"`
	YearBuilt           string `form:"year_built"`
	TaxURL              string `form:"tax_assessor_url"`
	ZillowURL           string `form:"zillow_url"`
	validator.Validator `form:"-"`
}

func (app *application) propertyCreatePost(w http.ResponseWriter, r *http.Request) {
	var form propertyCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	prop := models.Property{}

	// Address
	prop.Address = sql.NullString{String: r.PostForm.Get("address"), Valid: true}
	form.CheckField(validator.NotBlank(prop.Address.String), "address", "Cannot be left blank")

	// City
	prop.City = sql.NullString{String: r.PostForm.Get("city"), Valid: true}
	form.CheckField(validator.NotBlank(prop.City.String), "city", "Cannot be left blank")

	// State
	prop.State = sql.NullString{String: r.PostForm.Get("state"), Valid: true}
	form.CheckField(validator.NotBlank(prop.State.String), "state", "Cannot be left blank")

	// Zip Code
	prop.Zip = sql.NullString{String: r.PostForm.Get("zip_code"), Valid: true}
	form.CheckField(validator.MaxChars(prop.Zip.String, 7), "zip_code", "Cannot be longer than 7 characters")
	form.CheckField(validator.NotBlank(prop.Zip.String), "zip_code", "Cannot be left blank")

	// County
	prop.County = sql.NullString{String: r.PostForm.Get("county_id"), Valid: true}
	form.CheckField(validator.NotBlank(prop.County.String), "county_id", "Cannot be left blank")

	// Parcel ID
	prop.ParcelID = sql.NullString{String: r.PostForm.Get("parcel_id"), Valid: true}
	form.CheckField(validator.NotBlank(prop.ParcelID.String), "parcel_id", "Cannot be left blank")
	form.CheckField(app.properties.GetByParcel(prop.ParcelID.String), "parcel_id", "Property already exists")

	// Property Type
	prop.PropertyType = sql.NullString{String: r.PostForm.Get("property_type"), Valid: true}
	form.CheckField(validator.NotBlank(prop.PropertyType.String), "property_type", "Cannot be left blank")

	// Numeric fields
	landValue, err := strconv.ParseInt(r.PostForm.Get("land_value"), 10, 64)
	if err != nil {
		form.AddFieldError("land_value", "Enter an integer")
	}
	prop.LandValue = sql.NullInt64{Int64: landValue, Valid: err == nil}

	buildingValue, err := strconv.ParseInt(r.PostForm.Get("building_value"), 10, 64)
	if err != nil {
		form.AddFieldError("building_value", "Enter an integer")
	}
	prop.LandValue = sql.NullInt64{Int64: buildingValue, Valid: err == nil}

	fairMarketValue, err := strconv.ParseInt(r.PostForm.Get("appraisal_value"), 10, 64)
	if err != nil {
		form.AddFieldError("appraisal_value", "Enter an integer")
	}
	prop.AppraisalValue = sql.NullInt64{Int64: fairMarketValue, Valid: err == nil}

	lotSize, err := strconv.ParseFloat(r.PostForm.Get("lot_size"), 64)
	if err != nil {
		form.AddFieldError("lot_size", "Enter a number")
	}
	prop.LotSize = sql.NullFloat64{Float64: lotSize, Valid: err == nil}

	squareFootage, err := strconv.ParseInt(r.PostForm.Get("square_footage"), 10, 64)
	if err != nil {
		form.AddFieldError("square_footage", "Enter a number")
	}
	prop.SquareFt = sql.NullInt64{Int64: squareFootage, Valid: err == nil}

	bedrooms, err := strconv.ParseInt(r.PostForm.Get("bedrooms"), 10, 16)
	if err != nil {
		form.AddFieldError("bedrooms", "Enter a number")
	}
	prop.Bedrooms = sql.NullInt16{Int16: int16(bedrooms), Valid: err == nil}

	bathrooms, err := strconv.ParseFloat(r.PostForm.Get("bathrooms"), 64)
	if err != nil {
		form.AddFieldError("bathrooms", "Enter a number")
	}
	prop.Bathrooms = sql.NullFloat64{Float64: float64(bathrooms), Valid: err == nil}

	yearBuilt, err := strconv.ParseInt(r.PostForm.Get("year_built"), 10, 64)
	if err != nil {
		form.AddFieldError("year_built", "Enter a number")
	}
	prop.YearBuilt = sql.NullInt16{Int16: int16(yearBuilt), Valid: err == nil}

	// URLs
	prop.TaxURL = sql.NullString{String: r.PostForm.Get("tax_assessor_url"), Valid: true}

	prop.ZillowURL = sql.NullString{String: r.PostForm.Get("zillow_url"), Valid: true}

	// Check for field errors
	if !form.Valid() {
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

type userSignupForm struct {
	Email               string `form:"email"`
	Username            string `form:"username"`
	Password            string `form:"password"`
	PwConfirm           string `form:"passwordConfirm"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Please enter a valid email")
	form.CheckField(!app.users.CheckUsers(form.Email), "email", "This email is already in use")

	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(!app.users.CheckUsers(form.Username), "username", "Username already in use")

	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 digits long")

	form.CheckField(form.Password != form.PwConfirm, "passwordConfirm", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Sign-up Successful! Please login")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Enter a valid email address")

	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
			//TODO: search for error of no email found then redirect user to sign up page
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	//TODO: change redirect path
	http.Redirect(w, r, "/property/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) scrapeHandler(w http.ResponseWriter, r *http.Request) {
	gwinnettData := scraper.CountyScraper{
		Name:    "Gwinnett",
		Webpage: "https://www.gwinnetttaxcommissioner.com/property-tax/delinquent_tax/tax-liens-tax-sales",
		Domain:  "www.gwinnetttaxcommissioner.com",
	}
	err := gwinnettData.ScrapeAuctionData()
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
	} else {
		fmt.Fprintf(w, "Scraping completed\n")
	}

	parcelIDs, err := scraper.ProcessCSV("./scraper/Gwinnett")
	if err != nil {
		app.serverError(w, r, fmt.Errorf("error: %w", err))
	}

	err = scraper.ScrapeGwinnettParcelData(parcelIDs)
	if err != nil {
		app.serverError(w, r, err)
	}
	/*
		pauldingData := scraper.CountyScraper{
			Name:    "Paulding",
			Webpage: "",
			Domain:  "",
		}
	*/
}
