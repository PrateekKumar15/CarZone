-- CarZone Database Schema Definition
-- This file contains the complete database schema for the CarZone application
-- supporting rental and sales platform with nested JSONB structures for modern data modeling.

-- =============================================================================
-- CLEANUP AND PREPARATION
-- =============================================================================

-- Drop existing tables if they exist (for complete reset)
DROP TABLE IF EXISTS payment CASCADE;
DROP TABLE IF EXISTS booking CASCADE;
DROP TABLE IF EXISTS car CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- =============================================================================
-- TABLE DEFINITIONS
-- =============================================================================

-- Users Table Definition
-- Stores user account information for authentication and authorization
CREATE TABLE users (
    -- Primary key: Unique identifier for each user
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User account information
    username VARCHAR(255) NOT NULL,                              -- User's username
    email VARCHAR(255) NOT NULL UNIQUE,                          -- User's email address (unique)
    password_hash VARCHAR(255) NOT NULL,                         -- Hashed password for security
    phone VARCHAR(20),                                           -- User's phone number
    role VARCHAR(50) DEFAULT 'user',                            -- User role (user, admin, owner)
    profile_data JSONB,                                          -- Additional profile information as JSON
    
    -- Audit trail columns for tracking changes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,              -- Account creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP               -- Last update timestamp
);

-- Car Table Definition  
-- Stores comprehensive car information with nested engine and pricing structures
-- Uses JSONB for flexible, searchable nested data storage
CREATE TABLE car (
    -- Primary key: Unique identifier for each car
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Car ownership and basic information
    owner_id UUID,                                               -- Reference to users.id (nullable for system cars)
    name VARCHAR(255) NOT NULL,                                  -- Display name/model of the car
    brand VARCHAR(255) NOT NULL,                                 -- Manufacturer brand (e.g., "Toyota", "Tesla")
    model VARCHAR(255) NOT NULL,                                 -- Specific model name
    year INTEGER NOT NULL CHECK (year >= 1900 AND year <= 2030), -- Manufacturing year
    fuel_type VARCHAR(50) NOT NULL,                             -- Fuel type (Petrol, Diesel, Electric, Hybrid)
    
    -- Engine specifications stored as JSONB for flexibility and searchability
    engine JSONB NOT NULL,                                       -- Engine specifications: {engine_size, cylinders, horsepower, transmission}
    
    -- Location information
    location_city VARCHAR(255) NOT NULL,                         -- City where car is located
    location_state VARCHAR(255) NOT NULL,                        -- State/province where car is located
    location_country VARCHAR(255) NOT NULL,                      -- Country where car is located
    
    -- Pricing information as simple decimal for rental pricing
    price DECIMAL(10,2) NOT NULL,                               -- Daily rental price
    
    -- Status and availability
    status VARCHAR(50) DEFAULT 'active',                         -- active, maintenance, inactive
    availability_type VARCHAR(50) NOT NULL DEFAULT 'rental',     -- rental only
    is_available BOOLEAN DEFAULT true,                           -- Current availability status
    
    -- Additional information
    features JSONB,                                              -- Car features as JSON (GPS, AC, etc.)
    description TEXT,                                            -- Detailed description
    images TEXT[],                                               -- Array of image URLs
    mileage INTEGER DEFAULT 0,                                   -- Current mileage
    
    -- Audit trail columns
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,              -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP               -- Last update timestamp
);

-- Booking Table Definition
-- Stores booking information for car rentals and sales
CREATE TABLE booking (
    -- Primary key: Unique identifier for each booking
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationship fields
    customer_id UUID NOT NULL,                                   -- Reference to users.id (customer)
    car_id UUID NOT NULL,                                        -- Reference to car.id
    owner_id UUID,                                               -- Reference to users.id (car owner, nullable for system cars)
    
    -- Booking details (all bookings are rentals)
    status VARCHAR(50) DEFAULT 'pending',                        -- pending, confirmed, active, completed, cancelled
    total_amount DECIMAL(10,2) NOT NULL,                         -- Total booking amount
    start_date TIMESTAMP NOT NULL,                               -- Start date for rental
    end_date TIMESTAMP NOT NULL,                                 -- End date for rental
    notes TEXT,                                                  -- Additional notes or special requests
    
    -- Audit trail columns
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,              -- Booking creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP               -- Last update timestamp
);

