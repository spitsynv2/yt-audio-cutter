-- init.sql

CREATE SCHEMA IF NOT EXISTS yt_convertor;

CREATE TABLE IF NOT EXISTS yt_convertor.jobs (
    id TEXT PRIMARY KEY,
    youtube_url TEXT NOT NULL,
    start_time INTERVAL NOT NULL,
    end_time INTERVAL NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS yt_convertor.files (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    file_url TEXT NOT NULL,
    file_size BIGINT NOT NULL,  -- size in bytes
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_task
      FOREIGN KEY (task_id)
      REFERENCES yt_convertor.jobs(id)
      ON DELETE CASCADE
);