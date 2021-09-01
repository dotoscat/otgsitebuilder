
CREATE TABLE IF NOT EXISTS Page (
    id        INTEGER       PRIMARY KEY AUTOINCREMENT,
    name      VARCHAR (255) UNIQUE
                            NOT NULL,
    reference VARCHAR (255) DEFAULT ""
);

CREATE TABLE IF NOT EXISTS Post (
    id   INTEGER       PRIMARY KEY AUTOINCREMENT,
    name VARCHAR (255) UNIQUE
                       NOT NULL,
    date DATETIME      NOT NULL
                       DEFAULT (CURRENT_TIMESTAMP)
);
