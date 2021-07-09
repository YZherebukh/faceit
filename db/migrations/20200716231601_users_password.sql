-- migrate:up

CREATE TABLE users_password (
    password_id INTEGER REFERENCES users(user_id),
    pwd varchar(100),
    salt varchar(10),
    PRIMARY KEY (password_id)  
);

-- migrate:down

DROP TABLE users_password;