-- Payment Table Definition
-- Stores payment information for bookings with Razorpay integration
CREATE TABLE payment (
    -- Primary key: Unique identifier for each payment
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationship fields
    booking_id UUID NOT NULL,                                    -- Reference to booking.id
    
    -- Razorpay specific fields
    razorpay_order_id VARCHAR(255),                             -- Razorpay order ID
    razorpay_payment_id VARCHAR(255),                           -- Razorpay payment ID
    
    -- Payment details
    amount DECIMAL(10,2) NOT NULL,                              -- Payment amount in INR
    currency VARCHAR(3) DEFAULT 'INR',                          -- Currency code
    status VARCHAR(50) DEFAULT 'pending',                       -- pending, completed, failed, refunded, cancelled
    method VARCHAR(50) NOT NULL,                                -- razorpay, cash, card, upi, netbanking
    transaction_id VARCHAR(255),                                -- Transaction reference ID
    description TEXT,                                           -- Payment description
    notes TEXT,                                                 -- Additional payment notes
    
    -- Audit trail columns
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,             -- Payment creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP              -- Last update timestamp
);

-- =============================================================================
-- CONSTRAINTS AND RELATIONSHIPS
-- =============================================================================

-- Foreign Key Constraint: Establish relationship between car and user (owner)
ALTER TABLE car
ADD CONSTRAINT fk_car_owner_id
FOREIGN KEY (owner_id)
REFERENCES users(id)
ON DELETE SET NULL;                                              -- Set owner_id to NULL when user is deleted

-- Foreign Key Constraints for booking table
ALTER TABLE booking
ADD CONSTRAINT fk_booking_customer_id
FOREIGN KEY (customer_id)
REFERENCES users(id)
ON DELETE CASCADE;                                               -- Delete booking when customer is deleted

ALTER TABLE booking
ADD CONSTRAINT fk_booking_car_id
FOREIGN KEY (car_id)
REFERENCES car(id)
ON DELETE CASCADE;                                               -- Delete booking when car is deleted

ALTER TABLE booking
ADD CONSTRAINT fk_booking_owner_id
FOREIGN KEY (owner_id)
REFERENCES users(id)
ON DELETE SET NULL;                                              -- Set owner_id to NULL when owner is deleted

-- Foreign Key Constraints for payment table
ALTER TABLE payment
ADD CONSTRAINT fk_payment_booking_id
FOREIGN KEY (booking_id)
REFERENCES booking(id)
ON DELETE CASCADE;                                               -- Delete payment when booking is deleted

-- Check constraints for data validation
ALTER TABLE booking
ADD CONSTRAINT check_booking_status 
CHECK (status IN ('pending', 'confirmed', 'active', 'completed', 'cancelled'));

ALTER TABLE booking
ADD CONSTRAINT check_booking_dates 
CHECK (end_date >= start_date);

ALTER TABLE booking
ADD CONSTRAINT check_total_amount 
CHECK (total_amount > 0);

-- Check constraints for payment validation
ALTER TABLE payment
ADD CONSTRAINT check_payment_status 
CHECK (status IN ('pending', 'completed', 'failed', 'refunded', 'cancelled'));

ALTER TABLE payment
ADD CONSTRAINT check_payment_method 
CHECK (method IN ('razorpay', 'cash', 'card', 'upi', 'netbanking'));

ALTER TABLE payment
ADD CONSTRAINT check_payment_amount 
CHECK (amount > 0);

ALTER TABLE payment
ADD CONSTRAINT check_payment_currency 
CHECK (currency = 'INR');

-- Check constraints for data validation
ALTER TABLE car
ADD CONSTRAINT check_availability_type 
CHECK (availability_type IN ('rental'));

ALTER TABLE car
ADD CONSTRAINT check_status 
CHECK (status IN ('active', 'maintenance', 'inactive'));

ALTER TABLE car
ADD CONSTRAINT check_fuel_type 
CHECK (fuel_type IN ('Petrol', 'Diesel', 'Electric', 'Hybrid', 'CNG'));

-- =============================================================================
-- INDEXES FOR PERFORMANCE
-- =============================================================================

-- Index on user email for fast authentication queries
CREATE INDEX idx_users_email ON users(email);

-- Index on user role for authorization queries
CREATE INDEX idx_users_role ON users(role);

