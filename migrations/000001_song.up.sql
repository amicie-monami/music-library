-- Songs table
CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    group_name VARCHAR(32) NOT NULL,
    song_title VARCHAR(32) NOT NULL
);

-- Song details table
CREATE TABLE song_details (
    id BIGSERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    release_date VARCHAR(16),
    text VARCHAR(64),
    link VARCHAR(64)
);

-- Trigger function
CREATE OR REPLACE FUNCTION add_song_detail()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO song_details (song_id) VALUES (NEW.id); 
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger creating
CREATE TRIGGER after_insert_song
AFTER INSERT ON songs
FOR EACH ROW
EXECUTE FUNCTION add_song_detail();
