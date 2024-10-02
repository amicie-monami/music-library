DELETE FROM songs;
DELETE FROM song_details;

INSERT INTO songs (group_name, song_name) VALUES 
('The Beatles', 'Hey Jude'),
('Queen', 'Bohemian Rhapsody'),
('Pink Floyd', 'Comfortably Numb'),
('Led Zeppelin', 'Stairway to Heaven'),
('The Rolling Stones', 'Paint It Black'),
('The Who', 'Baba O''Riley'),
('Nirvana', 'Smells Like Teen Spirit'),
('Metallica', 'Enter Sandman'),
('AC/DC', 'Thunderstruck'),
('Guns N'' Roses', 'Sweet Child O'' Mine');

WITH random_services AS (
    SELECT unnest(array['spotify.com', 'apple.com', 'soundcloud.com']) as service
)

UPDATE song_details SET
    release_date = (
        CASE song_id
        WHEN 1  THEN TO_DATE('26-08-1968', 'dd.mm.yyyy')
        WHEN 2  THEN TO_DATE('31-10-1975', 'dd.mm.yyyy')
        WHEN 3  THEN TO_DATE('23-11-1979', 'dd.mm.yyyy')
        WHEN 4  THEN TO_DATE('08-11-1971', 'dd.mm.yyyy')
        WHEN 5  THEN TO_DATE('07-05-1966', 'dd.mm.yyyy')
        WHEN 6  THEN TO_DATE('14-04-1971', 'dd.mm.yyyy')
        WHEN 7  THEN TO_DATE('10-09-1991', 'dd.mm.yyyy')
        WHEN 8  THEN TO_DATE('30-07-1991', 'dd.mm.yyyy')
        WHEN 9  THEN TO_DATE('24-09-1990', 'dd.mm.yyyy')
        WHEN 10 THEN TO_DATE('21-07-1987', 'dd.mm.yyyy')
        ELSE null END
    ),
    
    link = (
        SELECT service || '/track/' || song_id 
        FROM random_services 
        ORDER BY random() 
        LIMIT 1
    ),

    text = (
        CASE song_id
        WHEN 1 THEN 'Hey Jude, don''t make it bad\n\nTake a sad song and make it better\n\nRemember to let her into your heart\n\nThen you can start to make it better'
        WHEN 2 THEN 'Is this the real life?\n\nIs this just fantasy?\n\nCaught in a landslide,\n\nNo escape from reality'
        WHEN 3 THEN 'Hello?\n\nIs there anybody in there?\n\nJust nod if you can hear me\n\nIs there anyone home?'
        WHEN 4 THEN 'There''s a lady who''s sure\n\nAll that glitters is gold\n\nAnd she''s buying a stairway to heaven'
        WHEN 5 THEN 'I see a red door and I want it painted black\n\nNo colors anymore I want them to turn black'
        WHEN 6 THEN 'Out here in the fields\n\nI fight for my meals\n\nI get my back into my living'
        WHEN 7 THEN 'With the lights out, it''s less dangerous\n\nHere we are now, entertain us\n\nI feel stupid and contagious'
        WHEN 8 THEN 'Say your prayers, little one\n\nDon''t forget, my son\n\nTo include everyone'
        WHEN 9 THEN 'Thunder, thunder, thunder, thunder\n\nI was caught in the middle of a railroad track'
        WHEN 10 THEN 'She''s got a smile that it seems to me\n\nReminds me of childhood memories\n\nWhere everything was as fresh as the bright blue sky'
        ELSE 'No lyrics available' END
    );
