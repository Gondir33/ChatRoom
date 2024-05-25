CREATE TABLE room (
    Id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE message (
    Id SERIAL PRIMARY KEY,
    mess VARCHAR(MAX),
    room_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES room (Id)
);