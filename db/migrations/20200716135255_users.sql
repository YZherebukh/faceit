-- migrate:up
CREATE TABLE users (
    user_id SERIAL,
    first_name varchar(50),
    last_name varchar(50),
    nick_name varchar(100),
    email varchar(100),
    country INTEGER REFERENCES countries(country_id),
    PRIMARY KEY (user_id)
);

-- migrate:down

DROP TABLE users;
