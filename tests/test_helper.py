#!.env/bin/python3

import pytest

from src.helper import Helper


def test_format_err_msg_output():
    assert Helper.format_err_msg(
        "TestCase") == "Error: TestCase\n", "Test failed"
    assert (Helper.format_err_msg(
        "Test1", "Test2", "Test3") == "Error: Test1\nTest2\nTest3\n",
        "Test failed")


def test_format_err_msg_values():
    with pytest.raises(ValueError):
        Helper.format_err_msg()


def test_format_err_msg_types():
    with pytest.raises(TypeError):
        Helper.format_err_msg(True)
        Helper.format_err_msg(123)
        Helper.format_err_msg("Test", 123)
