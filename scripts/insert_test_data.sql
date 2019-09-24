--
-- INSERT INTO TABLES
--
INSERT INTO
  project_status(project_status, comment)
VALUES
  (
    "active",
    "I'm actively working on this project."
  ),
  (
    "inactive",
    "I stopped working on this project for now."
  ),
  (
    "done",
    "I finished this project. There is nothing to do."
  ),
  (
    "dead",
    "This project is dead. I'm not planning to work on it anytime soon."
  );
--
INSERT INTO
  class('name')
VALUES
  ("c-plus-plus"),
  ("python");
--
INSERT INTO
  class_label(class_id, label)
VALUES
  (1, "cpp"),
  (1, "c++"),
  (1, "cc"),
  (2, "python"),
  (2, "py");
--
INSERT INTO
  class_folder(class_id, 'target', template)
VALUES
  (1, "build/debug", NULL),
  (1, "build/release", NULL),
  (1, "doc", NULL),
  (1, "include", NULL),
  (1, "lib", NULL),
  (1, "src", NULL),
  (1, "test", NULL),
  (2, "src", NULL),
  (2, "doc", NULL),
  (2, "test", NULL);
--
INSERT INTO
  class_file(class_id, 'target', template)
VALUES
  (1, "src/main.cpp", "main.cpp"),
  (
    1,
    "CMakeLists.txt",
    "CMakeLists-cpp.txt"
  ),
  (2, "src/__main__.py", NULL),
  (2, "src/__init__.py", NULL);
--
INSERT INTO
  class_script(class_id, 'name', run_as_sudo)
VALUES
  (2, "init_virtualenv.sh", 0);