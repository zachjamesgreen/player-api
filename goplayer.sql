DROP TABLE artists CASCADE;
DROP TABLE albums CASCADE;
DROP TABLE songs;

CREATE TABLE IF NOT EXISTS artists (
  id serial PRIMARY KEY,
  name varchar NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS albums (
  id serial PRIMARY KEY,
  title varchar NOT NULL UNIQUE,
  year int,
  artistId integer REFERENCES artists (id)
);

CREATE TABLE IF NOT EXISTS songs (
  id serial PRIMARY KEY,
  title varchar NOT NULL,
  file_type varchar,
  file_path varchar,
  file_name varchar,
  artistId integer REFERENCES artists (id),
  albumId integer REFERENCES albums (id)
);


INSERT INTO artists (name)
VALUES ('Flume')
RETURNING id, name;

INSERT INTO albums (title)
VALUES ('AIM')
RETURNING id, title;



cd ../fileuploader; go build && go install; cd ../models; go build && go install; cd ../restapi; go build; ./restapi;

SELECT songs.id, songs.title, songs.file_type, songs.file_name, albums.title as album_title, albums.year, artists.name
from songs
inner join albums on songs.albumId = albums.id
inner join artists on songs.artistId = artists.id;
