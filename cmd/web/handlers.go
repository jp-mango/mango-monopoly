package main

import (
	"database/sql"
	"errors"
	"fmt"
	"mango-monopoly/internal/models"
	"mango-monopoly/internal/validator"
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
	validator.Validator
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
	form.CheckField(!app.properties.GetByParcel(prop.ParcelID.String), "parcel_id", "Property already exists")

	// Property Type
	prop.PropertyType = sql.NullString{String: r.PostForm.Get("property_type"), Valid: true}
	form.CheckField(validator.NotBlank(prop.PropertyType.String), "property_type", "Cannot be left blank")

	// Numeric fields
	valid, errMsg := validator.ValidateInt(r.PostForm.Get("land_value"), 64)
	if !valid {
		form.AddFieldError("land_value", errMsg)
	}

	valid, errMsg = validator.ValidateInt(r.PostForm.Get("building_value"), 64)
	if !valid {
		form.AddFieldError("building_value", errMsg)
	}

	valid, errMsg = validator.ValidateInt(r.PostForm.Get("fair_market_value"), 64)
	if !valid {
		form.AddFieldError("fair_market_value", errMsg)
	}

	valid, errMsg = validator.ValidateInt(r.PostForm.Get("lot_size"), 64)
	if !valid {
		form.AddFieldError("lot_size", errMsg)
	}

	valid, errMsg = validator.ValidateInt(r.PostForm.Get("square_footage"), 64)
	if !valid {
		form.AddFieldError("square_footage", errMsg)
	}

	valid, errMsg = validator.ValidateInt(r.PostForm.Get("bedrooms"), 64)
	if !valid {
		form.AddFieldError("bedrooms", errMsg)
	}

	valid, errMsg = validator.ValidateFloat(r.PostForm.Get("bathrooms"), 64)
	if !valid {
		form.AddFieldError("bathrooms", errMsg)
	}

	valid, errMsg = validator.ValidateFloat(r.PostForm.Get("year_built"), 64)
	if !valid {
		form.AddFieldError("year_built", errMsg)
	}

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
