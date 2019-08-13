#!.env/bin/python3

import pytest

from src.create_project import CreateProject


def test_init():
    name = "name"
    lang = "lang"
    cp = CreateProject(name, lang)
    assert cp.get_project_name() == name
    assert cp.get_language() == lang


def test_init_types():
    with pytest.raises(TypeError):
        CreateProject(1, 2)
        CreateProject(True, False)
        CreateProject("Test", True)
        CreateProject(True, "Test")
