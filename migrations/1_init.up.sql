CREATE TABLE
    IF NOT EXISTS jobs (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        cron TEXT NOT NULL,
        url TEXT NOT NULL,
        cron_id INT NOT NULL
    );