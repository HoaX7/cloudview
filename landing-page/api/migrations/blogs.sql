CREATE TABLE IF NOT EXISTS blogs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content VARCHAR(255) NOT NULL,
    sub_title TEXT,
    slug TEXT NOT NULL,
    metadata JSONB,
	description TEXT,
	keywords TEXT,
	featured_image_url TEXT,
	created_at integer(4) not null default (strftime('%s','now')),
	updated_at integer(4) NULL,
	deleted_at integer(4) NULL
);

CREATE TRIGGER IF NOT EXISTS update_blogs_updated_at
	AFTER UPDATE ON blogs
	FOR EACH ROW
	BEGIN
    UPDATE blogs SET updated_at = strftime('%s','now') WHERE id = OLD.id;
	END;

--needed because of https://github.com/cloudflare/workers-sdk/pull/2912
select count(*) from blogs;