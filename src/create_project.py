#!.env/bin/python3

import os
import sqlite3
import subprocess
import sys

from src.helper import Helper


class CreateProject:

    # Path to the database
    __conf_dir = "/home/niko/.config/create_project/"
    __db = __conf_dir + "db/cp.sqlite"

    def __init__(self, project_name, lang):
        if type(project_name) != str:
            raise TypeError("Project name must be a string.")
        if type(lang) != str:
            raise TypeError("Language name must be a string.")
        self.__project_name = project_name
        self.__lang = lang
        self.__lang_id = 0
        self.__conn = None
        self.__cur = None

    def run(self):
        # Check if directory already exists
        self.__does_dir_exist()

        # Connect to database
        with sqlite3.connect(CreateProject.__db) as self.conn:
            if not self.conn:
                err_msg = Helper.format_err_msg(
                    "Could not connect to database.")
                print(err_msg)
                sys.exit(1)

            # Create cursor
            self.__cur = self.conn.cursor()

            # Check if provided language is supported
            self.__is_lang_supported()

            # Create the project folder
            self.__create_project_folder()

            # Create language specific sub folders
            self.__create_sub_folders()

            # Create language specific files
            self.__create_files()

            # Copy template files
            self.__copy_templates()

    def get_project_name(self):
        ''' Get the project name. '''
        return self.__project_name

    def get_language(self):
        ''' Get the language. '''
        return self.__lang

    def __does_dir_exist(self):
        ''' Check if directory already exists. '''

        if os.path.exists(self.__project_name):
            err_msg = Helper.format_err_msg("Directory already exists.")
            print(err_msg)
            sys.exit(2)

    def __is_lang_supported(self):
        ''' Check if the provided language is supported. '''

        for lang_short in self.__cur.execute('''
                                            SELECT
                                                language_id,
                                                name_short
                                            FROM
                                                languages_short
                                            '''):
            if lang_short[1] == self.__lang:
                self.lang_id = lang_short[0]
                return

        langs = self.__cur.execute(
            'SELECT name_short FROM languages_short').fetchall()

        err_msg = Helper.format_err_msg(
            "You have to specify a supported language.",
            ("Currently supported languages: " + str(langs)))
        print(err_msg)
        sys.exit(3)

    def __create_project_folder(self):
        ''' Create the main project directory. '''

        # Create main directory
        print("Creating project folder...")
        subprocess.run(["mkdir", "-p", self.__project_name])

    def __create_sub_folders(self):
        ''' Create sub folders depending on the specified language. '''

        # Create subfolders
        print("Creating subfolders...")

        for sub_folder in self.__cur.execute(
            '''
            SELECT
                relative_dest_path
            FROM
                folders
            WHERE
                language_id=?
            ''',
                (self.lang_id,)):

            sub_folder = "./" + self.__project_name + "/" + sub_folder[0]
            subprocess.run(["mkdir", "-p", sub_folder])

    def __create_files(self):
        ''' Create language specific files. '''
        # Create files
        print("Creating files...")

        for file in self.__cur.execute(
            '''
            SELECT
                relative_dest_path
            FROM
                files
            WHERE
                language_id=? and
                is_template=?
            ''',
                (self.lang_id, 0,)):

            file = "./" + self.__project_name + "/" + file[0]
            subprocess.run(["touch", file])

    def __copy_templates(self):
        ''' Create language specific files. '''
        # Create files
        print("Copying templates...")

        for template in self.__cur.execute(
            '''
            SELECT
                relative_dest_path,
                absolute_orig_path
            FROM
                files
            WHERE
                language_id=? and
                is_template=?''',
                (self.lang_id, 1,)):

            dest = self.__project_name + "/" + template[0]
            template = CreateProject.__conf_dir + template[1]
            subprocess.run(["cp", template, dest])
