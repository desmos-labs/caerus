/**
 * Table that holds the deep links that have been created by applications.
 */
CREATE TABLE deep_links
(
    id             TEXT                     NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ID of the application that has generated the link
    application_id TEXT                     NOT NULL,

    -- Link that has been generated
    link_url       TEXT                     NOT NULL,

    -- Configuration used to generate the link
    link_config    JSONB                    NOT NULL,

    -- Time when the link was created
    creation_time  TIMESTAMP WITH TIME ZONE NOT NULL             DEFAULT NOW()
)