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

func (m *PropertyModel) Insert(address, city, state, zip, parcelID, propertyType string, landValue, buildingValue, fmv, lotsize float32) (int, error) {

	//TODO:fix query, maybe?
	stmt := `INSERT INTO properties (situs, city, "state", zip_code, parcel_id, property_type, land_value, building_value, fair_market_value, lot_size)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING property_id`

	var id int

	err := m.DB.QueryRow(stmt, address, city, state, zip, parcelID, propertyType, landValue, buildingValue, fmv, lotsize).Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (m *PropertyModel) Get(id int) (Property, error) {
	return Property{}, nil
}

func (m *PropertyModel) Latest() ([]Property, error) {
	return nil, nil
}
