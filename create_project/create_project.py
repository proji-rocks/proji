#!.env/bin/python3

import os
import sqlite3
import subprocess
import sys

from .helper import Helper


class CreateProject:

    # Path to the database
    conf_dir = os.environ["HOME"] + "/.config/create_project/"
    template_path = conf_dir + "templates/"
    script_path = conf_dir + "scripts/"
    db = conf_dir + "db/cp.sqlite3"

    def __init__(self, lang, project_name):
        if not type(lang) is str:
            raise TypeError("Language name must be a string.")
        if not type(project_name) is str:
            raise TypeError("Project name must be a string.")
        self.lang = lang
        self.project_name = project_name
        self.project_id = None
        self.cwd = os.getcwd()
        self.conn = None
        self.cur = None

    def run(self):
        # Connect to database
        with sqlite3.connect(CreateProject.db) as self.conn:
            if not self.conn:
                err_msg = Helper.format_err_msg("Could not connect to database.")
                print(err_msg)
                return 2

            # Create cursor
            self.cur = self.conn.cursor()

            if not self.cur:
                return 3

            # Check if provided language is supported
            if not self._is_extension_supported():
                return 4

            # Create the project folder
            if not self._create_project_folder():
                return 5

            # Change directory into new project folder
            os.chdir(self.project_name)

            # Create language specific sub folders
            if not self._create_sub_folders():
                return 6

            # Create language specific files
            if not self._create_files():
                return 7

            # Copy template files
            if not self._copy_templates():
                return 8

            # Run custom scripts
            if not self._run_scripts():
                return 9

            # Cd back to old cwd
            os.chdir(self.cwd)
        return 0

    def _is_extension_supported(self):
        """ Check if the provided file extension is supported. """

        try:
            self.cur.execute(
                """
                SELECT
                    project_id,
                    extension
                FROM
                    file_extension
                WHERE
                    extension=?
                """,
                (self.lang,),
            )

            lang = self.cur.fetchone()

        except sqlite3.Error as err:
            print(err)
            return False

        if lang:
            self.project_id = lang[0]
            return True

        # Language not supported
        extensions = self.cur.execute("SELECT extension FROM file_extension").fetchall()

        err_msg = Helper.format_err_msg(
            "You have to specify a supported file extension.",
            ("Currently supported languages: " + str(extensions)),
        )

        print(err_msg)
        return False

    def _create_project_folder(self):
        """ Create the main project directory. """

        print("> Creating project folder...")
        try:
            subprocess.run(["mkdir", self.project_name], timeout=10.0, check=True)
        except (subprocess.CalledProcessError, TimeoutError) as err:
            print(err)
            return False
        return True

    def _create_sub_folders(self):
        """ Create sub folders depending on the specified language. """

        print("> Creating subfolders...")

        try:
            self.cur.execute(
                """
                SELECT
                    target_path
                FROM
                    project_folder
                WHERE
                    (project_id=? OR project_id IS NULL)
                    AND template_name IS NULL
                """,
                (self.project_id,),
            )

            subfolders = self.cur.fetchall()

            for subfolder in subfolders:
                subprocess.run(["mkdir", "-p", subfolder[0]], timeout=10.0, check=True)

        except (sqlite3.Error, subprocess.CalledProcessError, TimeoutError) as err:
            print(err)
            return False
        return True

    def _create_files(self):
        """ Create language specific files. """

        print("> Creating files...")

        try:
            self.cur.execute(
                """
                SELECT
                    target_path
                FROM
                    project_file
                WHERE
                    (project_id=? OR project_id IS NULL)
                    AND template_name IS NULL
                """,
                (self.project_id,),
            )

            files = self.cur.fetchall()

            for file in files:
                subprocess.run(["touch", file[0]], timeout=10.0, check=True)

        except (sqlite3.Error, subprocess.CalledProcessError, TimeoutError) as err:
            print(err)
            return False
        return True

    def _copy_templates(self):
        """ Create language specific files. """

        print("> Copying templates...")

        try:
            self.cur.execute(
                """
                SELECT
                    target_path, template_name
                FROM
                    project_file
                WHERE
                    (project_id=? OR project_id IS NULL)
                    AND template_name IS NOT NULL
                """,
                (self.project_id,),
            )

            template_files = self.cur.fetchall()

            self.cur.execute(
                """
                SELECT
                    target_path, template_name
                FROM
                    project_folder
                WHERE
                    (project_id=? OR project_id IS NULL)
                    AND template_name IS NOT NULL
                """,
                (self.project_id,),
            )

            template_folders = self.cur.fetchall()

            for file in template_files:
                template = CreateProject.template_path + file[1]
                target = file[0]
                subprocess.run(["cp", template, target], timeout=30.0, check=True)

            for folder in template_folders:
                template = CreateProject.template_path + folder[1]
                target = folder[0]
                subprocess.run(["cp", "-r", template, target], timeout=30.0, check=True)

        except (sqlite3.Error, subprocess.CalledProcessError, TimeoutError) as err:
            print(err)
            return False
        return True

    def _run_scripts(self):
        """ Run custom scripts. """

        try:
            self.cur.execute(
                """
                SELECT
                    script_name,
                    run_as_sudo
                FROM
                    project_script
                WHERE
                    project_id is NULL
                    OR project_id=?
                ORDER BY project_id DESC
                """,
                (self.project_id,),
            )

            scripts = self.cur.fetchall()

            for script in scripts:
                script_path = self.script_path + script[0]
                run_as_sudo = bool(script[1])

                if run_as_sudo:
                    subprocess.run(["sudo", script_path], timeout=60.0, check=True)
                else:
                    subprocess.run([script_path], timeout=60.0, check=True)

        except (sqlite3.Error, subprocess.CalledProcessError, TimeoutError) as err:
            print(err)
            return False
        return True
