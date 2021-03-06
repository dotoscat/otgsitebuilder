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

CREATE TABLE IF NOT EXISTS Document (
    id        INTEGER       PRIMARY KEY AUTOINCREMENT,
    title      VARCHAR (255) UNIQUE
                            NOT NULL,
    date DATE          NOT NULL
                       DEFAULT (CURRENT_DATE),
    content TEXT       DEFAULT ""
                        NOT NULL,
);

CREATE TABLE IF NOT EXISTS Option (
    title           VARCHAR (255)   NOT NULL
                                    DEFAULT "My Site",
    posts_per_page  INTEGER         DEFAULT 3,
    output          VARCHAR (255)   NOT NULL
                                    DEFAULT "output",
    license         VARCHAR (1024)  NOT NULL
                                    DEFAULT ""
);

CREATE TABLE IF NOT EXISTS Category (
    id   INTEGER       PRIMARY KEY AUTOINCREMENT,
    name VARCHAR (128)  UNIQUE
);

CREATE TABLE IF NOT EXISTS Tag (
    id   INTEGER       PRIMARY KEY AUTOINCREMENT,
    name VARCHAR (128)  UNIQUE
);

-- Reserved tags: Post, Page and Draf

CREATE TABLE IF NOT EXISTS Category_Document (
    category_id INT REFERENCES Category (id) ON DELETE CASCADE
                    NOT NULL,
    document_id     INT REFERENCES Document (id) ON DELETE CASCADE
                    NOT NULL
);

CREATE TABLE IF NOT EXISTS Tag_Document (
    category_id INT REFERENCES Category (id) ON DELETE CASCADE
                    NOT NULL,
    document_id     INT REFERENCES Document (id) ON DELETE CASCADE
                    NOT NULL
);
