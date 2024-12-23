package models

import (
	"context"
	"database/sql"
	"time"
)

type Property struct {
	ID               int64           `json:"id"`
	Address          sql.NullString  `json:"address"`
	City             sql.NullString  `json:"city"`
	Zip              sql.NullString  `json:"zip"`
	CountyID         sql.NullInt64   `json:"county_id"`
	ParcelID         sql.NullString  `json:"parcel_id"`
	PropertyType     sql.NullString  `json:"property_type"`
	PropertyClass    sql.NullString  `json:"property_class"`
	Grade            sql.NullString  `json:"grade"`
	RoofStructure    sql.NullString  `json:"roof_structure"`
	RoofCover        sql.NullString  `json:"roof_cover"`
	Heating          sql.NullString  `json:"heating"`
	Cooling          sql.NullString  `json:"cooling"`
	Floors           sql.NullFloat64 `json:"floors"`
	LandValue        sql.NullInt64   `json:"land_value"`
	ImprovementValue sql.NullInt64   `json:"improvement_value"`
	AppraisalValue   sql.NullInt64   `json:"appraisal_value"`
	StartingBid      sql.NullFloat64 `json:"starting_bid"`
	LotSize          sql.NullFloat64 `json:"lot_size"`
	SquareFt         sql.NullInt64   `json:"square_footage"`
	Bedrooms         sql.NullInt16   `json:"bedrooms"`
	Bathrooms        sql.NullFloat64 `json:"bathrooms"`
	YearBuilt        sql.NullInt16   `json:"year_built"`
	TaxURL           sql.NullString  `json:"tax_assessor_url"`
	ZillowURL        sql.NullString  `json:"zillow_url"`
	FloorPlanPhoto   sql.NullString  `json:"floorplan_photo"`
}

type PropertyModel struct {
	DB *sql.DB
}

func (m *PropertyModel) Insert(p *Property) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the property already exists
	exists, err := m.Exists(p.ParcelID.String)
	if err != nil {
		return 0, err
	}

	// If the property exists, return early without inserting
	if exists {
		return 0, nil // No new row inserted
	}

	query := `
        INSERT INTO properties (
            situs, 
            city, 
            zip_code,
            county_id,
            parcel_id,
            property_type,
            property_class,
            grade,
            roof_structure,
            roof_cover,
            heating,
            cooling,
            floors,
            land_value,
            improvement_value, 
            appraisal_value,
			starting_bid,
            lot_size, 
            square_footage,
            bedrooms,
            bathrooms, 
            year_built, 
            tax_assessor_url,
            zillow_url,
            floorplan_photo
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,$25
        ) RETURNING property_id`

	args := []any{
		p.Address, p.City, p.Zip, p.CountyID, p.ParcelID, p.PropertyType,
		p.PropertyClass, p.Grade, p.RoofStructure, p.RoofCover, p.Heating,
		p.Cooling, p.Floors, p.LandValue, p.ImprovementValue, p.AppraisalValue, p.StartingBid,
		p.LotSize, p.SquareFt, p.Bedrooms, p.Bathrooms, p.YearBuilt, p.TaxURL,
		p.ZillowURL, p.FloorPlanPhoto,
	}

	var id int64
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
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
        SELECT 
            property_id,
            situs,
            city,
            zip_code,
            county_id,
            parcel_id,
            property_type,
            property_class,
            grade,
            roof_structure,
            roof_cover,
            heating,
            cooling,
            floors,
            land_value,
            improvement_value,
            appraisal_value,
			starting_bid,
            lot_size,
            square_footage,
            bedrooms,
            bathrooms,
            year_built,
            tax_assessor_url,
            zillow_url,
            floorplan_photo
        FROM properties
        WHERE property_id = $1`

	var p Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Address, &p.City, &p.Zip, &p.CountyID, &p.ParcelID, &p.PropertyType,
		&p.PropertyClass, &p.Grade, &p.RoofStructure, &p.RoofCover, &p.Heating,
		&p.Cooling, &p.Floors, &p.LandValue, &p.ImprovementValue, &p.AppraisalValue, &p.StartingBid,
		&p.LotSize, &p.SquareFt, &p.Bedrooms, &p.Bathrooms, &p.YearBuilt,
		&p.TaxURL, &p.ZillowURL, &p.FloorPlanPhoto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &p, nil
}

func (m *PropertyModel) GetByParcel(parcelID string) (*Property, error) {
	if parcelID == "" {
		return nil, ErrNoRecord
	}

	query := `
        SELECT 
            property_id,
            situs,
            city,
            zip_code,
            county_id,
            parcel_id,
            property_type,
            property_class,
            grade,
            roof_structure,
            roof_cover,
            heating,
            cooling,
            floors,
            land_value,
            improvement_value,
            appraisal_value,
			starting_bid,
            lot_size,
            square_footage,
            bedrooms,
            bathrooms,
            year_built,
            tax_assessor_url,
            zillow_url,
            floorplan_photo
        FROM properties
        WHERE parcel_id = $1`

	var p Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, parcelID).Scan(
		&p.ID, &p.Address, &p.City, &p.Zip, &p.CountyID, &p.ParcelID, &p.PropertyType,
		&p.PropertyClass, &p.Grade, &p.RoofStructure, &p.RoofCover, &p.Heating,
		&p.Cooling, &p.Floors, &p.LandValue, &p.ImprovementValue, &p.AppraisalValue, &p.StartingBid,
		&p.LotSize, &p.SquareFt, &p.Bedrooms, &p.Bathrooms, &p.YearBuilt,
		&p.TaxURL, &p.ZillowURL, &p.FloorPlanPhoto,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &p, nil
}

func (m *PropertyModel) Latest() ([]Property, error) {
	query := `
        SELECT 
            property_id,
            situs,
            city,
            zip_code,
            county_id,
            parcel_id,
            property_type,
            property_class,
            grade,
            roof_structure,
            roof_cover,
            heating,
            cooling,
            floors,
            land_value,
            improvement_value,
            appraisal_value,
			starting_bid,
            lot_size,
            square_footage,
            bedrooms,
            bathrooms,
            year_built,
            tax_assessor_url,
            zillow_url,
            floorplan_photo
        FROM properties
        ORDER BY property_id DESC
        LIMIT 30`

	var properties []Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Property
		err := rows.Scan(
			&p.ID, &p.Address, &p.City, &p.Zip, &p.CountyID, &p.ParcelID, &p.PropertyType,
			&p.PropertyClass, &p.Grade, &p.RoofStructure, &p.RoofCover, &p.Heating,
			&p.Cooling, &p.Floors, &p.LandValue, &p.ImprovementValue, &p.AppraisalValue, &p.StartingBid,
			&p.LotSize, &p.SquareFt, &p.Bedrooms, &p.Bathrooms, &p.YearBuilt,
			&p.TaxURL, &p.ZillowURL, &p.FloorPlanPhoto,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return properties, nil
}

func (m *PropertyModel) Exists(parcelID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM properties WHERE parcel_id = $1)`
	err := m.DB.QueryRow(query, parcelID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
