-- migrate:up
CREATE TABLE "user"
(
    id         BIGSERIAL PRIMARY KEY,
    username   VARCHAR(128)        NOT NULL,
    email      VARCHAR(256) UNIQUE NOT NULL,
    is_active  BOOLEAN   DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user_auth"
(
    id                    BIGSERIAL PRIMARY KEY,
    user_id               BIGINT       NOT NULL,
    hash                  VARCHAR(256) NOT NULL,
    is_locked             BOOLEAN   DEFAULT FALSE,
    failed_login_attempts SMALLINT  DEFAULT 0,
    last_failed_attempt   TIMESTAMP,
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX idx_user_auth_user_id ON "user_auth" USING HASH (user_id);

CREATE TABLE "user_profile"
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    first_name VARCHAR(128),
    last_name  VARCHAR(128),
    height     DECIMAL(3, 2),
    weight     DECIMAL(5, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX idx_user_profile_user_id ON "user_profile" USING HASH (user_id);

CREATE TABLE "session"
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    is_active  BOOLEAN   DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX idx_session_user_id ON "session" USING HASH (user_id);

-- migrate:down
DROP TABLE "user_profile";
DROP TABLE "user_auth";
DROP TABLE "session";
DROP TABLE "user";