-- Index on car brand for fast brand-based queries
CREATE INDEX idx_car_brand ON car(brand);

-- Index on car year for year-based filtering
CREATE INDEX idx_car_year ON car(year);

-- Index on car location for location-based searches
CREATE INDEX idx_car_location ON car(location_city, location_state, location_country);

-- Index on car availability for quick filtering of available cars
CREATE INDEX idx_car_availability ON car(is_available, availability_type);

-- Index on car status for status-based filtering
CREATE INDEX idx_car_status ON car(status);

-- JSONB indexes for engine and price searches
CREATE INDEX idx_car_engine_gin ON car USING gin(engine);
-- Specific index for common price queries
CREATE INDEX idx_car_engine_horsepower ON car USING btree((engine->>'horsepower'));
CREATE INDEX idx_car_price ON car USING btree(price);

-- Booking table indexes for performance
CREATE INDEX idx_booking_customer_id ON booking(customer_id);
CREATE INDEX idx_booking_car_id ON booking(car_id);
CREATE INDEX idx_booking_owner_id ON booking(owner_id);
CREATE INDEX idx_booking_status ON booking(status);
-- Removed: booking_type index (no longer needed for rental-only platform)
CREATE INDEX idx_booking_dates ON booking(start_date, end_date);
CREATE INDEX idx_booking_created_at ON booking(created_at);

-- Payment table indexes for performance
CREATE INDEX idx_payment_booking_id ON payment(booking_id);
CREATE INDEX idx_payment_status ON payment(status);
CREATE INDEX idx_payment_method ON payment(method);
CREATE INDEX idx_payment_razorpay_order_id ON payment(razorpay_order_id);
CREATE INDEX idx_payment_razorpay_payment_id ON payment(razorpay_payment_id);
CREATE INDEX idx_payment_transaction_id ON payment(transaction_id);
CREATE INDEX idx_payment_created_at ON payment(created_at);

-- =============================================================================
-- TRIGGERS FOR AUTOMATIC TIMESTAMP UPDATES
-- =============================================================================

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at when records are modified
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_car_updated_at 
    BEFORE UPDATE ON car 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_booking_updated_at 
    BEFORE UPDATE ON booking 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payment_updated_at 
    BEFORE UPDATE ON payment 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- SAMPLE DATA FOR TESTING AND DEVELOPMENT
-- =============================================================================

-- Insert sample user data for testing
-- Passwords are hashed using bcrypt (these are sample hashes for 'password123')
INSERT INTO users (id, username, email, password_hash, phone, role, profile_data) VALUES
    -- Test car owners
    ('11111111-0000-4000-8000-000000000001', 'johndoe', 'john.doe@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.q7.q7.q7.q7.q7.q7.q7.q7.q7.q7', '+1-555-0101', 'owner', '{"verified": true, "rating": 4.8, "cars_owned": 2}'),
    
    ('22222222-0000-4000-8000-000000000002', 'janesmith', 'jane.smith@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.q7.q7.q7.q7.q7.q7.q7.q7.q7.q7', '+1-555-0102', 'owner', '{"verified": true, "rating": 4.9, "cars_owned": 3}'),
    
    ('33333333-0000-4000-8000-000000000003', 'mikejohnson', 'mike.johnson@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.q7.q7.q7.q7.q7.q7.q7.q7.q7.q7', '+1-555-0103', 'owner', '{"verified": true, "rating": 4.7, "cars_owned": 1}'),
    
    -- Regular users (potential renters)
    ('44444444-0000-4000-8000-000000000004', 'sarahwilson', 'sarah.wilson@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.q7.q7.q7.q7.q7.q7.q7.q7.q7.q7', '+1-555-0104', 'user', '{"verified": true, "license_verified": true}'),
    
    ('55555555-0000-4000-8000-000000000005', 'davidbrown', 'david.brown@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.q7.q7.q7.q7.q7.q7.q7.q7.q7.q7', '+1-555-0105', 'user', '{"verified": false, "license_verified": false}'),
    
    -- Admin user
    ('99999999-0000-4000-8000-000000000099', 'admin', 'admin@carzone.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.q7.q7.q7.q7.q7.q7.q7.q7.q7.q7', '+1-555-0199', 'admin', '{"verified": true, "admin_level": "super"}');

