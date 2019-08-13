--
-- DROP TABLES
--
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS languages_short;
DROP TABLE IF EXISTS default_folders;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS default_files;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS commands_general;
DROP TABLE IF EXISTS commands_specific;
--
-- CREATE TABLES
--
CREATE TABLE IF NOT EXISTS languages(
  id INTEGER PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  comment VARCHAR(255)
);
--
CREATE TABLE IF NOT EXISTS languages_short(
  id INTEGER PRIMARY KEY,
  language_id INTEGER NOT NULL,
  name_short VARCHAR(50) NOT NULL,
  comment VARCHAR(255),
  FOREIGN KEY(language_id) REFERENCES languages(id),
  UNIQUE(name_short)
);
--
CREATE TABLE IF NOT EXISTS default_folders(
  id INTEGER PRIMARY KEY,
  relative_dest_path VARCHAR(255) NOT NULL,
  absolute_orig_path VARCHAR(255),
  comment VARCHAR(255),
  UNIQUE(relative_dest_path, absolute_orig_path)
);
--
CREATE TABLE IF NOT EXISTS folders(
  id INTEGER PRIMARY KEY,
  language_id INTEGER NOT NULL,
  relative_dest_path VARCHAR(255) NOT NULL,
  comment VARCHAR(255),
  FOREIGN KEY(language_id) REFERENCES languages(id),
  UNIQUE(language_id, relative_dest_path)
);
--
CREATE TABLE IF NOT EXISTS default_files(
  id INTEGER PRIMARY KEY,
  relative_dest_path VARCHAR(255) NOT NULL,
  absolute_orig_path VARCHAR(255),
  comment VARCHAR(255),
  UNIQUE(relative_dest_path, absolute_orig_path)
);
--
CREATE TABLE IF NOT EXISTS files(
  id INTEGER PRIMARY KEY,
  language_id INTEGER NOT NULL,
  is_template BOOLEAN NOT NULL CHECK (is_template IN (0, 1)),
  relative_dest_path VARCHAR(255) NOT NULL,
  absolute_orig_path VARCHAR(255),
  comment VARCHAR(255),
  FOREIGN KEY(language_id) REFERENCES languages(id),
  UNIQUE(
    language_id,
    relative_dest_path,
    absolute_orig_path
  )
);
--
-- INSERT INTO TABLES
--
INSERT INTO
  languages(name)
VALUES
  ("C-Plus-Plus"),
  ("C"),
  ("Python"),
  ("Web");
--
INSERT INTO
  languages_short(language_id, name_short)
VALUES
  (1, "cpp"),
  (1, "c++"),
  (1, "cc"),
  (2, "c"),
  (3, "python"),
  (3, "py"),
  (4, "html"),
  (4, "php"),
  (4, "javascript"),
  (4, "js"),
  (4, "typescript"),
  (4, "ts"),
  (4, "web");
--
INSERT INTO
  default_folders(relative_dest_path, absolute_orig_path)
VALUES(".vscode", "templates/vscode/");
--
INSERT INTO
  folders(language_id, relative_dest_path)
VALUES
  (1, "build/debug/"),
  (1, "build/release/"),
  (1, "doc/"),
  (1, "include/"),
  (1, "lib/"),
  (1, "src/"),
  (1, "test/"),
  (2, "build/debug/"),
  (2, "build/release/"),
  (2, "doc/"),
  (2, "include/"),
  (2, "lib/"),
  (2, "src/"),
  (2, "test/"),
  (3, "src/"),
  (3, "doc/"),
  (3, "test/"),
  (4, "public_html/css/"),
  (4, "public_html/img/"),
  (4, "public_html/js/"),
  (4, "public_html/fonts/"),
  (4, "public_html/include/"),
  (4, "resources/library"),
  (4, "resources/templates/");
--
INSERT INTO
  default_files(relative_dest_path, absolute_orig_path)
VALUES
  (".gitignore", "templates/gitignore"),
  ("README.md", "templates/README.md");
--
INSERT INTO
  files(
    language_id,
    is_template,
    relative_dest_path,
    absolute_orig_path
  )
VALUES
  (
    1,
    1,
    "src/main.cpp",
    "templates/template.cpp"
  ),
  (
    1,
    1,
    "CMakeLists.txt",
    "templates/CMakeLists-cpp.txt"
  ),
  (
    2,
    1,
    "src/main.c",
    "templates/template.c"
  ),
  (
    2,
    1,
    "CMakeLists.txt",
    "templates/CMakeLists-c.txt"
  ),
  (3, 0, "./src/__main__.py", NULL),
  (3, 0, "./src/__init__.py", NULL),
  (
    4,
    1,
    "public_html/index.php",
    "templates/template.php"
  ),
  (4, 0, "public_html/css/main.css", NULL),
  (4, 0, "public_html/css/util.css", NULL),
  (4, 0, "public_html/js/main.js", NULL);