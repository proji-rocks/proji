#!.env/bin/python3

import sqlite3

import pytest

from create_project.create_project import CreateProject


"""

FIXTURES

"""


@pytest.fixture
def valid_db_conn():
    # Db Connection
    return sqlite3.connect(CreateProject.db)


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

    # CWD
    cwd = "/tmp/create_project/valid/"

    # Valid
    cps.append(CreateProject("cpp", f"{cwd}Project1"))
    cps.append(CreateProject("py", f"{cwd}Project2"))
    cps.append(CreateProject("js", f"{cwd}Project3"))

    # Add connection to cp instances
    for cp in cps:
        cp.conn = conn
        cp.cur = cp.conn.cursor()
    return cps


@pytest.fixture
def invalid_cp_languages(valid_db_conn):
    # Db Connection
    conn = valid_db_conn

    # List of instances
    cps = []

    # CWD
    cwd = "/tmp/create_project/invalid/"

    # Invalid
    cps.append(CreateProject("x", f"{cwd}Project1"))
    cps.append(CreateProject("fk", f"{cwd}Project2"))
    cps.append(CreateProject("", f"{cwd}Project3"))

    # Add connection to cp instances
    for cp in cps:
        cp.conn = conn
        cp.cur = cp.conn.cursor()
    return cps


@pytest.fixture
def invalid_cp_folders(valid_db_conn, valid_cps):
    # Db Connection
    conn = valid_db_conn

    # List of instances
    cps = []

    # Invalid
    cps.append(CreateProject("py", "/test_project123"))
    cps.append(CreateProject("cpp", ""))

    # Add connection to cp instances
    for cp in cps:
        cp.conn = conn
        cp.cur = cp.conn.cursor()
    return cps


"""

TESTS

"""


def test_init():
    lang = "lang"
    name = "name"
    cp = CreateProject(lang, name)
    assert cp.lang == lang
    assert cp.project_name == name


def test_init_types():
    with pytest.raises(TypeError):
        CreateProject(1, 2)
        CreateProject(True, False)
        CreateProject("Test", True)
        CreateProject(True, "Test")


def test_db_conn(valid_db_conn):
    assert valid_db_conn


def test_lang_supported(valid_cps, invalid_cp_languages):
    # Supported languages
    for valid_cp in valid_cps:
        assert valid_cp._is_lang_supported()

    # Not supported languages
    for invalid_cp in invalid_cp_languages:
        assert not invalid_cp._is_lang_supported()


def test_create_project_folder(valid_cps, invalid_cp_folders):

    for valid_cp in valid_cps:
        assert valid_cp._create_project_folder()

    for invalid_cp in invalid_cp_folders:
        assert not invalid_cp._create_project_folder()


def test_create_sub_folder(valid_cps):
    for valid_cp in valid_cps:
        assert valid_cp._create_sub_folders()


def test_create_files(valid_cps):
    for valid_cp in valid_cps:
        assert valid_cp._create_files()


def test_copy_templates(valid_cps):
    for valid_cp in valid_cps:
        assert valid_cp._copy_templates()
