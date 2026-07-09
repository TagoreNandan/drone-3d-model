-- Enable the TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Control-plane tables ---------------------------------------------------

CREATE TABLE vehicles (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           TEXT NOT NULL,
    type           TEXT,
    firmware       TEXT,
    mavlink_sysid  INTEGER,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE missions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id        UUID REFERENCES vehicles(id),
    name              TEXT NOT NULL,
    plan_json         JSONB NOT NULL,
    geofence_json     JSONB,
    rally_points_json JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    uploaded_at       TIMESTAMPTZ
);

CREATE TABLE flight_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id   UUID REFERENCES vehicles(id),
    mission_id   UUID REFERENCES missions(id),
    started_at   TIMESTAMPTZ NOT NULL,
    ended_at     TIMESTAMPTZ,
    log_file_ref TEXT
);

-- Telemetry hypertable ----------------------------------------------------

CREATE TABLE telemetry_frames (
    vehicle_id   UUID NOT NULL REFERENCES vehicles(id),
    ts           TIMESTAMPTZ NOT NULL,
    lat          DOUBLE PRECISION,
    lon          DOUBLE PRECISION,
    alt          DOUBLE PRECISION,
    heading      DOUBLE PRECISION,
    groundspeed  DOUBLE PRECISION,
    battery_pct  SMALLINT,
    voltage      DOUBLE PRECISION,
    current      DOUBLE PRECISION,
    flight_mode  TEXT,
    armed        BOOLEAN,
    cpu_pct      REAL,
    ram_pct      REAL,
    disk_pct     REAL,
    cpu_temp     REAL,
    node_status  JSONB
);

SELECT create_hypertable('telemetry_frames', 'ts');
CREATE INDEX idx_telemetry_vehicle_ts ON telemetry_frames (vehicle_id, ts DESC);

-- Alerts (mirrors MAVLink STATUSTEXT) --------------------------------------

CREATE TABLE alerts (
    id          BIGSERIAL PRIMARY KEY,
    vehicle_id  UUID REFERENCES vehicles(id),
    ts          TIMESTAMPTZ NOT NULL DEFAULT now(),
    severity    SMALLINT,
    text        TEXT
);
