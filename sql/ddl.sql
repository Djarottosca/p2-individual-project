DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS properties;

DROP TABLE IF EXISTS users;

CREATE DATABASE property_market;

-- users: yang masang & yang transaksi
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        full_name VARCHAR(100) NOT NULL,
        email VARCHAR(150) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        is_verified BOOLEAN NOT NULL DEFAULT FALSE,
        verification_token VARCHAR(255),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- properties: listing-nya
CREATE TABLE
    properties (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id),
        title VARCHAR(150) NOT NULL,
        property_type VARCHAR(20) NOT NULL, -- rumah | apartemen | tanah
        transaction_type VARCHAR(20) NOT NULL, -- dijual | disewakan
        price BIGINT NOT NULL,
        land_area INT, -- buat rumah & tanah
        building_area INT, -- buat apartemen
        bedrooms INT NOT NULL DEFAULT 0,
        bathrooms INT NOT NULL DEFAULT 0,
        certificate VARCHAR(20), -- SHM | HGB
        city VARCHAR(100) NOT NULL,
        district VARCHAR(100) NOT NULL,
        description TEXT,
        image_urls JSONB NOT NULL DEFAULT '[]',
        status VARCHAR(20) NOT NULL DEFAULT 'available', -- available | booked | sold
        featured_until TIMESTAMPTZ, -- kosong = nggak lagi featured
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- transactions: promote, booking, pelunasan
CREATE TABLE
    transactions (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id),
        property_id UUID NOT NULL REFERENCES properties (id),
        booking_id UUID REFERENCES transactions (id), -- diisi pas pelunasan, nunjuk booking-nya
        type VARCHAR(20) NOT NULL, -- promote | booking | pelunasan
        amount BIGINT NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending | paid | expired
        external_id VARCHAR(100) NOT NULL, -- id yang kita kirim ke xendit
        invoice_id VARCHAR(100), -- id balikan dari xendit
        invoice_url TEXT, -- link pembayaran
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );