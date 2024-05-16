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

CREATE TABLE IF NOT EXISTS appointments
(
    id              bigserial PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    doctor_id       bigserial                        NOT NULL,
    patient_id      bigserial                        NOT NULL,
    date_time       text                        NOT NULL,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id)
);

CREATE TABLE IF NOT EXISTS users (
                                     id bigserial PRIMARY KEY,
                                     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), name text NOT NULL,
                                     email text UNIQUE NOT NULL,
                                     password_hash bytea NOT NULL,
                                     activated bool NOT NULL,
                                     version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS tokens (
                                      hash bytea PRIMARY KEY,
                                      user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                      expiry timestamp(0) with time zone NOT NULL,
                                      scope text NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
                                           id bigserial PRIMARY KEY,
                                           code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
                                                 user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                                 permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
                                                 PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES
    ('user.create'),
    ('user.read'),
    ('user.update'),
    ('user.delete'),
    ('patient.create'),
    ('patient.read'),
    ('patient.update'),
    ('patient.delete'),
    ('doctor.create'),
    ('doctor.read'),
    ('doctor.update'),
    ('doctor.delete'),
    ('appointment.create'),
    ('appointment.read'),
    ('appointment.update'),
    ('appointment.delete');