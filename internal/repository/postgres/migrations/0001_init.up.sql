CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL DEFAULT ''::VARCHAR,
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    access_level INTEGER DEFAULT 1 NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms(
    id SERIAL PRIMARY KEY,
    room_name VARCHAR(255) NOT NULL,
    day_price  FLOAT NOT NULL DEFAULT 0,       
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS reservations(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) DEFAULT '',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS restrictions(
    id SERIAL PRIMARY KEY,
    restriction_name VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS room_restrictions(
    id SERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE ON UPDATE CASCADE,
    reservation_id INTEGER REFERENCES reservations(id) ON DELETE CASCADE ON UPDATE CASCADE,
    restriction_id INTEGER NOT NULL REFERENCES restrictions(id),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS services(
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255),
    day_price  FLOAT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS user_services(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    service_id INTEGER REFERENCES services(id),
    discount INTEGER,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

CREATE INDEX room_restrictions_reservation_id_idx ON room_restrictions (reservation_id);
CREATE INDEX room_restrictions_room_id_idx ON room_restrictions (room_id);
CREATE INDEX room_restrictions_start_date_end_date_idx ON room_restrictions (start_date,end_date);
CREATE INDEX room_restrictions_email_idx ON reservations (email);
CREATE INDEX room_restrictions_last_name_idx ON reservations (last_name);