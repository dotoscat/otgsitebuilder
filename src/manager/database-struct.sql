CREATE TABLE IF NOT EXISTS Content (
    id             INTEGER       PRIMARY KEY AUTOINCREMENT,
    file           VARCHAR (255) UNIQUE
                                 NOT NULL,
    date           DATETIME      NOT NULL
                                 DEFAULT (CURRENT_TIMESTAMP),
    contenttype_id INTEGER       REFERENCES ContentType (id) ON DELETE NO ACTION
                                                             ON UPDATE NO ACTION
                                 DEFAULT (1)
);

CREATE TABLE IF NOT EXISTS ContentType(
    id   INTEGER     PRIMARY KEY AUTOINCREMENT,
    name VARCHAR (8) UNIQUE
                     NOT NULL
);

INSERT OR IGNORE INTO ContentType (id, name) VALUES (1, 'POST');
INSERT OR IGNORE INTO ContentType (id, name) VALUES (2, 'PAGE');
