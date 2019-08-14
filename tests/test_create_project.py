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

    # CWD
    cwd = "/tmp/create_project/valid/"

    # Valid
    cps.append(CreateProject(f"{cwd}Project1", "cpp"))
    cps.append(CreateProject(f"{cwd}Project2", "py"))
    cps.append(CreateProject(f"{cwd}Project3", "js"))

    # Add connection to cp instances
    for cp in cps:
        cp._conn = conn
        cp._cur = cp._conn.cursor()

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
    cps.append(CreateProject(f"{cwd}Project1", "x"))
    cps.append(CreateProject(f"{cwd}Project2", "fk"))
    cps.append(CreateProject(f"{cwd}Project3", ""))

    # Add connection to cp instances
    for cp in cps:
        cp._conn = conn
        cp._cur = cp._conn.cursor()

    return cps


@pytest.fixture
def invalid_cp_folders(valid_db_conn, valid_cps):
    # Db Connection
    conn = valid_db_conn

    # List of instances
    cps = []

    # Invalid
    cps.append(CreateProject("/test_project123", "py"))
    cps.append(CreateProject("", "cpp"))

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


def test_db_conn(valid_db_conn):
    assert valid_db_conn


def test_does_dir_exist():
    cp_valid1 = CreateProject("new_project", "py")
    cp_valid2 = CreateProject("i_dont_exist", "js")

    assert not cp_valid1._does_dir_exist()
    assert not cp_valid2._does_dir_exist()

    cp_invalid1 = CreateProject("/bin", "py")
    cp_invalid2 = CreateProject("/tmp", "js")

    assert cp_invalid1._does_dir_exist()
    assert cp_invalid2._does_dir_exist()


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
