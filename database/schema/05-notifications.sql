/**
 * Table that holds the notification tokens that can be used in order to send
 * push notifications to the users.
 */
CREATE TABLE user_notifications_tokens
(
    id           SERIAL                   NOT NULL PRIMARY KEY,
    user_address TEXT                     NOT NULL,
    device_token TEXT                     NOT NULL,
    timestamp    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_notification_token UNIQUE (user_address, device_token)
);

/**
 * Table that holds the notifications that have been sent to the users.
 */
CREATE TABLE notifications
(
    id             TEXT                     NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ID of the application that has sent the notification
    application_id TEXT                     NOT NULL,

    -- JSON-encoded list of addresses of the users that has received the notification
    user_addresses JSONB                    NOT NULL,

    -- JSON-encoded notification that has been sent
    notification   JSONB                    NOT NULL,

    -- Time when the notification has been sent
    send_time      TIMESTAMP WITH TIME ZONE NOT NULL             DEFAULT NOW()
);