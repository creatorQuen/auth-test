BEGIN;

CREATE TABLE IF NOT EXISTS auth_users (
    id serial primary key not null ,
    created_at timestamp without time zone not null,
    email varchar(255) not null unique,
    password_hash varchar(255) not null,
    login varchar(255) not null unique,
    phone varchar(30) not null
    );

COMMIT;