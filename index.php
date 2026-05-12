<?php
require_once 'includes/db.php';
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SendTheSong - A bunch of untold words, sent through the song</title>
    <link href="https://fonts.googleapis.com/css2?family=Kalam:wght@300;400;700&family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="assets/style.css">
</head>
<body>
    <!-- Navigation -->
    <nav class="navbar">
        <div class="nav-container">
            <a href="index.php" class="logo">sendthesong</a>
            <div class="nav-menu">
                <a href="index.php" class="nav-link active">Home</a>
                <a href="submit.php" class="nav-link">Tell Your Story</a>
                <a href="browse.php" class="nav-link">Browse</a>
                <a href="history.php" class="nav-link">History</a>
                <a href="support.php" class="nav-link">Support</a>
            </div>
        </div>
    </nav>

    <!-- Hero Section -->
    <section class="hero">
        <div class="hero-content">
            <h1 class="hero-title">a bunch of the untold words, sent through the song</h1>
            <p class="hero-subtitle">Express your untold message through the song.</p>
            <div class="cta-buttons">
                <a href="submit.php" class="btn btn-primary">
                    <span class="btn-icon">✏️</span> Tell Your Story
                </a>
                <a href="browse.php" class="btn btn-secondary">
                    <span class="btn-icon">🔍</span> Browse the Stories
                </a>
            </div>
        </div>
    </section>

    <!-- How It Works Section -->
    <section class="how-it-works">
        <h2>How It Works</h2>
        <div class="features-grid">
            <div class="feature-card">
                <div class="feature-number">1</div>
                <h3>Share your Messages</h3>
                <p>Choose a song and write a heartfelt message to someone special or save it as a little gift for yourself.</p>
            </div>
            <div class="feature-card">
                <div class="feature-number">2</div>
                <h3>Browse Messages</h3>
                <p>Find messages that were written for you. Search your name and uncover heartfelt messages written just for you.</p>
            </div>
            <div class="feature-card">
                <div class="feature-number">3</div>
                <h3>Detail Messages</h3>
                <p>Tap on any message card to discover the full story behind it and listen to the song that captures the emotion.</p>
            </div>
        </div>
    </section>

    <!-- Stories Feed -->
    <section class="stories-feed">
        <h2>Messages Feed</h2>
        
        <div class="feed-container">
            <div class="carousel-scroll">
                <?php
                $messages = $db->getLatestMessages(8);
                if(empty($messages)) {
                    echo '<div class="no-messages"><p>No messages yet. Be the first to share!</p></div>';
                } else {
                    foreach($messages as $msg) {
                        echo '<div class="message-card">';
                        echo '<div class="card-header">';
                        echo '<span class="card-to">To: ' . htmlspecialchars($msg['to_name']) . '</span>';
                        echo '</div>';
                        
                        echo '<div class="card-content">';
                        echo '<p class="card-message">' . htmlspecialchars(substr($msg['message'], 0, 150)) . (strlen($msg['message']) > 150 ? '...' : '') . '</p>';
                        echo '</div>';
                        
                        echo '<div class="card-footer">';
                        if($msg['spotify_album_art']) {
                            echo '<img loading="lazy" src="' . htmlspecialchars($msg['spotify_album_art']) . '" alt="Album art" class="album-thumb">';
                        }
                        echo '<div class="song-info">';
                        echo '<div class="song-title">' . htmlspecialchars($msg['spotify_track_name']) . '</div>';
                        echo '<div class="song-artist">' . htmlspecialchars($msg['spotify_artist']) . '</div>';
                        echo '</div>';
                        echo '<img src="assets/spotify-logo.png" alt="Spotify" class="spotify-icon">';
                        echo '</div>';
                        
                        echo '<a href="detail.php?id=' . urlencode($msg['id']) . '" class="card-link"></a>';
                        echo '</div>';
                    }
                }
                ?>
            </div>
        </div>

        <div class="feed-cta">
            <a href="browse.php" class="btn btn-secondary">See All Messages</a>
        </div>
    </section>

    <footer class="footer">
        <p>&copy; 2024 SendTheSong. This is an anonymous platform. We do not store personal identity.</p>
    </footer>

    <script src="assets/script.js"></script>
</body>
</html>
