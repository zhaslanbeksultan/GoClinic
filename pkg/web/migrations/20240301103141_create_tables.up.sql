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
    phone           int                         NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments
(
    id               bigserial PRIMARY KEY,
    created_at       timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at       timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    appointment_date timestamp(0) with time zone NOT NULL,
    doctor_id        bigserial,
    patient_id       bigserial,
    FOREIGN KEY (doctor_id)
        REFERENCES doctors(id),
    FOREIGN KEY (patient_id)
        REFERENCES patients(id)
);