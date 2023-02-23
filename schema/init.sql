CREATE DATABASE video_uploader;
\c video_uploader;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; 

CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  email varchar (255),
  password varchar (255)
);

CREATE TABLE videos (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  author_id uuid REFERENCES users(id),
  is_stream_ready BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);