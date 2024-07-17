-- Trigger function to insert into properties if it doesn't exist

CREATE OR REPLACE FUNCTION insert_property_if_not_exists() RETURNS TRIGGER AS $$
BEGIN
    -- Insert into Properties if the parcel_id does not exist
    IF NOT EXISTS (SELECT 1 FROM Properties WHERE parcel_id = NEW.parcel_id) THEN
        INSERT INTO Properties (parcel_id)
        VALUES (NEW.parcel_id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for Upcoming_Sales

CREATE TRIGGER insert_property_on_upcoming_sales_insert AFTER
INSERT ON Upcoming_Sales
FOR EACH ROW EXECUTE FUNCTION insert_property_if_not_exists();

-- Trigger for Past_Sales

CREATE TRIGGER insert_property_on_past_sales_insert AFTER
INSERT ON Past_Sales
FOR EACH ROW EXECUTE FUNCTION insert_property_if_not_exists();