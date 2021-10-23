-- Copyright 2021 Oscar Triano García

-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at

--    http://www.apache.org/licenses/LICENSE-2.0

-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

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
    date DATE          NOT NULL
                       DEFAULT (CURRENT_DATE)
);

CREATE TABLE IF NOT EXISTS Option (
    title           VARCHAR (255)   DEFAULT "My Site",
    posts_per_page  INTEGER         DEFAULT 3,
    output          VARCHAR (255)   DEFAULT "output",
    license         VARCHAR (1024)  DEFAULT ""
);

CREATE TABLE IF NOT EXISTS Category (
    id   INTEGER       PRIMARY KEY AUTOINCREMENT,
    name VARCHAR (64)  UNIQUE
);

CREATE TABLE IF NOT EXISTS Category_Post (
    category_id INT REFERENCES Category (id) ON DELETE CASCADE
                    NOT NULL,
    post_id     INT REFERENCES Post (id) ON DELETE CASCADE
                    NOT NULL
);
