CREATE TABLE IF NOT EXISTS doctors (
                                     doctorId SERIAL PRIMARY KEY,
                                     firstName VARCHAR(255) NOT NULL,
                                     secondName VARCHAR(255) NOT NULL,
                                     speciality VARCHAR(255) NOT NULL,
                                     phone VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS patients (
                                     patientId SERIAL PRIMARY KEY,
                                     firstName VARCHAR(255) NOT NULL,
                                     secondName VARCHAR(255) NOT NULL,
                                     phone VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments
(
    appointmentId               SERIAL PRIMARY KEY,
    created_at       timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at       timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    appointmentDate timestamp(0) with time zone NOT NULL,
    doctorId        bigserial,
    patientId       bigserial,
    FOREIGN KEY (doctorId)
        REFERENCES doctors(doctorId),
    FOREIGN KEY (patientId)
        REFERENCES patients(patientId)
);

CREATE TABLE IF NOT EXISTS model (
                                     modelId SERIAL PRIMARY KEY,
                                     doctorId INT NOT NULL,
                                     patientId INT NOT NULL,
                                     FOREIGN KEY (doctorId) REFERENCES doctors(doctorId),
                                     FOREIGN KEY (patientId) REFERENCES patients(patientId)
);
CREATE TABLE IF NOT EXISTS users (
                                     id bigserial PRIMARY KEY,
                                     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                     name text NOT NULL,
                                     email citext UNIQUE NOT NULL,
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
    ('doctor.create'),
    ('doctor.read'),
    ('doctor.update'),
    ('doctor.delete'),
    ('patient.create'),
    ('patient.read'),
    ('patient.update'),
    ('patient.delete');