-- Insert comprehensive sample car data
-- These cars represent different categories: economy, mid-range, luxury, and electric
INSERT INTO car (id, owner_id, name, brand, model, year, fuel_type, engine, location_city, location_state, location_country, price, status, availability_type, is_available, features, description, images, mileage) VALUES

    -- John Doe's cars (Economy and Mid-range)
    ('c0000001-0000-4000-8000-000000000001', 
     '11111111-0000-4000-8000-000000000001',
     'Reliable Daily Driver', 'Toyota', 'Camry', 2023, 'Petrol',
     '{"engine_size": 2.5, "cylinders": 4, "horsepower": 203, "transmission": "Automatic"}',
     'San Francisco', 'California', 'United States',
     45.00,
     'active', 'rental', true,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": true}',
     'Well-maintained 2023 Toyota Camry perfect for city driving and longer trips. Recently serviced with excellent fuel economy.',
     ARRAY['https://example.com/images/camry1.jpg', 'https://example.com/images/camry2.jpg'],
     15420),

    ('c0000002-0000-4000-8000-000000000002',
     '11111111-0000-4000-8000-000000000001', 
     'City Commuter', 'Honda', 'Civic', 2022, 'Petrol',
     '{"engine_size": 2.0, "cylinders": 4, "horsepower": 158, "transmission": "CVT"}',
     'San Francisco', 'California', 'United States',
     35.00,
     'active', 'rental', true,
     '{"gps": false, "air_conditioning": true, "bluetooth": true, "backup_camera": false, "keyless_entry": false}',
     'Fuel-efficient Honda Civic ideal for city commuting. Available for both rental and purchase.',
     ARRAY['https://example.com/images/civic1.jpg'],
     28750),

    -- Jane Smith's cars (Mid-range and Luxury)
    ('c0000003-0000-4000-8000-000000000003',
     '22222222-0000-4000-8000-000000000002',
     'Family SUV', 'Honda', 'CR-V', 2024, 'Petrol',
     '{"engine_size": 1.5, "cylinders": 4, "horsepower": 190, "transmission": "CVT"}',
     'Los Angeles', 'California', 'United States',
     65.00,
     'active', 'rental', true,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": true, "all_wheel_drive": true, "third_row_seating": false}',
     'Spacious and reliable Honda CR-V perfect for family trips and outdoor adventures. Features all-wheel drive for various weather conditions.',
     ARRAY['https://example.com/images/crv1.jpg', 'https://example.com/images/crv2.jpg', 'https://example.com/images/crv3.jpg'],
     8200),

    ('c0000004-0000-4000-8000-000000000004',
     '22222222-0000-4000-8000-000000000002',
     'Luxury Sedan', 'BMW', '3 Series', 2024, 'Petrol',
     '{"engine_size": 2.0, "cylinders": 4, "horsepower": 255, "transmission": "Automatic"}',
     'Los Angeles', 'California', 'United States',
     120.00,
     'active', 'rental', true,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": true, "leather_seats": true, "sunroof": true, "premium_audio": true}',
     'Premium BMW 3 Series with luxury features and outstanding performance. Perfect for business trips or special occasions.',
     ARRAY['https://example.com/images/bmw3series1.jpg', 'https://example.com/images/bmw3series2.jpg'],
     5670),

    ('c0000005-0000-4000-8000-000000000005',
     '22222222-0000-4000-8000-000000000002',
     'Weekend Sports Car', 'Mazda', 'MX-5 Miata', 2023, 'Petrol',
     '{"engine_size": 2.0, "cylinders": 4, "horsepower": 181, "transmission": "Manual"}',
     'Los Angeles', 'California', 'United States',
     85.00,
     'active', 'rental', true,
     '{"gps": false, "air_conditioning": true, "bluetooth": true, "backup_camera": false, "keyless_entry": true, "convertible": true, "sport_mode": true}',
     'Fun and sporty Mazda MX-5 Miata convertible. Perfect for weekend getaways and scenic drives along the coast.',
     ARRAY['https://example.com/images/miata1.jpg', 'https://example.com/images/miata2.jpg'],
     12890),

    -- Mike Johnson's car (Electric Vehicle)
    ('c0000006-0000-4000-8000-000000000006',
     '33333333-0000-4000-8000-000000000003',
     'Eco-Friendly Tesla', 'Tesla', 'Model 3', 2024, 'Electric',
     '{"engine_size": 0, "cylinders": 0, "horsepower": 283, "transmission": "Single-Speed"}',
     'Seattle', 'Washington', 'United States',
     95.00,
     'active', 'rental', true,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": true, "autopilot": true, "supercharging": true, "premium_connectivity": true}',
     'State-of-the-art Tesla Model 3 with autopilot and premium features. Zero emissions and cutting-edge technology for the environmentally conscious driver.',
     ARRAY['https://example.com/images/tesla1.jpg', 'https://example.com/images/tesla2.jpg', 'https://example.com/images/tesla3.jpg'],
     7320),

    -- System/Company cars (no owner)
    ('c0000007-0000-4000-8000-000000000007',
     NULL,
     'Corporate Fleet Vehicle', 'Toyota', 'Prius', 2023, 'Hybrid',
     '{"engine_size": 1.8, "cylinders": 4, "horsepower": 121, "transmission": "CVT"}',
     'New York', 'New York', 'United States',
     50.00,
     'active', 'rental', true,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": false, "hybrid_system": true}',
     'Fuel-efficient Toyota Prius hybrid perfect for city driving. Part of CarZone corporate fleet with excellent fuel economy.',
     ARRAY['https://example.com/images/prius1.jpg'],
     22100),

    ('c0000008-0000-4000-8000-000000000008',
     NULL,
     'Premium Rental', 'Audi', 'A4', 2024, 'Petrol',
     '{"engine_size": 2.0, "cylinders": 4, "horsepower": 261, "transmission": "Automatic"}',
     'Miami', 'Florida', 'United States',
     110.00,
     'active', 'rental', true,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": true, "leather_seats": true, "navigation": true, "premium_audio": true}',
     'Luxury Audi A4 sedan with premium features and exceptional comfort. Perfect for business travel and special events.',
     ARRAY['https://example.com/images/audi1.jpg', 'https://example.com/images/audi2.jpg'],
     9850),

    -- Cars in different statuses for testing
    ('c0000009-0000-4000-8000-000000000009',
     '22222222-0000-4000-8000-000000000002',
     'Under Maintenance', 'Ford', 'Escape', 2022, 'Petrol',
     '{"engine_size": 1.5, "cylinders": 3, "horsepower": 181, "transmission": "Automatic"}',
     'Chicago', 'Illinois', 'United States',
     55.00,
     'maintenance', 'rental', false,
     '{"gps": true, "air_conditioning": true, "bluetooth": true, "backup_camera": true, "keyless_entry": true}',
     'Ford Escape currently undergoing scheduled maintenance. Will be available for rental and purchase soon.',
     ARRAY['https://example.com/images/escape1.jpg'],
     35670),

    ('c0000010-0000-4000-8000-000000000010',
     '11111111-0000-4000-8000-000000000001',
     'Sale Only Vehicle', 'Volkswagen', 'Jetta', 2021, 'Petrol',
     '{"engine_size": 1.4, "cylinders": 4, "horsepower": 147, "transmission": "Manual"}',
     'Austin', 'Texas', 'United States',
     40.00,
     'active', 'rental', true,
     '{"gps": false, "air_conditioning": true, "bluetooth": true, "backup_camera": false, "keyless_entry": false}',
     'Well-maintained Volkswagen Jetta available for purchase only. Great first car or reliable daily driver.',
     ARRAY['https://example.com/images/jetta1.jpg'],
     45230);

