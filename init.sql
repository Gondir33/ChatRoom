CREATE TABLE room (
    Id SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE message (
    Id SERIAL PRIMARY KEY,
    mess TEXT,
    room_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES room (Id)
);