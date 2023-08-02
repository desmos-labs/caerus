CREATE TABLE notifications_tokens
(
    user_address TEXT                     NOT NULL,
    device_token TEXT                     NOT NULL,
    timestamp    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_notification_token UNIQUE (user_address, device_token)
);

/**
 * This table is used to store the notification senders, which are all the users that are allowed to send notifications.
 * The content of this table will only be used if authentication for notifications sending is enabled.
 */
CREATE TABLE notification_senders
(
    /*
     * Token used to authenticate the sender.
     * This will have to be sent inside the Authentication header of the request to send a notification
     */
    token TEXT NOT NULL PRIMARY KEY,

    /*
     * Human readable name for the sender.
     * This is actually never used and it's present only to identify the ownership of different tokens
     */
    name  TEXT
);