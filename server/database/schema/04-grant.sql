/**
 * Table that holds the requests for on-chain fee grants.
 */
CREATE TABLE fee_grant_requests
(
    -- ID of the application that has requested the grant
    application_id  TEXT                     NOT NULL REFERENCES applications (id),

    -- Desmos address of the user who should receive the grant
    grantee_address TEXT                     NOT NULL PRIMARY KEY,

    -- Time at which the request was made
    request_time    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Time at which the authorizations were granted.
    -- This will be updated after the on-chain authorization has been created
    grant_time      TIMESTAMP WITH TIME ZONE
);