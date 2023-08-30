SELECT 'CREATE DATABASE segmentation' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'segmentation')\gexec
\c segmentation;

ALTER DATABASE segmentation SET timezone TO 'Europe/Moscow';

CREATE TABLE segments (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE segments_users (
    segments_id SERIAL NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (segments_id, user_id)
);

CREATE TABLE report (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    segments_id SERIAL NOT NULL,
    action VARCHAR(6),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
