package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Property struct {
	ID              int64           `json:"id"`
	Address         sql.NullString  `json:"address"`
	City            sql.NullString  `json:"city"`
	State           sql.NullString  `json:"state"`
	Zip             sql.NullString  `json:"zip"`
	ParcelID        sql.NullString  `json:"parcel_id"`
	PropertyType    sql.NullString  `json:"property_type"`
	LandValue       sql.NullInt64   `json:"land_value"`
	BuildingValue   sql.NullInt64   `json:"building_value"`
	FairMarketValue sql.NullInt64   `json:"fair_market_value"`
	LotSize         sql.NullFloat64 `json:"lot_size"`
}

type PropertyModel struct {
	DB *sql.DB
}

func (m *PropertyModel) Insert(property *Property) (int64, error) {

	query := `
		INSERT INTO properties (situs, city, "state", zip_code, parcel_id, 	property_type, land_value, building_value, fair_market_value, lot_size)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING property_id`

	args := []any{property.Address, property.City, property.State, property.Zip, property.ParcelID, property.PropertyType, property.LandValue, property.BuildingValue, property.FairMarketValue, property.LotSize}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, ErrNoRecord
	}

	fmt.Println(id)

	return id, nil
}

func (m *PropertyModel) Get(id int64) (*Property, error) {
	if id < 1 {
		return nil, ErrNoRecord
	}

	query := `
		SELECT property_id, situs, city, state, zip_code, parcel_id, property_type, land_value, building_value, fair_market_value, lot_size FROM properties
		WHERE property_id = $1`

	var property Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&property.ID,
		&property.Address,
		&property.City,
		&property.State,
		&property.Zip,
		&property.ParcelID,
		&property.PropertyType,
		&property.LandValue,
		&property.BuildingValue,
		&property.FairMarketValue,
		&property.LotSize,
	)
	if err != nil {
		if err == ErrNoRecord {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &property, nil
}

func (m *PropertyModel) Latest() ([]Property, error) {
	return nil, nil
}

func (model *PropertyModel) FormatMoney(price int64) string {
	m := strconv.FormatInt(price, 10)
	n := len(m)
	if n <= 3 {
		return m
	}

	var result string
	for i := n - 1; i >= 0; i-- {
		result = string(m[i]) + result
		if (n-i)%3 == 0 && i != 0 {
			result = "," + result
		}
	}
	return result
}
