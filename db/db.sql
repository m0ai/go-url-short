-- show current database
select current_database();

-- show current user
select current_user;

-- create user for shorturl database
CREATE USER shorturl_app WITH PASSWORD 'your-password-here';

-- create database schema for shorturl
CREATE DATABASE shorturl;
GRANT ALL PRIVILEGES ON DATABASE shorturl TO shorturl_app;

-- create table
CREATE TABLE shorturl.shorturl
(
    id         SERIAL PRIMARY KEY,
    url        VARCHAR(1024) NOT NULL,
    short      VARCHAR(16)   NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- show shorturl table owner
SELECT * from pg_tables WHERE tablename = 'shorturl';

-- chanege owner of table
ALTER TABLE shorturl OWNER TO shorturl_app;

-- truncate table
-- truncate table shorturl;
