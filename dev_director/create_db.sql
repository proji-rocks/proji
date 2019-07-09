
-- Create tables
 CREATE TABLE
	IF NOT EXISTS languages(id INTEGER PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	comment VARCHAR(255));

CREATE TABLE
	IF NOT EXISTS languages_short(id INTEGER PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	language_id INTEGER NOT NULL,
	comment VARCHAR(255),
	FOREIGN KEY(language_id) REFERENCES languages(id));

CREATE TABLE
	IF NOT EXISTS folders(id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(50) NOT NULL,
	language_id INTEGER NOT NULL,
	relative_path VARCHAR(255) NOT NULL,
	comment VARCHAR(255),
	FOREIGN KEY(language_id) REFERENCES languages(id));

CREATE TABLE
	IF NOT EXISTS commands_general(id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(50) NOT NULL,
	command VARCHAR(255) NOT NULL,
	comment VARCHAR(255));

CREATE TABLE
	IF NOT EXISTS commands_specific(id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(50) NOT NULL,
	language_id INTEGER NOT NULL,
	command VARCHAR(255) NOT NULL,
	comment VARCHAR(255),
	FOREIGN KEY(language_id) REFERENCES languages(id));
-- Inserts
 INSERT
	INTO
		languages(name)
	VALUES ('C-Plus-Plus'),
	('C'),
	('Python'),
	('HTML'),
	('PHP'),
	('JavaScript');

INSERT
	INTO
		languages_short(name,
		language_id)
	VALUES('cpp',
	1),
	('c++',
	1),
	('cc',
	1),
	('c',
	2),
	('python',
	3),
	('py',
	3),
	('html',
	4),
	('php',
	5),
	('javascript',
	6),
	('js',
	6);

INSERT
	INTO
		folders(name,
		language_id,
		relative_path,
		comment)
	VALUES ('build',
	);
