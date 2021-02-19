-- таблица зарегистрированных
CREATE TABLE users
(
    id        BIGSERIAL PRIMARY KEY,
    name      TEXT      NOT NULL,
    phone     TEXT      NOT NULL UNIQUE,
    password  TEXT      NOT NULL,
    active    BOOLEAN   NOT NULL DEFAULT TRUE,
    created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- таблица токенов зарегистрированных
CREATE TABLE users_tokens (
    token       TEXT      NOT NULL UNIQUE,
    user_id BIGINT    NOT NULL references users,
    expire      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO users (name, phone, password)
VALUES ('Hakim', '+992000000001', '123456789');


INSERT INTO users (name, phone, password)
VALUES ('Eydzhpee', '+992000000002', '23345');