/**
 * Table that holds the notification tokens that can be used in order to send
 * push notifications to the users.
 */
CREATE TABLE notifications_tokens
(
    user_address TEXT                     NOT NULL,
    device_token TEXT                     NOT NULL,
    timestamp    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_notification_token UNIQUE (user_address, device_token)
);