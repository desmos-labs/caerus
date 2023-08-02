/*
 * Contains all the hashes of the images generated using BlurHash (https://blurha.sh/).
 */
CREATE TABLE images_hashes
(
    image_url TEXT NOT NULL PRIMARY KEY,
    hash      TEXT NOT NULL
);