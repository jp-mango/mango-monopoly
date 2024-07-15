CREATE TABLE IF NOT EXISTS Counties (county_id SERIAL PRIMARY KEY,
																																						name VARCHAR(100) NOT NULL,
																																						state VARCHAR(50) NOT NULL);


CREATE TABLE IF NOT EXISTS Properties (property_id SERIAL PRIMARY KEY,
																																								situs TEXT NOT NULL UNIQUE,
																																								county_id INT NOT NULL,
																																								parcel_id VARCHAR(150) NOT NULL UNIQUE,
																																								property_type TEXT, land_value NUMERIC, building_value NUMERIC, fair_market_value NUMERIC, lot_size NUMERIC, square_footage NUMERIC, bedrooms INT, bathrooms INT, year_built INT,
																																							FOREIGN KEY (county_id) REFERENCES Counties(county_id));


CREATE TABLE IF NOT EXISTS Property_Images (image_id SERIAL PRIMARY KEY,
																																													parcel_id VARCHAR(150) NOT NULL,
																																													image_path TEXT NOT NULL,
																																													image_description TEXT,
																																												FOREIGN KEY (parcel_id) REFERENCES Properties (parcel_id));


CREATE TABLE IF NOT EXISTS Upcoming_Sales (upcoming_sale_id SERIAL PRIMARY KEY,
																																												parcel_id VARCHAR(150) NOT NULL,
																																												"owner" VARCHAR(100),
																																												auction_date DATE, "address" VARCHAR(100),
																																												amount_due VARCHAR(20));


CREATE TABLE IF NOT EXISTS Past_Sales (sale_id SERIAL PRIMARY KEY,
																																								auction_date DATE NOT NULL,
																																								parcel_id VARCHAR(150) NOT NULL,
																																								previous_owner VARCHAR(150),
																																								addr TEXT, starting_bid VARCHAR(20),
																																								tax_deed_purchaser VARCHAR(100),
																																								winning_bid_amount VARCHAR(20));

