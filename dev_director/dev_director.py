#!/usr/bin/env python3

import os
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
                    "Syntax: dev_director <projectname> <language>")
        sys.exit(1)

    if not str(sys.argv[1]).strip():
        print_error("Error: Projectname needs to be specified.")
        sys.exit(2)
    if not str(sys.argv[2]).strip():
        print_error("Error: Language needs to be specified.")
        sys.exit(3)


class DevDirector:

    # Path to the database
    db = "~/.config/dev_director/db/dd.sqlite"

    def __init__(self, project_name, lang):
        self.project_name = project_name
        self.lang = lang

    def run(self):
        # Check if provided language is supported
        self.is_lang_supported()

        # Create the project folder
        self.create_project_folder()

        # Create language specific sub folders
        self.create_sub_folders()

    def does_dir_exist(self):
        ''' Check if directory already exists. '''

        if os.path.exists(self.project_name):
            print_error("Error: Directory already exists.")
            sys.exit(2)

    def is_lang_supported(self):
        ''' Check if the provided language is supported. '''

        langs = []

        if not self.lang in langs:
            print_error("You have to specify a supported language.",
                        ("Currently supported languages: " + str(langs)))
            sys.exit(3)

    def create_project_folder(self):
        ''' Create the main project directory. '''

        # Create main directory
        print("Creating project folder...")
        subprocess.run(["mkdir", "-p", self.project_name])

        # Create ReadMe
        print("Creating ReadMe...")
        cmd = 'echo "# ' + self.project_name + '" > ' + self.project_name + '/README.md'
        subprocess.run(cmd, shell=True)

    def create_sub_folders(self):
        ''' Create sub folders depending on the specified language. '''

        lang_folders = {}

        # Create subfolders
        print("Creating subfolders...")

        for langs, folders in lang_folders.items():
            if self.lang in langs:
                for folder in folders:
                    folder = self.project_name + '/' + folder
                    subprocess.run(["mkdir", "-p", folder])


def main():
    ''' Main function. '''
    # Check cli args
    are_args_valid(sys.argv)

    # Define name and language
    project_name = sys.argv[1]
    lang = str(sys.argv[2]).lower()

    # Create and run the dev_director
    dd = DevDirector(project_name, lang)
    dd.run()


if __name__ == '__main__':
    main()
