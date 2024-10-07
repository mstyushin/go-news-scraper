INSERT INTO rss_feeds (url) VALUES ('http://rss.slashdot.org/Slashdot/slashdot');
INSERT INTO rss_feeds (url) VALUES ('https://news.ycombinator.com/rss');

INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (1, 'test article', 'test content', 'http://some.site.com/2024/whatever', 1724877777);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'another article', 'another content', 'http://another.site.com/2024/something', 1724866666);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post1', 'wow check it out!', 'https://hackernews.com/post/1', 1728118904);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post2', 'wow check it out!', 'https://hackernews.com/post/2', 1728128904);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post3', 'wow check it out!', 'https://hackernews.com/post/3', 1728138904);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post4', 'wow check it out!', 'https://hackernews.com/post/4', 1728148904);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post5', 'wow check it out!', 'https://hackernews.com/post/5', 1728158904);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post6', 'wow check it out!', 'https://hackernews.com/post/6', 1728168904);
INSERT INTO articles (rss_feed_id, title, content, link, pub_time) VALUES (2, 'Hackernews post7', 'wow check it out!', 'https://hackernews.com/post/7', 1728168904);
