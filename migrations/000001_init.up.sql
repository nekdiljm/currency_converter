CREATE TABLE IF NOT EXISTS users (
    id          SERIAL          PRIMARY KEY,
    first_name  VARCHAR(100)    NOT NULL,
    last_name   VARCHAR(100)    NOT NULL,
    email       VARCHAR(200)    NOT NULL UNIQUE,
    password    VARCHAR(300)    NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS conversion_history (
    id                  SERIAL          PRIMARY KEY,
    user_id             INT             REFERENCES users(id),
    currency_from       VARCHAR(3)      NOT NULL,
    currency_to         VARCHAR(3)      NOT NULL,
    amount              DECIMAL(10,5)   NOT NULL,
    result              DECIMAL(10,5)   NOT NULL,
    rate                DECIMAL(10,5),
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS favorite_pairs (
    id                  SERIAL        PRIMARY KEY,
    user_id             INT           REFERENCES users(id),
    currency_from       VARCHAR(3)    NOT NULL,
    currency_to         VARCHAR(3)    NOT NULL,

    UNIQUE (user_id, currency_from, currency_to)
);


