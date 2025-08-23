-- CarZone Database Schema Definition
-- This file contains the complete database schema for the CarZone application
-- including table definitions, relationships, constraints, and sample data.

-- =============================================================================
-- CLEANUP AND PREPARATION
-- =============================================================================

-- Remove existing foreign key constraint to allow table recreation
-- IF EXISTS prevents errors if the constraint doesn't exist
ALTER TABLE IF EXISTS car 
DROP CONSTRAINT IF EXISTS fk_engine_id;

-- Clear existing data to ensure clean slate for schema setup
-- TRUNCATE is faster than DELETE for removing all rows
TRUNCATE TABLE IF EXISTS car CASCADE;
TRUNCATE TABLE IF EXISTS engine CASCADE;

-- Drop existing tables if they exist (for complete reset)
DROP TABLE IF EXISTS car;
DROP TABLE IF EXISTS engine;

-- =============================================================================
-- TABLE DEFINITIONS
-- =============================================================================

-- Engine Table Definition
-- Stores engine specifications and technical details
-- This table is referenced by the car table (master table in the relationship)
CREATE TABLE engine (
    -- Primary key: Unique identifier for each engine
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Engine technical specifications
    displacement INT NOT NULL CHECK (displacement > 0),           -- Engine displacement in cubic centimeters (cc)
    no_of_cylinders INT NOT NULL CHECK (no_of_cylinders > 0),     -- Number of cylinders (e.g., 4, 6, 8)
    car_range INT NOT NULL CHECK (car_range > 0),                 -- Vehicle range in kilometers
    
    -- Audit trail columns for tracking changes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,               -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP                -- Last update timestamp
);

-- Car Table Definition  
-- Stores car information including specifications, pricing, and engine association
-- This table has a foreign key relationship with the engine table
CREATE TABLE car (
    -- Primary key: Unique identifier for each car
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Car basic information
    name VARCHAR(255) NOT NULL,                                   -- Car model name (e.g., "Camry", "Model S")
    year VARCHAR(4) NOT NULL CHECK (year ~ '^[0-9]{4}$'),       -- Manufacturing year (4-digit string)
    brand VARCHAR(255) NOT NULL,                                  -- Manufacturer brand (e.g., "Toyota", "Tesla")
    fuel_type VARCHAR(50) NOT NULL,                              -- Fuel type (Petrol, Diesel, Electric, Hybrid)
    
    -- Foreign key relationship to engine table
    engine_id UUID NOT NULL,                                      -- Reference to engine.id
    
    -- Pricing information
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),           -- Car price with 2 decimal precision
    
    -- Audit trail columns
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,              -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP               -- Last update timestamp
);

-- =============================================================================
-- CONSTRAINTS AND RELATIONSHIPS
-- =============================================================================

-- Foreign Key Constraint: Establish relationship between car and engine
-- This ensures referential integrity and prevents orphaned records
ALTER TABLE car
ADD CONSTRAINT fk_engine_id
FOREIGN KEY (engine_id)
REFERENCES engine(id)
ON DELETE CASCADE;                                               -- Delete cars when associated engine is deleted

-- =============================================================================
-- INDEXES FOR PERFORMANCE
-- =============================================================================

-- Index on car brand for fast brand-based queries
CREATE INDEX idx_car_brand ON car(brand);

-- Index on car year for year-based filtering
CREATE INDEX idx_car_year ON car(year);

-- Index on engine displacement for engine specification queries
CREATE INDEX idx_engine_displacement ON engine(displacement);

-- Index on car price for price-range queries
CREATE INDEX idx_car_price ON car(price);

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
CREATE TRIGGER update_engine_updated_at 
    BEFORE UPDATE ON engine 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_car_updated_at 
    BEFORE UPDATE ON car 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- SAMPLE DATA FOR TESTING AND DEVELOPMENT
-- =============================================================================

-- Insert sample engine data
-- These engines represent different types: economy, performance, and electric
INSERT INTO engine (id, displacement, no_of_cylinders, car_range) VALUES
    -- Economy 4-cylinder engine
    ('e1f86b1a-0873-4c19-bae2-fc60329d0140', 2000, 4, 600),
    
    -- Compact 4-cylinder engine  
    ('f4a9c66b-8e38-419b-93c4-215d5cefb318', 1600, 4, 550),
    
    -- Performance V6 engine
    ('cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 3000, 6, 700),
    
    -- Small displacement 4-cylinder
    ('9746be12-07b7-42a3-b8ab-7d1f209b63d7', 1800, 4, 500),
    
    -- Electric motor (0 displacement, 0 cylinders, high range)
    ('a1b2c3d4-e5f6-7890-abcd-123456789012', 0, 0, 800);

-- Insert sample car data
-- These cars represent different brands, fuel types, and price ranges
INSERT INTO car (id, name, year, brand, fuel_type, engine_id, price) VALUES
    -- Toyota vehicles
    ('550e8400-e29b-41d4-a716-446655440000', 'Camry', '2024', 'Toyota', 'Petrol', 'e1f86b1a-0873-4c19-bae2-fc60329d0140', 28999.99),
    ('550e8400-e29b-41d4-a716-446655440001', 'Corolla', '2023', 'Toyota', 'Petrol', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', 24999.99),
    
    -- Honda vehicles  
    ('550e8400-e29b-41d4-a716-446655440002', 'Accord', '2024', 'Honda', 'Petrol', 'e1f86b1a-0873-4c19-bae2-fc60329d0140', 32999.99),
    ('550e8400-e29b-41d4-a716-446655440003', 'Civic', '2023', 'Honda', 'Petrol', 'f4a9c66b-8e38-419b-93c4-215d5cefb318', 26999.99),
    
    -- Luxury BMW vehicle
    ('550e8400-e29b-41d4-a716-446655440004', 'X5', '2024', 'BMW', 'Petrol', 'cc2c2a7d-2e21-4f59-b7b8-bd9e5e4cf04c', 65999.99),
    
    -- Electric Tesla vehicle
    ('550e8400-e29b-41d4-a716-446655440005', 'Model 3', '2024', 'Tesla', 'Electric', 'a1b2c3d4-e5f6-7890-abcd-123456789012', 45999.99);

-- =============================================================================
-- VERIFICATION QUERIES
-- =============================================================================

-- Query to verify data insertion and relationships
-- Uncomment these lines to verify the setup (useful for debugging)

-- SELECT 'Engine Count' as info, COUNT(*) as count FROM engine;
-- SELECT 'Car Count' as info, COUNT(*) as count FROM car;  
-- SELECT c.name, c.brand, c.year, c.fuel_type, c.price, e.displacement, e.no_of_cylinders, e.car_range
-- FROM car c 
-- JOIN engine e ON c.engine_id = e.id 
-- ORDER BY c.brand, c.name;