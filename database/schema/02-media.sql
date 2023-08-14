/*
 * Contains all the hashes of the images generated using BlurHash (https://blurha.sh/).
 */
CREATE TABLE files_hashes
(
    file_name TEXT NOT NULL PRIMARY KEY,
    hash      TEXT NOT NULL
);