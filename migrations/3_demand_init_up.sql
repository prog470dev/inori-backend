CREATE TABLE ino.demand_school (
    rider_id    INT NOT NULL,
    day         INT NOT NULL,
    start       INT NOT NULL,
    end         INT NOT NULL,
    FOREIGN KEY(rider_id) REFERENCES riders (id)
);

CREATE TABLE ino.demand_home (
    rider_id    INT NOT NULL,
    day         INT NOT NULL,
    start       INT NOT NULL,
    end         INT NOT NULL,
    FOREIGN KEY(rider_id) REFERENCES riders (id)
);

CREATE TABLE ino.demand_aggregate_school (
    time_zone   INT NOT NULL,
    value       INT NOT NULL
);

CREATE TABLE ino.demand_aggregate_home (
    time_zone   INT NOT NULL,
    value       INT NOT NULL
);