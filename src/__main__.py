#!.env/bin/python3

import sys

from create_project import CreateProject
from helper import are_args_valid


def main():
    ''' Main function. '''
    # Check cli args
    are_args_valid(sys.argv)

    # Define name and language
    project_name = sys.argv[1]
    lang = str(sys.argv[2]).lower()

    cp = CreateProject(project_name, lang)
    cp.run()


if __name__ == "__main__":
    main()
