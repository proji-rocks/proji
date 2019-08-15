#!.env/bin/python3

"""A file with helper functions."""

import sys


class Helper:
    @staticmethod
    def format_err_msg(*err_msg):
        """ Print an error. """

        if not err_msg:
            raise ValueError("You have to specify an error message.")

        if not all(isinstance(i, str) for i in err_msg):
            raise TypeError("The error message must be a string.")

        err_out = "> Error: "

        for err in err_msg:
            err_out += err + "\n"

        return err_out

    @staticmethod
    def are_args_valid(args):
        """ Check number of cli arguments. """

        if len(args) < 3:
            print(
                Helper.format_err_msg(
                    "Missing arguments.",
                    "Syntax: create_project <language> <project name> [more projects]",
                )
            )
            return False

        if not str(sys.argv[1]).strip():
            print(Helper.format_err_msg("Language needs to be specified."))
            return False

        if not str(sys.argv[2]).strip():
            print(Helper.format_err_msg("Projectname needs to be specified."))
            return False

        return True

    @staticmethod
    def create_header(title):
        """ Create an individiual header for every project that's being created. """
        return 50 * "#" + "\n#" + f"\n# Project {title}\n" + "#\n" + 50 * "#" + "\n"
