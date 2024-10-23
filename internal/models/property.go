package models

import (
	"database/sql"
)

type Property struct {
	ID              int
	Address         string
	City            string
	State           string
	Zip             string
	County          string
	ParcelID        string
	PropertyType    string
	LandValue       float32
	BuildingValue   float32
	FairMarketValue float32
	LotSize         float32
}

type PropertyModel struct {
	DB *sql.DB
}

func (m *PropertyModel) Insert(address, city, state, zip, county, parcelID, propertyType string, landValue, buildingValue, fmv, lotsize float32) (int, error) {
	return 0, nil
}

func (m *PropertyModel) Get(id int) (Property, error) {
	return Property{}, nil
}

func (m *PropertyModel) Latest() ([]Property, error) {
	return nil, nil
}
