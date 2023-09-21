CREATE OR REPLACE FUNCTION haversine_distance(
    lat1 FLOAT,
    lon1 FLOAT,
    lat2 FLOAT,
    lon2 FLOAT
) RETURNS FLOAT AS $$
DECLARE
    radius FLOAT;
    dlat FLOAT;
    dlon FLOAT;
    a FLOAT;
    c FLOAT;
BEGIN
    radius := 6371; -- Earth's radius in kilometers

    dlat := RADIANS(lat2 - lat1);
    dlon := RADIANS(lon2 - lon1);

    a := SIN(dlat / 2) * SIN(dlat / 2) + COS(RADIANS(lat1)) * COS(RADIANS(lat2)) * SIN(dlon / 2) * SIN(dlon / 2);
    c := 2 * ATAN2(SQRT(a), SQRT(1 - a));

    RETURN radius * c; -- Distance in kilometers
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS bus_lines (
    id varchar,
    external_id varchar,
    full_name varchar,
    short_name varchar,
    origin varchar,

    PRIMARY KEY(id),
    UNIQUE (external_id)
);

CREATE TABLE IF NOT EXISTS bus_stops (
    id varchar,
    external_id varchar,
    name varchar,
    latitude FLOAT,
    longitude FLOAT,

    PRIMARY KEY(id),
    UNIQUE (external_id)
);

CREATE TABLE IF NOT EXISTS bus_line_paths (
    id varchar,
    sequence int,
    bus_line_id VARCHAR,
    latitude FLOAT,
    longitude FLOAT,
    distance_to_next_path FLOAT,

    PRIMARY KEY(id),
    FOREIGN KEY(bus_line_id) REFERENCES bus_lines(id)
);

CREATE TABLE IF NOT EXISTS bus_line_stops (
    id varchar,
    bus_line_id VARCHAR,
    nearest_bus_line_path_id VARCHAR,
    bus_stop_id VARCHAR,

    PRIMARY KEY(id),
    FOREIGN KEY (bus_line_id) REFERENCES bus_lines(id),
    FOREIGN KEY(nearest_bus_line_path_id) REFERENCES bus_line_paths(id),
    FOREIGN KEY(bus_stop_id) REFERENCES bus_stops(id)
);

CREATE TABLE IF NOT EXISTS bus_informations (
    id varchar,
    plat_number varchar,
    current_line_id VARCHAR,
    last_location_id varchar,

    PRIMARY KEY(id),
    UNIQUE (plat_number),
    FOREIGN KEY (current_line_id) REFERENCES bus_lines(id)
);

CREATE TABLE IF NOT EXISTS bus_location_histories (
    id varchar,
    bus_id VARCHAR,
    latitude FLOAT,
    longitude FLOAT,
    bearing FLOAT,
    crowd_level varchar,

    PRIMARY KEY(id)
);
