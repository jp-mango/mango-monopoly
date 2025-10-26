package main

import (
	"fmt"
	"mango-monopoly/internal/data"
	"net/http"
	"time"
)

func (app *application) createPropertyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "add a new property")
}

func (app *application) showPropertyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	//DUMMY DATA
	property := data.Property{
		ID:        id,
		CreatedAt: time.Now(),
		Version:   1,

		County: "Gwinnett",
		PropID: "R7269A060",
		AltID:  "1423341",
		Owners: []string{"MENZIES ALLYSON", "MENZIES ERIC"},

		SitusAddress: "4585 S LEE ST",
		OwnerAddress: "4585 S LEE ST",
		City:         "BUFORD",
		State:        "GA",
		ZIP:          "30518-5735",

		Class:        "Residential SFR",
		Neighborhood: "12143002",
		Acres:        0.25,
		Type:         "Ranch",
		Grade:        "C",
		YearBuilt:    1959,

		// Building Attributes
		Occupancy:     "Single Family",
		RoofStruct:    "Gable-Hip",
		RoofCover:     "Composition Shingle 240-260#",
		Heating:       "Central Heat",
		AC:            "Central Air",
		Stories:       1.0,
		Bedrooms:      2,
		Bathrooms:     1,
		ExteriorWall:  "Wood Siding",
		InteriorFloor: "Carpet/Vinyl",

		// Valuation (latest year)
		ValueYear:        2024,
		ValueReason:      "Adjusted for Market Conditions",
		LandValue:        50000,
		ImprovementValue: 95000,
		TotalAppraised:   145000,
		TotalAssessed:    58000,

		// Transfer (latest deed)
		TransferBook: "59866",
		TransferPage: "781",
		TransferDate: time.Date(2022, 4, 8, 0, 0, 0, 0, time.UTC),
		Grantor:      "WE-FLIP LLC",
		Grantee:      "MENZIES ALLYSON",
		DeedType:     "Fee Simple Warranty Deed",
		VacantLand:   false,
		SalePrice:    305000,

		// Land / Features
		LandPrimaryUse: "R01 - Primary Site",
		LandType:       "Residential",
		EffFrontage:    80,
		EffDepth:       135,
		FloorAreaCode:  "1.0",
		FloorDesc:      "Main Floor",
		GrossAreaSF:    1032,
		FinishedSF:     1032,
		FloorConstr:    "Wood Frame",
		ExtFeatCode:    "OFP",
		ExtFeatDesc:    "Open Frame Porch",
		ExtFeatSizeSF:  132,
		ExtFeatConstr:  "Wood",

		// Legal / References
		LegalDescription: []string{
			"LOT 6 BLK A SUB 2 RIDGEWOOD",
			"PB H PG 76",
		},
		SourceURL: "https://gwinnettassessor.manatron.com/IWantTo/PropertyGISSearch/PropertyDetail.aspx?p=R7269A060&a=1423341",
		Photos: []string{
			"http://gis.aumentumtech.com/arcgis/rest/directories/arcgisoutput/GwinnettGa/ExportWebMap_GPServer/_ags_9d0f8dc2-b2a1-11f0-a6e2-000d3ae60619.jpg",
			"http://sotf.publicaccessnow.com/OmniSketchASPNet/sotf_omniSketch.ashx?a=gwinnett&d=2&p=1423341&c=R01&s=md&e=24",
		},
	}

	err = app.writeJSON(w, http.StatusOK, property, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
