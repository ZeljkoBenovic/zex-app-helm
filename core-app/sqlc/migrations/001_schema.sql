-- +migrate Up
CREATE TABLE assets (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title VARCHAR(255),
    aboutMe TEXT
);


-- +migrate Down
DROP TABLE assets;