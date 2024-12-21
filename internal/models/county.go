package models

import (
	"database/sql"
)

// County represents a record in the Counties table.
type County struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// CountyModel wraps a sql.DB connection pool.
type CountyModel struct {
	DB *sql.DB
}

// All fetches all counties from the database.
func (m *CountyModel) All() ([]County, error) {
	query := `
		SELECT county_id, name, state
		FROM counties
		ORDER BY name ASC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counties []County

	for rows.Next() {
		var c County
		err := rows.Scan(&c.ID, &c.Name, &c.State)
		if err != nil {
			return nil, err
		}
		counties = append(counties, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counties, nil
}

func (m *CountyModel) Get(id int64) (*County, error) {
	query := `
		SELECT county_id, name, state
		FROM counties
		WHERE county_id = $1`

	var county County

	err := m.DB.QueryRow(query, id).Scan(&county.ID, &county.Name, &county.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &county, nil
}