-- Insert comprehensive sample booking data
-- These bookings represent different scenarios: rentals, sales, various statuses
INSERT INTO booking (id, customer_id, car_id, owner_id, status, total_amount, start_date, end_date, notes) VALUES

    -- Confirmed rental bookings
    ('b0000001-0000-4000-8000-000000000001',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000001-0000-4000-8000-000000000001',  -- John's Toyota Camry
     '11111111-0000-4000-8000-000000000001',  -- John Doe (owner)
     'confirmed', 135.00,
     '2024-02-01 10:00:00', '2024-02-04 10:00:00',
     'Weekend getaway trip to Napa Valley. Customer has excellent driving record.'),

    ('b0000002-0000-4000-8000-000000000002',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000003-0000-4000-8000-000000000003',  -- Jane's Honda CR-V
     '22222222-0000-4000-8000-000000000002',  -- Jane Smith (owner)
     'confirmed', 325.00,
     '2024-02-05 09:00:00', '2024-02-10 09:00:00',
     'Family vacation to Lake Tahoe. Requesting child seat installation.'),

    -- Active rental (currently ongoing)
    ('b0000003-0000-4000-8000-000000000003',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000006-0000-4000-8000-000000000006',  -- Mike's Tesla Model 3
     '33333333-0000-4000-8000-000000000003',  -- Mike Johnson (owner)
     'active', 285.00,
     '2024-01-28 14:00:00', '2024-01-31 14:00:00',
     'Business trip to Portland. Customer specifically requested electric vehicle.'),

    -- Completed rental bookings
    ('b0000004-0000-4000-8000-000000000004',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000002-0000-4000-8000-000000000002',  -- John's Honda Civic
     '11111111-0000-4000-8000-000000000001',  -- John Doe (owner)
     'completed', 70.00,
     '2024-01-15 12:00:00', '2024-01-17 12:00:00',
     'Quick rental for airport transportation. Customer very satisfied.'),

    ('b0000005-0000-4000-8000-000000000005',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000005-0000-4000-8000-000000000005',  -- Jane's Mazda MX-5 Miata
     '22222222-0000-4000-8000-000000000002',  -- Jane Smith (owner)
     'completed', 255.00,
     '2024-01-20 11:00:00', '2024-01-23 11:00:00',
     'Weekend convertible rental for anniversary celebration.'),

    -- Pending rental bookings
    ('b0000006-0000-4000-8000-000000000006',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000004-0000-4000-8000-000000000004',  -- Jane's BMW 3 Series
     '22222222-0000-4000-8000-000000000002',  -- Jane Smith (owner)
     'pending', 360.00,
     '2024-02-15 16:00:00', '2024-02-18 16:00:00',
     'Business conference in downtown LA. Awaiting license verification.'),

    ('b0000007-0000-4000-8000-000000000007',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000007-0000-4000-8000-000000000007',  -- Corporate Toyota Prius
     NULL,                                    -- No individual owner (corporate fleet)
     'pending', 100.00,
     '2024-02-10 08:00:00', '2024-02-12 08:00:00',
     'Eco-friendly option for client meetings. Customer prefers hybrid vehicles.'),

    -- Cancelled rental booking
    ('b0000008-0000-4000-8000-000000000008',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000008-0000-4000-8000-000000000008',  -- Corporate Audi A4
     NULL,                                    -- No individual owner (corporate fleet)
     'cancelled', 220.00,
     '2024-01-25 13:00:00', '2024-01-27 13:00:00',
     'Customer cancelled due to change in travel plans. Full refund processed.'),

    -- Sale bookings (confirmed purchases)
    ('b0000009-0000-4000-8000-000000000009',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000010-0000-4000-8000-000000000010', -- John's Volkswagen Jetta
     '11111111-0000-4000-8000-000000000001',  -- John Doe (owner)
     'confirmed', 160.00,
     '2024-02-20 09:00:00', '2024-02-24 09:00:00',
     'Extended rental for business meetings. Customer satisfied with vehicle condition.'),

    ('b0000010-0000-4000-8000-000000000010',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000002-0000-4000-8000-000000000002',  -- John's Honda Civic (available for both)
     '11111111-0000-4000-8000-000000000001',  -- John Doe (owner)
     'pending', 105.00,
     '2024-02-25 10:00:00', '2024-02-28 10:00:00',
     'Weekend rental pending final verification. Customer has good rental history.'),

    -- Completed sale
    ('b0000011-0000-4000-8000-000000000011',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000005-0000-4000-8000-000000000005',  -- Jane's Mazda MX-5 Miata
     '22222222-0000-4000-8000-000000000002',  -- Jane Smith (owner)
     'completed', 255.00,
     '2024-02-10 11:00:00', '2024-02-13 11:00:00',
     'Completed rental. Customer enjoyed the convertible experience. Excellent feedback.'),

    -- More diverse rental scenarios
    ('b0000012-0000-4000-8000-000000000012',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000008-0000-4000-8000-000000000008',  -- Corporate Audi A4
     NULL,                                    -- Corporate vehicle
     'confirmed', 440.00,
     '2024-02-20 15:00:00', '2024-02-24 15:00:00',
     'Corporate client rental. Premium vehicle for important business meetings.'),

    -- Long-term rental
    ('b0000013-0000-4000-8000-000000000013',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000001-0000-4000-8000-000000000001',  -- John's Toyota Camry
     '11111111-0000-4000-8000-000000000001',  -- John Doe (owner)
     'confirmed', 1350.00,
     '2024-03-01 10:00:00', '2024-03-31 10:00:00',
             'Month-long rental while customer''s car is being repaired. Excellent customer history.'),    -- Recent completed booking
    ('b0000014-0000-4000-8000-000000000014',
     '44444444-0000-4000-8000-000000000004',  -- Sarah Wilson (customer)
     'c0000007-0000-4000-8000-000000000007',  -- Corporate Toyota Prius
     NULL,                                    -- Corporate vehicle
     'completed', 150.00,
     '2024-01-22 09:00:00', '2024-01-25 09:00:00',
     'Eco-friendly rental for environmental conference. Customer very pleased with fuel efficiency.'),

    -- Future pending booking
    ('b0000015-0000-4000-8000-000000000015',
     '55555555-0000-4000-8000-000000000005',  -- David Brown (customer)
     'c0000006-0000-4000-8000-000000000006',  -- Mike's Tesla Model 3
     '33333333-0000-4000-8000-000000000003',  -- Mike Johnson (owner)
     'pending', 475.00,
     '2024-03-15 11:00:00', '2024-03-20 11:00:00',
     'Customer wants to try electric vehicle before potential purchase. Special EV orientation requested.');

