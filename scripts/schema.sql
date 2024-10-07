DROP TABLE IF EXISTS articles, rss_feeds;

CREATE TABLE rss_feeds (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    UNIQUE(url)
);

CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    rss_feed_id INTEGER REFERENCES rss_feeds(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    link TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    UNIQUE(link)
);

CREATE INDEX idx_articles_pub_time ON articles (pub_time);
