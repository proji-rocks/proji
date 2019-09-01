--
-- DROP TABLES
--
DROP TABLE IF EXISTS project;
DROP TABLE IF EXISTS project_status;
DROP TABLE IF EXISTS class;
DROP TABLE IF EXISTS class_label;
DROP TABLE IF EXISTS class_folder;
DROP TABLE IF EXISTS class_file;
DROP TABLE IF EXISTS class_script;
DROP TABLE IF EXISTS global_folder;
DROP TABLE IF EXISTS global_file;
DROP TABLE IF EXISTS global_script;
--
-- CREATE TABLES
--
CREATE TABLE IF NOT EXISTS project(
  project_id INTEGER PRIMARY KEY,
  'name' TEXT NOT NULL,
  class_id INTEGER REFERENCES class(class_id),
  install_path TEXT,
  install_date TEXT,
  project_status_id INTEGER REFERENCES project_status(project_status_id)
);
--
CREATE TABLE IF NOT EXISTS project_status(
  project_status_id INTEGER PRIMARY KEY,
  project_status TEXT NOT NULL,
  comment TEXT
);
--
CREATE TABLE IF NOT EXISTS class(
  class_id INTEGER PRIMARY KEY,
  'name' TEXT NOT NULL
);
--
CREATE TABLE IF NOT EXISTS class_label(
  class_label_id INTEGER PRIMARY KEY,
  class_id INTEGER NOT NULL REFERENCES class(class_id),
  label TEXT NOT NULL
);
--
CREATE TABLE IF NOT EXISTS class_folder(
  class_folder_id INTEGER PRIMARY KEY,
  class_id INTEGER NOT NULL REFERENCES class(class_id),
  'target' TEXT NOT NULL,
  template TEXT
);
--
CREATE TABLE IF NOT EXISTS class_file(
  class_file_id INTEGER PRIMARY KEY,
  class_id INTEGER NOT NULL REFERENCES class(class_id),
  'target' TEXT NOT NULL,
  template TEXT
);
--
CREATE TABLE IF NOT EXISTS class_script(
  class_script_id INTEGER PRIMARY KEY,
  class_id INTEGER NOT NULL REFERENCES class(class_id),
  'name' TEXT NOT NULL,
  run_as_sudo INTEGER NOT NULL
);
--
CREATE TABLE IF NOT EXISTS global_folder(
  global_folder_id INTEGER PRIMARY KEY,
  'target' TEXT NOT NULL,
  template TEXT
);
--
CREATE TABLE IF NOT EXISTS global_file(
  global_file_id INTEGER PRIMARY KEY,
  'target' TEXT NOT NULL,
  template TEXT
);
--
CREATE TABLE IF NOT EXISTS global_script(
  global_script_id INTEGER PRIMARY KEY,
  'name' TEXT NOT NULL,
  run_as_sudo INTEGER NOT NULL
);