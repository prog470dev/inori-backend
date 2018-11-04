CREATE DATABASE ino;
USE ino;

CREATE TABLE drivers (
    id INT NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    grade VARCHAR(255) NOT NULL,
    major VARCHAR(255) NOT NULL,
    mail VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    car_color VARCHAR(255) NOT NULL,
    car_number VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE riders (
    id INT NOT NULL AUTO_INCREMENT,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    grade VARCHAR(255) NOT NULL,
    major VARCHAR(255) NOT NULL,
    mail VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE offers (
    id INT NOT NULL AUTO_INCREMENT,
    driver_id INT NOT NULL,
    start VARCHAR(255) NOT NULL,
    goal VARCHAR(255) NOT NULL,
    departure_time DATETIME NOT NULL,
    rider_capacity INT NOT NULL,
    FOREIGN KEY(driver_id) REFERENCES drivers (id),
    PRIMARY KEY (id)
);

CREATE TABLE reservations (
    id INT NOT NULL  AUTO_INCREMENT,
    offer_id INT NOT NULL,
    rider_id INT NOT NULL,
    FOREIGN KEY(offer_id) REFERENCES offers (id),
    FOREIGN KEY(rider_id) REFERENCES riders (id),
    PRIMARY KEY (id)
);
