CREATE TABLE timestamps
(
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    first INTEGER NOT NULL CHECK (first >= 0),
    rest  INTEGER NOT NULL CHECK (rest >= 0),
    UNIQUE (first, rest)
);

CREATE TABLE platforms
(
    id    TEXT PRIMARY KEY CHECK (length(id) BETWEEN 1 AND 64),
    title TEXT NOT NULL,
    url   TEXT NOT NULL
);

CREATE TABLE courses
(
    id           TEXT PRIMARY KEY CHECK (length(id) BETWEEN 1 AND 64),
    slug         TEXT UNIQUE NOT NULL,
    timestamp_id INTEGER     NOT NULL,
    platform_id  TEXT        NOT NULL,
    FOREIGN KEY (timestamp_id) REFERENCES timestamps (id),
    FOREIGN KEY (platform_id) REFERENCES platforms (id)
);

CREATE TABLE course_urls
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    course_id TEXT                                       NOT NULL,
    url       TEXT                                       NOT NULL,
    category  TEXT CHECK (category IN ('url', 'gdrive', 'magnet')) NOT NULL,
    position  INTEGER                                    NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE,
    UNIQUE (course_id, position)
);