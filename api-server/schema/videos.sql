USE video_uploader;

CREATE TABLE videos (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  author_id uuid REFERENCES users(id),
  is_stream_ready BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);