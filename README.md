# SendTheSong 🎵
**A platform for sending anonymous heartfelt messages paired with songs from Spotify**

---

## 📖 Project Overview

SendTheSong is a minimalist web platform that allows users to:
- ✏️ Write and send anonymous messages to anyone
- 🎵 Pair messages with songs from Spotify
- 🔍 Browse and search for messages sent to them
- 📝 View message details with embedded Spotify player
- 📜 Explore all messages shared on the platform

**Key Principle**: Completely anonymous - no personal identity stored.

---

## 🚀 Quick Start Guide

### Prerequisites
- **Go 1.22+** (download from [golang.org](https://golang.org))
- **MySQL/MariaDB** (running on localhost:3306)
- **Laragon** (or any MySQL server)

### Installation Steps

#### 1. **Database Setup**
1. Open **Laragon** → MySQL admin (or use MySQL CLI)
2. Create a new database named `sendthesong`:
   ```sql
   CREATE DATABASE sendthesong;
   ```
3. Import the schema:
   ```bash
   mysql -u root sendthesong < db-setup.sql
   ```
   Or run the SQL file through Laragon's MySQL admin interface.

#### 2. **Database Configuration**
Database connection uses the following defaults (edit `goapp/main.go` if needed):
```go
db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/sendthesong")
```

Update the connection string if your MySQL password is different.

#### 3. **Start the Application**

**Option A: Using the batch file (Windows)**
```bash
run-go.bat
```

**Option B: Using terminal**
```bash
cd goapp
go mod tidy
go run .
```

The app will start on **http://localhost:8080**

---

## 📁 Project Structure

```
sendthesong/
├── goapp/                    # Go application
│   ├── main.go              # Main app with routes
│   ├── db.go                # Database layer
│   ├── media.go             # Media/link parsing
│   ├── go.mod               # Go module definition
│   ├── go.sum               # Go dependencies
│   └── templates/           # HTML templates
│       ├── index.html       # Home page
│       ├── submit.html      # Submit message
│       ├── browse.html      # Search messages
│       ├── detail.html      # View message
│       ├── history.html     # All messages
│       └── support.html     # FAQ & contact
│
├── assets/
│   ├── style.css            # Main stylesheet
│   ├── script.js            # Frontend JavaScript
│   └── spotify-logo.png     # Spotify logo
│
├── uploads/                 # User uploaded audio files
├── db-setup.sql             # Database schema
├── run-go.bat               # Windows launcher
└── README.md                # This file
```

---

## 🎨 Features & Pages

### 🏠 Home Page (`/`)
- Handwritten hero section with tagline
- "How It Works" section with 3 feature cards
- Animated horizontal carousel of recent messages
- Call-to-action buttons to submit and browse

### ✏️ Submit Page (`/submit`)
- Form to send anonymous message
- Recipient name input
- Message textarea with handwriting font
- Real-time Spotify song search
- Selected song display
- One-click publish

### 🔍 Browse Page (`/browse`)
- Search bar to find messages by recipient name
- Grid display of matching messages
- Click cards to view full message

### 📖 Detail Page (`/message/:id`)
- Full message with handwriting font
- Album art display
- Embedded Spotify player (iframe)
- Share button (copy link)
- Navigation back to browse

### 📜 History Page (`/history`)
- All messages ordered by date
- Sort options (latest/oldest)
- Card feed layout

### 💬 Support Page (`/support`)
- FAQ section with 8 common questions
- Contact form for inquiries
- Privacy notice
- Platform information

---

## 🎨 Design & Styling

### Fonts
- **Headings**: [Google Fonts - Kalam](https://fonts.google.com/specimen/Kalam) (handwriting)
- **Body**: Inter (clean, modern sans-serif)

### Colors
- **Primary**: #000 (Black)
- **Background**: #fafafa (Light Gray)
- **Accents**: #1DB954 (Spotify Green for icons)
- **Text**: #555 (Dark Gray)

### Animations
- ✨ Fade-in animations on page load
- 🎯 Hover effects on cards and buttons
- 📱 Smooth horizontal scroll for carousel
- ⬆️ Slide-up animations for form submissions

---

## 🎵 Spotify Integration

### How It Works
1. **Client Credentials Flow**: No user login needed
2. **Song Search**: Real-time search via Spotify API
3. **Track Data**: Stores track ID, name, artist, album art, preview URL
4. **Embedded Player**: Display Spotify iframe on detail page

### API Methods
```php
$spotify = new SpotifyAPI();
$results = $spotify->search('song name', 10);
```

### Response Format
```json
{
  "id": "track_id",
  "name": "Song Title",
  "artist": "Artist Name",
  "album_art": "url_to_image",
  "preview_url": "url_to_30s_preview",
  "uri": "spotify:track:..."
}
```

---

## 💾 Database Schema

### Messages Table
```sql
CREATE TABLE messages (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    to_name VARCHAR(255) NOT NULL,
    message LONGTEXT NOT NULL,
    spotify_track_id VARCHAR(255) NOT NULL,
    spotify_track_name VARCHAR(500),
    spotify_artist VARCHAR(500),
    spotify_album_art LONGTEXT,
    spotify_preview_url LONGTEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_to_name (to_name),
    INDEX idx_created_at (created_at)
);
```

---

## 🔧 API Endpoints

### POST `/api/submit.php`
Submit a new message
```php
POST /api/submit.php
Body:
- to_name: string
- message: string
- spotify_track_id: string
- spotify_data: json (name, artist, album_art, preview_url)

Response: {"success": true/false, "message": "..."}
```

### GET `/api/spotify-search.php?q=query`
Search Spotify tracks
```php
GET /api/spotify-search.php?q=song+name

Response: {"results": [array of tracks]}
```

### GET `/api/search.php?q=name`
Search messages by recipient name
```php
GET /api/search.php?q=john

Response: {"results": [array of messages], "count": number}
```

---

## 🚀 Deployment Guide

### Option 1: Local Laragon Development
1. Place project in `Laragon/www/sendthesong`
2. Follow installation steps above
3. Access at `http://localhost/sendthesong`

### Option 2: Web Hosting (Shared/VPS)
1. **Upload Files**
   ```bash
   # Via FTP or File Manager
   # Upload all files to public_html or desired directory
   ```

2. **Create Database**
   ```bash
   # Via cPanel → Databases
   # Create new MySQL database
   # Import db-setup.sql
   ```

3. **Update Configuration**
   - Edit `includes/config.php`
   - Update DB_HOST (usually "localhost")
   - Update DB_USER and DB_PASS from cPanel
   - Add Spotify credentials
   - Update SITE_URL to your domain

4. **File Permissions**
   ```bash
   chmod 755 ./
   chmod 755 ./logs/
   chmod 644 ./includes/config.php
   ```

5. **SSL Certificate**
   - Use free SSL (Let's Encrypt via cPanel)
   - Update SITE_URL to use `https://`

### Option 3: Vercel/Cloud Deployment

> **Note**: Vercel is serverless and PHP-based apps need traditional hosting. For cloud deployment, consider:
- **Railway** (PHP + MySQL support)
- **Render** (PHP friendly)
- **Heroku** (traditional but being deprecated)

For Vercel, you would need to refactor to Next.js.

### Environment Variables (Production)
Create `.env` file (or use hosting control panel):
```env
DB_HOST=your_db_host
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=sendthesong
SPOTIFY_CLIENT_ID=your_spotify_id
SPOTIFY_CLIENT_SECRET=your_spotify_secret
SITE_URL=https://yourdomain.com
```

---

## 🔒 Security Best Practices

1. ✅ Never commit `.env` or config files with credentials
2. ✅ Use SSL/HTTPS in production
3. ✅ Validate all inputs on server side
4. ✅ Use prepared statements (already implemented with PDO)
5. ✅ Keep dependencies updated
6. ✅ Monitor access logs for suspicious activity
7. ✅ Implement rate limiting for API endpoints
8. ✅ Regular backups of database

---

## 📱 Responsive Design

The website is fully responsive and optimized for:
- 📱 Mobile phones (320px+)
- 📱 Tablets (768px+)
- 💻 Desktop (1200px+)

All carousels, forms, and layouts adapt automatically.

---

## 🔄 Customization Guide

### Change Colors
Edit `assets/style.css`:
```css
/* Primary colors */
--primary: #000;
--secondary: #1DB954;
--background: #fafafa;
```

### Change Fonts
Edit `includes/config.php` and `assets/style.css`:
```css
font-family: 'YourFont', cursive;  /* Handwriting */
font-family: 'YourFont', sans-serif; /* Body */
```

### Add More Fields
1. Update database schema in `db-setup.sql`
2. Modify `includes/db.php` methods
3. Update form in `submit.php`
4. Update display in `detail.php`

---

## 🐛 Troubleshooting

### Database Connection Error
- ✅ Verify MySQL is running in Laragon
- ✅ Check `includes/config.php` credentials
- ✅ Ensure `sendthesong` database exists

### Spotify Search Not Working
- ✅ Verify API credentials in `includes/config.php`
- ✅ Check Spotify API quotas in dashboard
- ✅ Verify cURL is enabled in PHP

### Messages Not Appearing
- ✅ Check database has messages table
- ✅ Verify messages were inserted correctly
- ✅ Check browser console for JavaScript errors

### Blank Page or 500 Error
- ✅ Check PHP error log
- ✅ Enable debug mode (set `display_errors` to 1)
- ✅ Verify file permissions are correct

---

## 📊 Performance Tips

1. **Database Optimization**
   - Indexes on `to_name` and `created_at` are included
   - Limit queries with LIMIT clause

2. **Caching**
   - Add browser caching for static assets
   - Cache Spotify API results for 1 hour

3. **Image Optimization**
   - Album art is served from Spotify CDN
   - Compress CSS and JavaScript in production

4. **Database Maintenance**
   - Regular backups
   - Archive old messages if database grows large

---

## 📝 Future Enhancements

- 🔐 User accounts and message management
- ❤️ Like/reaction system
- 🏆 Trending messages
- 🔔 Notifications for new messages
- 📸 Custom image upload
- 🎤 Voice message recording
- 🌍 Multi-language support
- 💬 Comments on messages
- 📊 Analytics dashboard

---

## 📄 License

This project is open source. Feel free to use, modify, and deploy!


## 🤝 Support

For issues or questions:
1. Check the FAQ on `/support` page
2. Review this README
3. Check PHP error logs
4. Test locally before deploying


## 🎉 Enjoy SendTheSong!

Express untold emotions through music and words.

**Made with ❤️**
# Go Migration

The Go version lives in `goapp/` and can be started on Windows with `run-go.bat`.

Quick start:

```bat
run-go.bat
```

Or manually:

```bash
cd goapp
go mod tidy
go run .
```

Use the Go version for the rewrite. The PHP files are kept as the old implementation for reference while we transition.
