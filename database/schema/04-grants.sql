/**
 * Table that holds the requests for on-chain fee grants.
 */
CREATE TABLE fee_grant_requests
(
    id              TEXT                     NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ID of the application that has requested the grant
    application_id  TEXT                     NOT NULL REFERENCES applications (id) ON DELETE CASCADE,

    -- JSON-encoded allowance that will be granted to the user
    allowance       JSONB                    NOT NULL,

    -- Desmos address of the user who should receive the grant
    grantee_address TEXT                     NOT NULL,

    -- Time at which the request was made
    request_time    TIMESTAMP WITH TIME ZONE NOT NULL             DEFAULT NOW(),

    -- Time at which the authorizations were granted.
    -- This will be updated after the on-chain authorization has been created
    grant_time      TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_application_user_fee_grant_request UNIQUE (application_id, grantee_address)
);