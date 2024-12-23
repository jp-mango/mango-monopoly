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

	if property.CountyID.Valid {
		county, err := app.counties.Get(property.CountyID.Int64)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.logger.Warn("County not found", "CountyID", property.CountyID.Int64)
			} else {
				app.serverError(w, r, err)
				return
			}
		} else {
			data.County = *county
		}
	}

	app.render(w, r, http.StatusOK, "property.tmpl", data)
}

func (app *application) createPropertyPage(w http.ResponseWriter, r *http.Request) {
	c, err := app.counties.All()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Counties = c
	data.Form = propertyCreateForm{}
	app.render(w, r, http.StatusOK, "createProperty.tmpl", data)
}

type propertyCreateForm struct {
	Address             string `form:"address"`
	City                string `form:"city"`
	Zip                 string `form:"zip_code"`
	CountyID            string `form:"county_id"`
	ParcelID            string `form:"parcel_id"`
	PropertyType        string `form:"property_type"`
	LandValue           string `form:"land_value"`
	ImprovementValue    string `form:"improvement_value"`
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
	prop.Address = sql.NullString{String: form.Address, Valid: validator.NotBlank(form.Address)}
	if !prop.Address.Valid {
		form.AddFieldError("address", "Cannot be left blank")
	}

	// City
	prop.City = sql.NullString{String: form.City, Valid: validator.NotBlank(form.City)}
	if !prop.City.Valid {
		form.AddFieldError("city", "Cannot be left blank")
	}

	// Zip Code
	prop.Zip = sql.NullString{String: form.Zip, Valid: validator.NotBlank(form.Zip) && validator.MaxChars(form.Zip, 7)}
	if !prop.Zip.Valid {
		if !validator.NotBlank(form.Zip) {
			form.AddFieldError("zip_code", "Cannot be left blank")
		} else {
			form.AddFieldError("zip_code", "Cannot be longer than 7 characters")
		}
	}

	// CountyID
	countyID, err := strconv.ParseInt(form.CountyID, 10, 64)
	if err != nil {
		form.AddFieldError("county_id", "Invalid County")
	}
	prop.CountyID = sql.NullInt64{Int64: countyID, Valid: err == nil}
	if !prop.CountyID.Valid {
		form.AddFieldError("county_id", "Cannot be left blank")
	}

	// Parcel ID
	prop.ParcelID = sql.NullString{String: form.ParcelID, Valid: validator.NotBlank(form.ParcelID)}
	if !prop.ParcelID.Valid {
		form.AddFieldError("parcel_id", "Cannot be left blank")
	} else {
		// Check if Parcel ID already exists
		exists, err := app.properties.GetByParcel(prop.ParcelID.String)
		if err != nil && err != models.ErrNoRecord {
			app.serverError(w, r, err)
			return
		}
		if exists != nil { // Property exists
			form.AddFieldError("parcel_id", "Property already exists")
		}
	}

	// Property Type
	prop.PropertyType = sql.NullString{String: form.PropertyType, Valid: validator.NotBlank(form.PropertyType)}
	if !prop.PropertyType.Valid {
		form.AddFieldError("property_type", "Cannot be left blank")
	}

	// Numeric Fields
	if form.LandValue != "" {
		landValue, err := strconv.ParseInt(form.LandValue, 10, 64)
		if err != nil {
			form.AddFieldError("land_value", "Enter an integer")
		} else {
			prop.LandValue = sql.NullInt64{Int64: landValue, Valid: true}
		}
	}

	if form.ImprovementValue != "" {
		improvementValue, err := strconv.ParseInt(form.ImprovementValue, 10, 64)
		if err != nil {
			form.AddFieldError("improvement_value", "Enter an integer")
		} else {
			prop.ImprovementValue = sql.NullInt64{Int64: improvementValue, Valid: true}
		}
	}

	if form.AppraisalValue != "" {
		appraisalValue, err := strconv.ParseInt(form.AppraisalValue, 10, 64)
		if err != nil {
			form.AddFieldError("appraisal_value", "Enter an integer")
		} else {
			prop.AppraisalValue = sql.NullInt64{Int64: appraisalValue, Valid: true}
		}
	}

	if form.LotSize != "" {
		lotSize, err := strconv.ParseFloat(form.LotSize, 64)
		if err != nil {
			form.AddFieldError("lot_size", "Enter a number")
		} else {
			prop.LotSize = sql.NullFloat64{Float64: lotSize, Valid: true}
		}
	}

	if form.SquareFootage != "" {
		squareFootage, err := strconv.ParseInt(form.SquareFootage, 10, 64)
		if err != nil {
			form.AddFieldError("square_footage", "Enter a number")
		} else {
			prop.SquareFt = sql.NullInt64{Int64: squareFootage, Valid: true}
		}
	}

	if form.Bedrooms != "" {
		bedrooms, err := strconv.ParseInt(form.Bedrooms, 10, 16)
		if err != nil {
			form.AddFieldError("bedrooms", "Enter a number")
		} else {
			prop.Bedrooms = sql.NullInt16{Int16: int16(bedrooms), Valid: true}
		}
	}

	if form.Bathrooms != "" {
		bathrooms, err := strconv.ParseFloat(form.Bathrooms, 64)
		if err != nil {
			form.AddFieldError("bathrooms", "Enter a number")
		} else {
			prop.Bathrooms = sql.NullFloat64{Float64: bathrooms, Valid: true}
		}
	}

	if form.YearBuilt != "" {
		yearBuilt, err := strconv.ParseInt(form.YearBuilt, 10, 16)
		if err != nil {
			form.AddFieldError("year_built", "Enter a number")
		} else {
			prop.YearBuilt = sql.NullInt16{Int16: int16(yearBuilt), Valid: true}
		}
	}

	// URLs
	prop.TaxURL = sql.NullString{String: form.TaxURL, Valid: form.TaxURL != ""}
	prop.ZillowURL = sql.NullString{String: form.ZillowURL, Valid: form.ZillowURL != ""}
	prop.FloorPlanPhoto = sql.NullString{String: r.PostForm.Get("floorplan_photo"), Valid: r.PostForm.Get("floorplan_photo") != ""}

	// Check for field errors
	if !form.Valid() {
		counties, err := app.counties.All()
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		data := app.newTemplateData(r)
		data.Form = form
		data.Counties = counties
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
		app.serverError(w, r, fmt.Errorf("error scraping auction data: %w", err))
		return
	}

	parcelIDs, startingBids, err := scraper.ProcessCSV("./scraper/Gwinnett")
	if err != nil {
		app.serverError(w, r, fmt.Errorf("error processing CSV: %w", err))
		return
	}

	properties, err := scraper.ScrapeGwinnettParcelData(parcelIDs, startingBids)
	if err != nil {
		app.serverError(w, r, fmt.Errorf("error scraping parcel data: %w", err))
		return
	}

	for _, prop := range properties {
		// Check if the property already exists
		exists, err := app.properties.Exists(prop.ParcelID.String)
		if err != nil {
			app.logger.Error("error checking property existence", "err", err)
			continue
		}

		if exists {
			app.logger.Info("property already exists", "parcel_id", prop.ParcelID.String)
			continue
		}

		// Insert only if it doesn't exist
		id, err := app.properties.Insert(prop)
		if err != nil {
			app.logger.Error("error inserting property", "parcel id", prop.ParcelID, "err", err)
			continue
		}

		app.logger.Info(fmt.Sprintf("successfully inserted property %d", id))
	}

	http.Redirect(w, r, "/properties", http.StatusSeeOther)
}
