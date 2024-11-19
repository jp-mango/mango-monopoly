package models

import (
	"context"
	"database/sql"
	"time"
)

// TODO: check is the json tag is even doing anything????
type Property struct {
	ID              int64           `json:"id"`
	Address         sql.NullString  `json:"address"`
	City            sql.NullString  `json:"city"`
	State           sql.NullString  `json:"state"`
	Zip             sql.NullString  `json:"zip"`
	County          sql.NullString  `json:"county_id"`
	ParcelID        sql.NullString  `json:"parcel_id"`
	PropertyType    sql.NullString  `json:"property_type"`
	LandValue       sql.NullInt64   `json:"land_value"`
	BuildingValue   sql.NullInt64   `json:"building_value"`
	FairMarketValue sql.NullInt64   `json:"fair_market_value"`
	LotSize         sql.NullFloat64 `json:"lot_size"`
	SquareFt        sql.NullInt64   `json:"square_footage"`
	Bedrooms        sql.NullInt16   `json:"bedrooms"`
	Bathrooms       sql.NullFloat64 `json:"bathrooms"`
	YearBuilt       sql.NullInt16   `json:"year_built"`
	TaxURL          sql.NullString  `json:"tax_assessor_url"`
	ZillowURL       sql.NullString  `json:"zillow_url"`
}

type PropertyModel struct {
	DB *sql.DB
}

func (m *PropertyModel) Insert(property *Property) (int64, error) {

	query := `
		INSERT INTO properties (situs, city, "state", "zip_code",county_id,parcel_id,property_type,land_value, building_value, fair_market_value,lot_size, square_footage,bedrooms,bathrooms, year_built, tax_assessor_url,zillow_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING property_id`

	args := []any{property.Address, property.City, property.State, property.Zip, property.County, property.ParcelID, property.PropertyType, property.LandValue, property.BuildingValue, property.FairMarketValue, property.LotSize, property.SquareFt, property.Bedrooms, property.Bathrooms, property.YearBuilt, property.TaxURL, property.ZillowURL}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *PropertyModel) Get(id int64) (*Property, error) {
	if id < 1 {
		return nil, ErrNoRecord
	}

	query := `
		SELECT property_id,
		situs,
		city,
		"state",
		zip_code,
		parcel_id,
		property_type,
		land_value,
		building_value,
		fair_market_value,
		lot_size
		FROM properties
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

func (m *PropertyModel) GetByParcel(parcel_id string) bool {
	if parcel_id == "" {
		return false
	}

	query := `
		SELECT property_id,
		WHERE parcel_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, parcel_id)
	if row != nil {
		return true
	}

	return false
}

func (m *PropertyModel) Latest() ([]Property, error) {
	query := `
		SELECT property_id,
		situs,
		city,
		"state",
		zip_code,
		parcel_id,
		property_type,
		land_value,
		building_value,
		fair_market_value,
		lot_size
		FROM properties
		ORDER BY property_id DESC
		LIMIT 15`

	var latestProperties []Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		if err == ErrNoRecord {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var property Property
		err := rows.Scan(
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
			return nil, err
		}
		latestProperties = append(latestProperties, property)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return latestProperties, nil
}
