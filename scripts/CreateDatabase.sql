-- Creates a database and a user.
-- Parameters:
--   database_name: Name of the database to create
--   username: Database username to create
--   password: Password for the user

\c postgres
CREATE DATABASE :database_name;

CREATE USER :user_name WITH PASSWORD :'password';
GRANT CONNECT ON DATABASE :database_name TO :user_name;


\c :database_name
GRANT USAGE, CREATE ON SCHEMA public TO :user_name;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO :user_name;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO :user_name;


ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO :user_name;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO :user_name;