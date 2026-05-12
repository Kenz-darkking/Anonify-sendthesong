package main

import (
	"embed"
	"html/template"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*.html
var templateFS embed.FS

type AppConfig struct {
	DSN      string
	SiteName string
	Port     string
}

type PageData struct {
	SiteName     string
	Title        string
	Active       string
	Error        string
	Message      string
	SearchQuery  string
	SortBy       string
	Page         int
	TotalPages   int
	TotalResults int
	PageLinks    []int
	Messages     []Message
	MessageItem  *Message
}

func main() {
	cfg := loadConfig()
	db, err := NewStore(cfg.DSN)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	
	// Create FuncMap
	funcMap := template.FuncMap{
		"trimMessage": trimMessage,
		"formatDate":  formatDate,
		"safeURL":     safeURL,
		"add":         func(a, b int) int { return a + b },
		"sub":         func(a, b int) int { return a - b },
	}
	
	// Parse templates with FuncMap
	tmpl := template.New("").Funcs(funcMap)
	loadedTemplates, err := tmpl.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatalf("template parse failed: %v", err)
	}
	r.SetHTMLTemplate(loadedTemplates)

	r.Static("/assets", "../assets")
	r.Static("/uploads", "../uploads")

	r.GET("/", func(c *gin.Context) {
		messages, err := db.LatestMessages(8)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to load homepage")
			return
		}
		c.HTML(http.StatusOK, "index.html", PageData{
			SiteName: cfg.SiteName,
			Title:    cfg.SiteName,
			Active:   "home",
			Messages: messages,
		})
	})

	r.GET("/submit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "submit.html", PageData{SiteName: cfg.SiteName, Title: "Tell Your Story", Active: "submit"})
	})

	r.POST("/submit", func(c *gin.Context) {
		toName := strings.TrimSpace(c.PostForm("to_name"))
		message := strings.TrimSpace(c.PostForm("message"))
		songLink := strings.TrimSpace(c.PostForm("song_link"))

		if toName == "" || message == "" || songLink == "" {
			c.HTML(http.StatusBadRequest, "submit.html", PageData{
				SiteName: cfg.SiteName,
				Title:    "Tell Your Story",
				Active:   "submit",
				Error:    "Please fill all required fields.",
			})
			return
		}

		media := ExtractSongMetadata(songLink)
		_, err := db.CreateMessage(MessageInput{
			ToName:        toName,
			Message:       message,
			SongLink:      media.SongLink,
			SongSource:    media.SongSource,
			SongTitle:     media.SongTitle,
			SongThumbnail: media.SongThumbnail,
			SongEmbedURL:  media.SongEmbedURL,
			SongAudioURL:  media.SongAudioURL,
			SongProviderID: media.SongProviderID,
		})
		if err != nil {
			c.HTML(http.StatusInternalServerError, "submit.html", PageData{
				SiteName: cfg.SiteName,
				Title:    "Tell Your Story",
				Active:   "submit",
				Error:    err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "submit.html", PageData{
			SiteName: cfg.SiteName,
			Title:    "Tell Your Story",
			Active:   "submit",
			Message:  "Message published successfully!",
		})
	})

	r.GET("/browse", func(c *gin.Context) {
		searchQuery := strings.TrimSpace(c.Query("search"))
		page := atoiDefault(c.Query("page"), 1)
		if page < 1 {
			page = 1
		}
		perPage := 10
		offset := (page - 1) * perPage

		messages := []Message{}
		total := 0
		var err error
		if searchQuery != "" {
			total, err = db.CountSearch(searchQuery)
			if err != nil {
				c.String(http.StatusInternalServerError, "failed to search messages")
				return
			}
			messages, err = db.SearchMessages(searchQuery, perPage, offset)
			if err != nil {
				c.String(http.StatusInternalServerError, "failed to load messages")
				return
			}
		}

		totalPages := 1
		if total > 0 {
			totalPages = (total + perPage - 1) / perPage
		}

		c.HTML(http.StatusOK, "browse.html", PageData{
			SiteName:     cfg.SiteName,
			Title:        "Browse Messages",
			Active:       "browse",
			SearchQuery:  searchQuery,
			Page:         page,
			TotalPages:   totalPages,
			TotalResults: total,
			PageLinks:    pageLinks(page, totalPages),
			Messages:     messages,
		})
	})

	r.GET("/detail/:id", func(c *gin.Context) {
		msg, err := db.GetMessageByID(c.Param("id"))
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to load message")
			return
		}
		if msg == nil {
			c.HTML(http.StatusNotFound, "detail.html", PageData{SiteName: cfg.SiteName, Title: "Message" , Active: "browse", Error: "Message not found"})
			return
		}
		c.HTML(http.StatusOK, "detail.html", PageData{SiteName: cfg.SiteName, Title: "Message for " + msg.ToName, Active: "browse", MessageItem: msg})
	})

	r.GET("/history", func(c *gin.Context) {
		sortBy := strings.ToLower(strings.TrimSpace(c.DefaultQuery("sort", "latest")))
		page := atoiDefault(c.Query("page"), 1)
		if page < 1 {
			page = 1
		}
		perPage := 10
		offset := (page - 1) * perPage

		total, err := db.CountMessages()
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to load messages")
			return
		}
		messages, err := db.MessagesBySort(sortBy, perPage, offset)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to load messages")
			return
		}
		totalPages := 1
		if total > 0 {
			totalPages = (total + perPage - 1) / perPage
		}

		c.HTML(http.StatusOK, "history.html", PageData{
			SiteName:     cfg.SiteName,
			Title:        "History",
			Active:       "history",
			SortBy:       sortBy,
			Page:         page,
			TotalPages:   totalPages,
			TotalResults: total,
			PageLinks:    pageLinks(page, totalPages),
			Messages:     messages,
		})
	})

	r.GET("/support", func(c *gin.Context) {
		success := c.Query("success") == "1"
		errorFlag := c.Query("error") == "1"
		page := PageData{SiteName: cfg.SiteName, Title: "Support", Active: "support"}
		if success {
			page.Message = "✅ Thank you for your message! We'll get back to you soon."
		}
		if errorFlag {
			page.Error = "❌ Please fill in all required fields correctly."
		}
		c.HTML(http.StatusOK, "support.html", page)
	})

	r.POST("/support", func(c *gin.Context) {
		email := strings.TrimSpace(c.PostForm("email"))
		message := strings.TrimSpace(c.PostForm("message"))
		if email == "" || message == "" {
			c.Redirect(http.StatusSeeOther, "/support?error=1")
			return
		}
		c.Redirect(http.StatusSeeOther, "/support?success=1")
	})

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() AppConfig {
	dsn := os.Getenv("SENDTHESONG_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/sendthesong?parseTime=true&loc=Local"
	}
	return AppConfig{
		DSN:      dsn,
		SiteName: envOr("SENDTHESONG_SITE_NAME", "SendTheSong"),
		Port:     envOr("PORT", "8080"),
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func atoiDefault(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	var n int
	_, err := fmt.Sscanf(value, "%d", &n)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func trimMessage(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 150 {
		return value
	}
	return value[:150] + "..."
}

func formatDate(t time.Time) string {
	return t.Format("Jan 02, 2006")
}

func safeURL(value string) string {
	if value == "" {
		return ""
	}
	return filepath.Clean(value)
}

func pageLinks(current, total int) []int {
	if total < 1 {
		return []int{1}
	}
	start := current - 2
	if start < 1 {
		start = 1
	}
	end := current + 2
	if end > total {
		end = total
	}
	links := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		links = append(links, i)
	}
	return links
}