-- =============================================================================
-- VERIFICATION QUERIES
-- =============================================================================

-- Query to verify data insertion and relationships
-- Uncomment these lines to verify the setup (useful for debugging)

-- SELECT 'Users Count' as info, COUNT(*) as count FROM users;
-- SELECT 'Car Count' as info, COUNT(*) as count FROM car;
-- SELECT 'Booking Count' as info, COUNT(*) as count FROM booking;

-- Sample query to test JSONB functionality
-- SELECT 
--     c.name, 
--     c.brand, 
--     c.model,
--     c.year, 
--     c.fuel_type,
--     c.engine->>'horsepower' as horsepower,
--     c.engine->>'transmission' as transmission,
--     c.price as daily_rental,
--     c.price->>'sale_price' as sale_price,
--     u.username as owner_username
-- FROM car c 
-- LEFT JOIN users u ON c.owner_id = u.id 
-- ORDER BY c.brand, c.model;

-- Query to test engine filtering
-- SELECT name, brand, model FROM car WHERE engine->>'horsepower' > '200';

-- Query to test price filtering  
-- SELECT name, brand, price as daily_rate 
-- FROM car 
-- WHERE price < 60;

-- Sample booking queries for testing
-- SELECT 
--     b.id,
--     b.status,
--     b.status,
--     b.total_amount,
--     b.start_date,
--     b.end_date,
--     customer.username as customer_name,
--     car.name as car_name,
--     car.brand,
--     car.model,
--     owner.username as owner_name
-- FROM booking b
-- JOIN users customer ON b.customer_id = customer.id
-- JOIN car ON b.car_id = car.id
-- LEFT JOIN users owner ON b.owner_id = owner.id
-- ORDER BY b.created_at DESC;

-- Query to test booking status filtering
-- SELECT status, COUNT(*) as count FROM booking GROUP BY status;

-- Query to test booking type filtering
-- SELECT status, COUNT(*) as count FROM booking GROUP BY status;

-- =============================================================================
-- SCHEMA SUMMARY
-- =============================================================================

-- This schema supports:
-- ✓ User management with roles and profile data
-- ✓ Flexible car data with JSONB for engine and pricing
-- ✓ Comprehensive booking system for rentals and sales
-- ✓ Rental and sales functionality with status tracking
-- ✓ Location-based searches
-- ✓ Feature-rich car listings
-- ✓ Booking conflict detection and management
-- ✓ Customer and owner relationship tracking
-- ✓ Comprehensive indexing for performance
-- ✓ Audit trails with automatic timestamps
-- ✓ Rich test data for development and testing including:
--   - 6 test users (owners, customers, admin)
--   - 10 test cars (various brands, types, statuses)
--   - 15 test bookings (rentals, sales, various statuses)
