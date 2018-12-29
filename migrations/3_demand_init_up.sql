CREATE TABLE ino.demand_school (
    rider_id    INT NOT NULL,
    day         INT NOT NULL,
    dir         INT NOT NULL,
    start       INT NOT NULL,
    end         INT NOT NULL,
    FOREIGN KEY(rider_id) REFERENCES riders (id)
);

CREATE TABLE ino.demand_aggregate (
    time_zone   INT NOT NULL,
    value       INT NOT NULL
);