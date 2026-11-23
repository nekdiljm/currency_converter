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
    currency_from       VARCHAR(10)     NOT NULL,
    currency_to         VARCHAR(10)     NOT NULL,
    amount              NUMERIC         NOT NULL,
    result              NUMERIC         NOT NULL,
    rate                NUMERIC         NOT NULL,
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS favorite_pairs (
    id                  SERIAL        PRIMARY KEY,
    user_id             INT           REFERENCES users(id),
    currency_from       VARCHAR(10)   NOT NULL,
    currency_to         VARCHAR(10)   NOT NULL,

    UNIQUE (user_id, currency_from, currency_to)
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id              SERIAL            PRIMARY KEY,
    currency_from   VARCHAR(10)       NOT NULL,
    currency_to     VARCHAR(10)       NOT NULL,
    rate            NUMERIC           NOT NULL,
    date            TIMESTAMP         NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP         NOT NULL DEFAULT NOW(),

    UNIQUE (currency_from, currency_to, date)
);


