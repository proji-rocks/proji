#!.env/bin/python3

import unittest

from src.helper import Helper


class TestHelper(unittest.TestCase):
    def test_format_err_msg_output(self):
        self.assertEqual(
            Helper.format_err_msg("TestCase"),
            "Error: TestCase\n")
        self.assertEqual(
            Helper.format_err_msg("Test1", "Test2", "Test3"),
            "Error: Test1\nTest2\nTest3\n")

    def test_format_err_msg_values(self):
        self.assertRaises(ValueError, Helper.format_err_msg)

    def test_format_err_msg_types(self):
        self.assertRaises(TypeError, Helper.format_err_msg, True)
        self.assertRaises(TypeError, Helper.format_err_msg, 123)
        self.assertRaises(TypeError, Helper.format_err_msg, "Test", 123)
