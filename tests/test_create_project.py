#!.env/bin/python3

import sqlite3

import pytest

from create_project.create_project import CreateProject


@pytest.fixture
def exmpl_cp():
    # Db Connection
    conn = sqlite3.connect(CreateProject.get_db_path())

    # List of instances
    cps = []

    cps.append(CreateProject("TestProject1", "cpp"))
    cps.append(CreateProject("TestProject2", "py"))
    cps.append(CreateProject("TestProject3", "js"))
    cps.append(CreateProject("TestProject4", "x"))
    cps.append(CreateProject("TestProject5", "fk"))
    cps.append(CreateProject("TestProject6", ""))

    # Add connection to cp instances
    for cp in cps:
        cp._conn = conn
        cp._cur = cp._conn.cursor()

    return cps


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


def test_lang_supported(exmpl_cp):
    # Supported languages
    assert exmpl_cp[0]._is_lang_supported()
    assert exmpl_cp[1]._is_lang_supported()
    assert exmpl_cp[2]._is_lang_supported()

    # Not supported languages
    assert not exmpl_cp[3]._is_lang_supported()
    assert not exmpl_cp[4]._is_lang_supported()
    assert not exmpl_cp[5]._is_lang_supported()
