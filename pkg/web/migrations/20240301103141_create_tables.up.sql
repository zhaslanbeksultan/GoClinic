CREATE TABLE IF NOT EXISTS doctors
(
    id          bigserial PRIMARY KEY,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    first_name  text                        NOT NULL,
    last_name   text                        NOT NULL,
    speciality  text                        NOT NULL,
    phone       text                        NOT NULL
);

CREATE TABLE IF NOT EXISTS patients
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    first_name      text                        NOT NULL,
    last_name       text                        NOT NULL,
    phone           text                        NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
                                     id bigserial PRIMARY KEY,
                                     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), name text NOT NULL,
                                     email text UNIQUE NOT NULL,
                                     password_hash bytea NOT NULL,
                                     activated bool NOT NULL,
                                     version integer NOT NULL DEFAULT 1
);