CREATE TABLE IF NOT EXISTS characteristics (
    ID SERIAL PRIMARY KEY,
    Age INTEGER NOT NULL,
    Gender VARCHAR(255) NOT NULL,
    Nationality VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS people (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Surname VARCHAR(255) NOT NULL,
    Patronymic VARCHAR(255),
    Characteristic_ID INTEGER,
    FOREIGN KEY (Characteristic_ID) REFERENCES characteristics(ID)
);
