#!.env/bin/python3

import sys

from .create_project import CreateProject
from .helper import Helper


def main():
    """ Main function. """
    # Check cli args
    are_args_valid = Helper.are_args_valid(sys.argv)

    if not are_args_valid:
        return are_args_valid

    # Define name and language
    lang = sys.argv[1]
    projects = set()

    for i in range(2, len(sys.argv)):
        projects.add(sys.argv[i])

    for project in projects:
        print(Helper.create_header(project))
        cp = CreateProject(lang, project)
        run_res = cp.run()

        if run_res > 0:
            print(f"Project creation failed with code {run_res}.\n")
            continue
        print(f"> Done...\n")


if __name__ == "__main__":
    main()
