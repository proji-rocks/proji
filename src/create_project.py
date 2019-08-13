#!/usr/bin/env python3

import os
import sqlite3
import subprocess
import sys


def print_error(*err_msg):
    ''' Print an error. '''

    err_out = "Error: "

    for err in err_msg:
        err_out += (err + '\n')

    print(err_out)


def are_args_valid(args):
    ''' Check number of cli arguments. '''

    if len(args) != 3:
        print_error("Missing arguments.",
                    "Syntax: create_project <projectname> <language>")
        sys.exit(1)

    if not str(sys.argv[1]).strip():
        print_error("Error: Projectname needs to be specified.")
        sys.exit(2)
    if not str(sys.argv[2]).strip():
        print_error("Error: Language needs to be specified.")
        sys.exit(3)


class CreateProject:

    # Path to the database
    conf_dir = "/home/niko/.config/create_project/"
    db = conf_dir + "db/dd.sqlite"

    def __init__(self, project_name, lang):
        self.project_name = project_name
        self.lang = lang
        self.lang_id = 0
        self.conn = None
        self.cur = None

    def run(self):
        # Check if directory already exists
        self.__does_dir_exist()

        # Connect to database
        with sqlite3.connect(cp.sqlite) as self.conn:
            if not self.conn:
                print_error("Could not connect to database.")
                sys.exit(1)

            # Create cursor
            self.cur = self.conn.cursor()

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

    def __does_dir_exist(self):
        ''' Check if directory already exists. '''

        if os.path.exists(self.project_name):
            print_error("Directory already exists.")
            sys.exit(2)

    def __is_lang_supported(self):
        ''' Check if the provided language is supported. '''

        for lang_short in self.cur.execute('SELECT language_id, name_short FROM languages_short'):
            if lang_short[1] == self.lang:
                self.lang_id = lang_short[0]
                return

        langs = self.cur.execute(
            'SELECT name_short FROM languages_short').fetchall()

        print_error("You have to specify a supported language.",
                    ("Currently supported languages: " + str(langs)))
        sys.exit(3)

    def __create_project_folder(self):
        ''' Create the main project directory. '''

        # Create main directory
        print("Creating project folder...")
        subprocess.run(["mkdir", "-p", self.project_name])

    def __create_sub_folders(self):
        ''' Create sub folders depending on the specified language. '''

        # Create subfolders
        print("Creating subfolders...")

        for sub_folder in self.cur.execute(
                'SELECT relative_dest_path FROM folders WHERE language_id=?', (self.lang_id,)):
            sub_folder = "./" + self.project_name + "/" + sub_folder[0]
            subprocess.run(["mkdir", "-p", sub_folder])

    def __create_files(self):
        ''' Create language specific files. '''
        # Create files
        print("Creating files...")

        for file in self.cur.execute(
                'SELECT relative_dest_path FROM files WHERE language_id=? and is_template=?', (self.lang_id, 0,)):
            file = "./" + self.project_name + "/" + file[0]
            subprocess.run(["touch", file])

    def __copy_templates(self):
        ''' Create language specific files. '''
        # Create files
        print("Copying templates...")

        for template in self.cur.execute(
                'SELECT relative_dest_path, absolute_orig_path FROM files WHERE language_id=? and is_template=?', (self.lang_id, 1,)):
            dest = self.project_name + "/" + template[0]
            template = CreateProject.conf_dir + template[1]
            subprocess.run(["cp", template, dest])


def main():
    ''' Main function. '''
    # Check cli args
    are_args_valid(sys.argv)

    # Define name and language
    project_name = sys.argv[1]
    lang = str(sys.argv[2]).lower()

    cp = CreateProject(project_name, lang)
    cp.run()


if __name__ == '__main__':
    main()
