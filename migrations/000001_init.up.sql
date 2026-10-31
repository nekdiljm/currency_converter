CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    first_name  VARCHAR(100)    NOT NULL,
    last_name   VARCHAR(100)    NOT NULL,
    email       VARCHAR(200)    NOT NULL UNIQUE,
    password    VARCHAR(300)    NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS currencies (
    id      SERIAL PRIMARY KEY,
    code    VARCHAR(3)  NOT NULL UNIQUE,
    name    VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS conversion_history (
    id                  SERIAL PRIMARY KEY,
    user_id             INT             REFERENCES users(id),
    currency_from_id    INT             REFERENCES currencies(id),
    currency_to_id      INT             REFERENCES currencies(id),
    amount              DECIMAL(10,5)   NOT NULL,
    result              DECIMAL(10,5)   NOT NULL,
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS favorite_pairs (
    id                  SERIAL PRIMARY KEY,
    user_id             INT REFERENCES users(id),
    currency_from_id    INT REFERENCES currencies(id),
    currency_to_id      INT REFERENCES currencies(id)
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id                  SERIAL PRIMARY KEY,
    currency_from_id    INT             REFERENCES currencies(id),
    currency_to_id      INT             REFERENCES currencies(id),
    conv_rate           DECIMAL(10,5)   NOT NULL,
    updated_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);
