-- init.sql

CREATE SCHEMA IF NOT EXISTS yt_convertor;

CREATE TABLE IF NOT EXISTS yt_convertor.jobs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    youtube_url TEXT NOT NULL,
    start_time INTERVAL NOT NULL,
    end_time INTERVAL NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file_url TEXT
);