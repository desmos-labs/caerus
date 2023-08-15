/*************************************************************************/
/***                               Login                              ****/
/*************************************************************************/

CREATE TABLE nonces
(
    desmos_address  TEXT NOT NULL,
    token           TEXT NOT NULL UNIQUE, /* SHA-256 encrypted value */
    expiration_time TIMESTAMP WITH TIME ZONE
);

CREATE TABLE sessions
(
    desmos_address  TEXT                     NOT NULL,
    token           TEXT                     NOT NULL UNIQUE, /* SHA-256 encrypted value */
    creation_time   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expiration_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + INTERVAL '1 week'
);

/*************************************************************************/
/***                               Users                              ****/
/*************************************************************************/

CREATE TABLE users
(
    desmos_address TEXT                     NOT NULL PRIMARY KEY,
    creation_time  TIMESTAMP WITH TIME ZONE NOT NULL                       DEFAULT NOW(),
    last_login     TIMESTAMP WITH TIME ZONE NOT NULL                       DEFAULT NOW()
);