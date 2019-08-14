#!.env/bin/python3

import sys

from .create_project import CreateProject
from .helper import Helper


def main():
    ''' Main function. '''
    # Check cli args
    are_args_valid = Helper.are_args_valid(sys.argv)

    if not are_args_valid:
        return are_args_valid

    # Define name and language
    project_name = sys.argv[1]
    lang = str(sys.argv[2]).lower()

    cp = CreateProject(project_name, lang)
    run_res = cp.run()

    if run_res > 0:
        print(Helper.format_err_msg(
            f"Project creation failed with code {run_res}."))


if __name__ == "__main__":
    main()
