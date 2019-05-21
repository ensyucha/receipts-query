CREATE TABLE IF NOT EXISTS systems (
    password TEXT,
    unusedusage INTEGER,
    apicode TEXT
);

CREATE TABLE IF NOT EXISTS users (
    userid INTEGER PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(100) unique,
    nickname VARCHAR(100) unique,
    password TEXT,
    usages INTEGER
);

CREATE TABLE IF NOT EXISTS logs (
    logid INTEGER PRIMARY KEY AUTO_INCREMENT,
    logtime TEXT,
    logtype TEXT,
    loginfo TEXT,
    username VARCHAR(100)
);


DELETE FROM systems WHERE apicode = 'empty';

INSERT INTO systems(password, unusedusage, apicode) VALUES ('123456', 0, 'empty');