#!.env/bin/python3

import sqlite3

import pytest

from create_project.create_project import CreateProject


'''

FIXTURES

'''


@pytest.fixture
def valid_db_conn():
    # Db Connection
    return sqlite3.connect(CreateProject.get_db_path())


@pytest.fixture
def invalid_db_conn():
    # Db Connection
    return sqlite3.connect("i_dont_exist.sqlite")


@pytest.fixture
def valid_cps(valid_db_conn):
    # Db Connection
    conn = valid_db_conn

    # List of instances
    cps = []

    # Valid
    cps.append(CreateProject("TestProject1", "cpp"))
    cps.append(CreateProject("TestProject2", "py"))
    cps.append(CreateProject("TestProject3", "js"))

    # Add connection to cp instances
    for cp in cps:
        cp._conn = conn
        cp._cur = cp._conn.cursor()

    return cps


@pytest.fixture
def invalid_cps(valid_db_conn):
    # Db Connection
    conn = valid_db_conn

    # List of instances
    cps = []

    # Invalid
    cps.append(CreateProject("TestProject4", "x"))
    cps.append(CreateProject("TestProject5", "fk"))
    cps.append(CreateProject("TestProject6", ""))

    # Add connection to cp instances
    for cp in cps:
        cp._conn = conn
        cp._cur = cp._conn.cursor()

    return cps


'''

TESTS

'''


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

def


def test_lang_supported(valid_cps, invalid_cps):
    # Supported languages
    assert valid_cps[0]._is_lang_supported()
    assert valid_cps[1]._is_lang_supported()
    assert valid_cps[2]._is_lang_supported()

    # Not supported languages
    assert not invalid_cps[0]._is_lang_supported()
    assert not invalid_cps[1]._is_lang_supported()
    assert not invalid_cps[2]._is_lang_supported()
