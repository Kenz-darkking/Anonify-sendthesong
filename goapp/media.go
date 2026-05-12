package main

import (
	"net/url"
	"path"
	"regexp"
	"strings"
)

type SongMedia struct {
	SongLink      string
	SongSource    string
	SongTitle     string
	SongThumbnail string
	SongEmbedURL  string
	SongAudioURL  string
	SongProviderID string
}

func ExtractSongMetadata(input string) SongMedia {
	link := strings.TrimSpace(input)
	media := SongMedia{SongLink: link}
	media.SongSource = detectSource(link)

	switch media.SongSource {
	case "youtube":
		id := extractYouTubeID(link)
		media.SongProviderID = id
		if id != "" {
			media.SongEmbedURL = "https://www.youtube.com/embed/" + id
			media.SongThumbnail = "https://img.youtube.com/vi/" + id + "/hqdefault.jpg"
		}
		media.SongTitle = "YouTube Song"
	case "spotify":
		id := extractSpotifyID(link)
		media.SongProviderID = id
		if id != "" {
			media.SongEmbedURL = "https://open.spotify.com/embed/track/" + id
		}
		media.SongTitle = "Spotify Track"
	case "audio":
		media.SongAudioURL = link
		media.SongTitle = baseName(link)
	default:
		media.SongTitle = link
	}

	return media
}

func detectSource(link string) string {
	if link == "" {
		return "unknown"
	}
	if strings.Contains(link, "open.spotify.com/track/") || strings.HasPrefix(link, "spotify:track:") {
		return "spotify"
	}
	if strings.Contains(link, "youtube.com/watch?v=") || strings.Contains(link, "youtu.be/") || strings.Contains(link, "youtube.com/shorts/") {
		return "youtube"
	}
	if regexp.MustCompile(`(?i)\.(mp3|m4a|wav|ogg|aac|webm)(\?.*)?$`).MatchString(link) {
		return "audio"
	}
	return "link"
}

func extractYouTubeID(link string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`youtu\.be/([A-Za-z0-9_-]{6,})`),
		regexp.MustCompile(`v=([A-Za-z0-9_-]{6,})`),
		regexp.MustCompile(`youtube\.com/shorts/([A-Za-z0-9_-]{6,})`),
	}
	for _, re := range patterns {
		if m := re.FindStringSubmatch(link); len(m) == 2 {
			return m[1]
		}
	}
	return ""
}

func extractSpotifyID(link string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`open\.spotify\.com/track/([A-Za-z0-9]+)`),
		regexp.MustCompile(`spotify:track:([A-Za-z0-9]+)`),
	}
	for _, re := range patterns {
		if m := re.FindStringSubmatch(link); len(m) == 2 {
			return m[1]
		}
	}
	return ""
}

func baseName(raw string) string {
	if u, err := url.Parse(raw); err == nil && u.Path != "" {
		if name := path.Base(u.Path); name != "." && name != "/" {
			return name
		}
	}
	return raw
}
