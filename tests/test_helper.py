#!.env/bin/python3

import pytest

from create_project.helper import Helper


def test_format_err_msg_output():
    assert Helper.format_err_msg("TestCase") == "Error: TestCase\n"
    assert Helper.format_err_msg("T1", "T2", "T3") == "Error: T1\nT2\nT3\n"


def test_format_err_msg_values():
    with pytest.raises(ValueError):
        Helper.format_err_msg()


def test_format_err_msg_types():
    with pytest.raises(TypeError):
        Helper.format_err_msg(True)
        Helper.format_err_msg(123)
        Helper.format_err_msg("Test", 123)
