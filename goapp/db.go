package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

type Message struct {
	ID             string
	ToName         string
	Message        string
	SpotifyTrackID  sql.NullString
	SpotifyTrackName sql.NullString
	SpotifyArtist   sql.NullString
	SpotifyAlbumArt sql.NullString
	SpotifyPreviewURL sql.NullString
	UploadedFilePath sql.NullString
	UploadedFileName sql.NullString
	UploadedFileMime sql.NullString
	SongLink       sql.NullString
	SongSource     sql.NullString
	SongTitle      sql.NullString
	SongThumbnail  sql.NullString
	SongEmbedURL   sql.NullString
	SongAudioURL   sql.NullString
	SongProviderID sql.NullString
	CreatedAt      time.Time
}

type MessageInput struct {
	ToName        string
	Message       string
	SongLink      string
	SongSource    string
	SongTitle     string
	SongThumbnail string
	SongEmbedURL  string
	SongAudioURL  string
	SongProviderID string
}

func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) LatestMessages(limit int) ([]Message, error) {
	return s.MessagesBySort("latest", limit, 0)
}

func (s *Store) MessagesByPage(limit, offset int) ([]Message, error) {
	return s.queryMessages("DESC", "", limit, offset)
}

func (s *Store) MessagesBySort(sortBy string, limit, offset int) ([]Message, error) {
	order := "DESC"
	if strings.ToLower(sortBy) == "oldest" {
		order = "ASC"
	}
	return s.queryMessages(order, "", limit, offset)
}

func (s *Store) SearchMessages(q string, limit, offset int) ([]Message, error) {
	return s.queryMessages("DESC", q, limit, offset)
}

func (s *Store) CountSearch(q string) (int, error) {
	row := s.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE to_name LIKE ?`, "%"+q+"%")
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) CountMessages() (int, error) {
	row := s.db.QueryRow(`SELECT COUNT(*) FROM messages`)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (s *Store) GetMessageByID(id string) (*Message, error) {
	query := `SELECT id, to_name, message, spotify_track_id, spotify_track_name, spotify_artist, spotify_album_art, spotify_preview_url, uploaded_file_path, uploaded_file_name, uploaded_file_mime, song_link, song_source, song_title, song_thumbnail, song_embed_url, song_audio_url, song_provider_id, created_at FROM messages WHERE id = ? LIMIT 1`
	row := s.db.QueryRow(query, id)
	msg, err := scanMessage(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *Store) CreateMessage(input MessageInput) (string, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	_, err := s.db.Exec(`INSERT INTO messages (
		id, to_name, message, spotify_track_id, spotify_track_name, spotify_artist, spotify_album_art, spotify_preview_url,
		uploaded_file_path, uploaded_file_name, uploaded_file_mime,
		song_link, song_source, song_title, song_thumbnail, song_embed_url, song_audio_url, song_provider_id,
		created_at
	) VALUES (?, ?, ?, '', '', '', '', '', NULL, NULL, NULL, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		id, input.ToName, input.Message,
		input.SongLink, input.SongSource, input.SongTitle, input.SongThumbnail, input.SongEmbedURL, input.SongAudioURL, input.SongProviderID,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Store) queryMessages(orderBy string, search string, limit, offset int) ([]Message, error) {
	base := `SELECT id, to_name, message, spotify_track_id, spotify_track_name, spotify_artist, spotify_album_art, spotify_preview_url, uploaded_file_path, uploaded_file_name, uploaded_file_mime, song_link, song_source, song_title, song_thumbnail, song_embed_url, song_audio_url, song_provider_id, created_at FROM messages`
	args := []any{}
	if search != "" {
		base += ` WHERE to_name LIKE ?`
		args = append(args, "%"+search+"%")
	}
	base += ` ORDER BY created_at ` + orderBy + ` LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	rows, err := s.db.Query(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Message
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, msg)
	}
	return out, rows.Err()
}

func scanMessage(scanner interface{ Scan(...any) error }) (Message, error) {
	var msg Message
	err := scanner.Scan(
		&msg.ID,
		&msg.ToName,
		&msg.Message,
		&msg.SpotifyTrackID,
		&msg.SpotifyTrackName,
		&msg.SpotifyArtist,
		&msg.SpotifyAlbumArt,
		&msg.SpotifyPreviewURL,
		&msg.UploadedFilePath,
		&msg.UploadedFileName,
		&msg.UploadedFileMime,
		&msg.SongLink,
		&msg.SongSource,
		&msg.SongTitle,
		&msg.SongThumbnail,
		&msg.SongEmbedURL,
		&msg.SongAudioURL,
		&msg.SongProviderID,
		&msg.CreatedAt,
	)
	return msg, err
}
