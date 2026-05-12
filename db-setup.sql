-- SendTheSong Database Setup
-- Create database
CREATE DATABASE IF NOT EXISTS sendthesong;
USE sendthesong;

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    to_name VARCHAR(255) NOT NULL,
    message LONGTEXT NOT NULL,
    spotify_track_id VARCHAR(255) NOT NULL,
    spotify_track_name VARCHAR(500),
    spotify_artist VARCHAR(500),
    spotify_album_art LONGTEXT,
    spotify_preview_url LONGTEXT,
    uploaded_file_path LONGTEXT,
    uploaded_file_name VARCHAR(500),
    uploaded_file_mime VARCHAR(255),
    song_link LONGTEXT,
    song_source VARCHAR(50),
    song_title VARCHAR(500),
    song_thumbnail LONGTEXT,
    song_embed_url LONGTEXT,
    song_audio_url LONGTEXT,
    song_provider_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_to_name (to_name),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
