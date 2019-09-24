-- DROP
DROP INDEX IF EXISTS u_class_idx;
DROP INDEX IF EXISTS u_class_label_idx;
DROP INDEX IF EXISTS u_class_folder_idx;
DROP INDEX IF EXISTS u_class_file_idx;
DROP INDEX IF EXISTS u_class_script_idx;
-- CLASS REGULAR
--DROP INDEX IF EXISTS class_folder_idx;
--DROP INDEX IF EXISTS class_file_idx;
--DROP INDEX IF EXISTS class_script_idx;
-- CREATE
-- PROJECT UNIQUE
CREATE UNIQUE INDEX u_project_idx ON project(install_path);
-- CLASS UNIQUE
CREATE UNIQUE INDEX u_class_idx ON class('name');
CREATE UNIQUE INDEX u_class_label_idx ON class_label(label);
CREATE UNIQUE INDEX u_class_folder_idx ON class_folder(class_id, 'target');
CREATE UNIQUE INDEX u_class_file_idx ON class_file(class_id, 'target');
CREATE UNIQUE INDEX u_class_script_idx ON class_script(class_id, 'name');
-- CLASS REGULAR
--CREATE INDEX class_folder_idx ON class_folder('target', template);
--CREATE INDEX class_file_idx ON class_file('target', template);
--CREATE INDEX class_script_idx ON class_script('name', run_as_sudo);