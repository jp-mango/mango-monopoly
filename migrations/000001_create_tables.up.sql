CREATE TABLE IF NOT EXISTS Counties (county_id SERIAL PRIMARY KEY,
																																						name VARCHAR(100) NOT NULL,
																																						state VARCHAR(50) NOT NULL,
																																						UNIQUE (name,
																																															state));


CREATE TABLE IF NOT EXISTS Properties
	(property_id SERIAL PRIMARY KEY,
			situs TEXT, city TEXT, zip_code VARCHAR(7),
			county_id INTEGER NOT NULL,
			parcel_id VARCHAR(150) NOT NULL UNIQUE,
			property_type TEXT, property_class TEXT, grade TEXT, roof_structure TEXT, roof_cover TEXT, heating TEXT, cooling TEXT, floors NUMERIC, land_value NUMERIC, improvement_value NUMERIC, starting_bid NUMERIC, appraisal_value NUMERIC, lot_size NUMERIC, square_footage NUMERIC, bedrooms INT, bathrooms NUMERIC, year_built INT, tax_assessor_url VARCHAR(200),
			zillow_url VARCHAR(200),
			floorplan_photo TEXT,
		FOREIGN KEY (county_id) REFERENCES Counties(county_id) ON UPDATE CASCADE ON DELETE RESTRICT);


CREATE TABLE IF NOT EXISTS Property_Images (image_id SERIAL PRIMARY KEY,
																																													parcel_id VARCHAR(150) NOT NULL,
																																													image_path TEXT NOT NULL,
																																													image_description TEXT,
																																												FOREIGN KEY (parcel_id) REFERENCES Properties (parcel_id));


CREATE TABLE IF NOT EXISTS Upcoming_Sales (upcoming_sale_id SERIAL PRIMARY KEY,
																																												parcel_id VARCHAR(150) NOT NULL,
																																												"owner" VARCHAR(100),
																																												auction_date DATE, "address" VARCHAR(100),
																																												amount_due NUMERIC);


CREATE TABLE IF NOT EXISTS Past_Sales (sale_id SERIAL PRIMARY KEY,
																																								auction_date DATE NOT NULL,
																																								parcel_id VARCHAR(150) NOT NULL,
																																								previous_owner VARCHAR(150),
																																								addr TEXT, starting_bid NUMERIC, tax_deed_purchaser VARCHAR(100),
																																								winning_bid_amount NUMERIC);


CREATE TABLE IF NOT EXISTS "sessions" (token CHAR(43) PRIMARY KEY,
																																								data BYTEA NOT NULL,
																																								expiry TIMESTAMP(6) NOT NULL);


CREATE INDEX sessions_expiry_idx ON "sessions" (expiry);


CREATE TABLE users (id SERIAL PRIMARY KEY,
																					username VARCHAR(255) NOT NULL,
																					email VARCHAR(255) NOT NULL,
																					hashed_password CHAR(60) NOT NULL,
																					created TIMESTAMP NOT NULL);