/**
 * Table that holds the details of the various subscription plans that can be
 * subscribed to by applications.
 */
CREATE TABLE application_subscriptions
(
    id                       SERIAL  NOT NULL PRIMARY KEY,

    -- Name of the subscription plan
    subscription             TEXT    NOT NULL,

    -- Number of fee grants that can be requested per day
    -- If set to 0, no limit is applied
    fee_grant_rate_limit     INTEGER NOT NULL,

    -- Number of notifications that can be sent per day
    -- If set to 0, no limit is applied
    notifications_rate_limit INTEGER NOT NULL
);

/**
 * Table that holds the details of various application that are registered
 * and can use the APIs.
 */
CREATE TABLE applications
(
    id              TEXT                     NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),

    -- Name of the application
    name            TEXT                     NOT NULL,

    -- Address of the wallet associated to the application.
    --
    -- This address is going to be used to grant fee allowances on behalf of the application.
    -- For this reason, it must have given an on-chain authorization to the APIs in order to sign
    -- MsgGrantAllowance transactions on its behalf.
    -- If its not set, or an on-chain authorization does not exist, the application will not be
    -- able to request fee grants.
    wallet_address  TEXT,

    -- ID of the subscription plan that the application is subscribed to
    subscription_id INTEGER REFERENCES application_subscriptions (id),

    -- Time at which the application was created
    creation_time   TIMESTAMP WITH TIME ZONE NOT NULL                    DEFAULT NOW()
);

/**
 * Table that holds the details of the various admins that can manage the
 * applications.
 */
CREATE TABLE application_admins
(
    id             SERIAL NOT NULL PRIMARY KEY,
    application_id TEXT   NOT NULL REFERENCES applications (id),
    user_address   TEXT   NOT NULL,
    CONSTRAINT unique_application_user_entry UNIQUE (application_id, user_address)
);

/**
 * Table that holds all the tokens that can be used to authenticate API requests
 * on behalf of different applications.
 */
CREATE TABLE application_tokens
(
    id             SERIAL                   NOT NULL PRIMARY KEY,

    -- ID of the application that owns the token
    application_id TEXT                     NOT NULL REFERENCES applications (id),

    -- Name of the token decided by the user
    token_name     TEXT                     NOT NULL,

    -- SHA-256 encrypted value of the token
    token_value    TEXT                     NOT NULL,

    -- Time at which the token was created
    creation_time  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

