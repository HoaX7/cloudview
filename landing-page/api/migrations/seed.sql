CREATE TABLE IF NOT EXISTS registrations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name varchar(255),
	email_address TEXT NOT NULL UNIQUE,
	notes text,
	is_demo_completed boolean default false,
	is_onboarded boolean default false,
	call_scheduled_at timestamptz,
	created_at integer(4) not null default (strftime('%s','now')),
	updated_at integer(4) NULL,
	deleted_at integer(4) NULL
);

CREATE TRIGGER IF NOT EXISTS update_registrations_updated_at
	AFTER UPDATE ON registrations
	FOR EACH ROW
	BEGIN
    UPDATE registrations SET updated_at = strftime('%s','now') WHERE id = OLD.id;
	END;

--needed because of https://github.com/cloudflare/workers-sdk/pull/2912
select count(*) from registrations;