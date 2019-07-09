#!/usr/bin/env python

import sys
import subprocess
import os


def print_error(*err_msg):
    ''' Print an error. '''

    err_out = "Error: "

    for err in err_msg:
        err_out += (err + '\n')

    print(err_out)


def check_args(args):
    ''' Check number of cli arguments. '''

    if len(args) != 3:
        print_error("Missing arguments.",
                    "Syntax: dev_director <projectname> <language>")
        sys.exit(1)


def check_dir(folder):
    ''' Check if a valid dir/project name was provided. '''

    if os.path.exists(folder):
        print_error("Directory already exists.")
        sys.exit(2)


def check_lang(lang):
    ''' Check if the provided language is supported. '''

    langs = ["c", "cpp", "c++", "cc", "python", "py", "web", "php",
             "js", "javascript", "ts", "typescript", "html", "css"]

    if not lang in langs:
        print_error("You have to specify a supported language.",
                    ("Currently supported languages: " + str(langs)))
        sys.exit(3)


def create_project_base(project_name):
    ''' Create the main project directory. '''

    # Create main directory
    print("Creating project folder...")
    subprocess.run(["mkdir", "-p", project_name])

    # Create ReadMe
    print("Creating ReadMe...")
    cmd = 'echo "# ' + project_name + '" > ' + project_name + '/README.md'
    subprocess.run(cmd, shell=True)


def create_sub_folders(project_name, lang):
    ''' Create sub folders depending on the specified language. '''

    lang_folders = {("c", "cpp", "c++", "cc"): ("build/debug/", "build/release/", "doc/", "include/" + project_name + "/", "lib/", "src/", "test/"), ("php", "js", "javascript", "ts", "typescript", "html",
                                                                                                                                                      "css", "web"): ("public_html/css/", "public_html/img/content/", "public_html/img/layout/", "public_html/js/", "resources/library/", "resources/templates/"), ("python", "py"): ("src/", "docs/", "tests/")}

    # Create subfolders
    print("Creating subfolders...")

    for langs, folders in lang_folders.items():
        if lang in langs:
            for folder in folders:
                folder = project_name + '/' + folder
                subprocess.run(["mkdir", "-p", folder])


def main():
    ''' Main function. '''
    args = sys.argv
    project_name = sys.argv[1]
    lang = str(sys.argv[2]).lower()

    # Check cli args
    check_args(args)

    # Check if directory already exists
    check_dir(project_name)

    # Check if supported language was provided
    check_lang(lang)

    # Create the project folder
    create_project_base(project_name)

    # Create language specific sub folders
    create_sub_folders(project_name, lang)


if __name__ == '__main__':
    main()
