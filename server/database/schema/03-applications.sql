/**
 * Table that holds the details of various application that are registered
 * and can use the APIs.
 */
CREATE TABLE applications
(
    id                     TEXT                     NOT NULL PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
    name                   TEXT                     NOT NULL,

    -- Address of the wallet associated to the application.
    --
    -- This address is going to be used to grant fee allowances on behalf of the application.
    -- For this reason, it must have given an on-chain authorization to the APIs in order to sign
    -- MsgGrantAllowance transactions on its behalf.
    -- If its not set, or an on-chain authorization does not exist, the application will not be
    -- able to request fee grants.
    wallet_address         TEXT,

    -- Whether the application can send notifications to the users or not
    can_send_notifications BOOLEAN                  NOT NULL                    DEFAULT FALSE,

    -- Time at which the application was created
    creation_time          TIMESTAMP WITH TIME ZONE NOT NULL                    DEFAULT NOW()
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
