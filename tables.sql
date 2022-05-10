CREATE TABLE users (
    id serial primary key,
    username varchar(100) unique,
    password varchar(100),
    created_at timestamp default now(),
    updated_at timestamp default now()
);

CREATE TABLE companies(
    id serial primary key,
    name varchar(100) not null,
    code varchar(50) not null,
    country varchar(20) not null,
    website varchar(100) not null,
    phone varchar(50) not null,
    status varchar(20) not null default 'active',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);