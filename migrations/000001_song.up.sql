CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    group_name VARCHAR(32) NOT NULL,
    song_name VARCHAR(32) NOT NULL
);

CREATE TABLE song_details (
    id BIGSERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    release_date DATE,
    text text,
    link varchar(256)
);

-- add_song_details automatically adds a row with default values to 
-- the song_details table after inserting a record into the songs table
CREATE OR REPLACE FUNCTION add_song_detail()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO song_details (song_id) VALUES (NEW.id); 
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Creating a trigger to insert into the songs table
CREATE TRIGGER after_insert_song
AFTER INSERT ON songs
FOR EACH ROW
EXECUTE FUNCTION add_song_detail();
