DROP TABLE IF EXISTS continents;
DROP TABLE IF EXISTS countries;

CREATE TABLE Continents (
    Continent_id serial,
    Name text NOT NULL,
    PRIMARY KEY (Continent_id)
);

CREATE TABLE Countries (
    Country_id serial,
    Name text NOT NULL,
    Population text NOT NULL,
    Capital text NOT NULL,
    Continent_id integer NOT NULL,
    PRIMARY KEY (Country_id),
    FOREIGN KEY (Continent_id) REFERENCES Continents(Continent_id)
);
