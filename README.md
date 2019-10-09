<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="assets/proji-freelogodesign-200x200.png" alt="Project logo"></a>
</p>

<!--<h3 align="center">proji</h3>-->

<div align="center">

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nikoksr/proji?sort=semver)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)
[![CircleCI](https://circleci.com/gh/nikoksr/proji/tree/master.svg?style=svg&circle-token=437a39b49c4fbc9656f7aed86aea369d584ecb87)](https://circleci.com/gh/nikoksr/proji/tree/master)

</div>

---

<p align="center">Proji is a simple and fast project creator and manager.
    <br>
</p>

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Built Using](#built_using)

## About <a name = "about"></a>

I always liked to have orderly project structures but it is so annoying to create all those folders and files each time I want to start a new project. It's so much redundant work that costs me a lot of time. This is why I created proji. Proji automates all of those tasks for you based on templates, configs and scripts you created once.

## Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Installing

A step by step series of examples that tell you how to get a development env running.

Say what the step will be

```
Probably run install.sh
```

End with an example of getting some data out of the system or using it for a little demo.

## Usage <a name="usage"></a>

Let's suppose I want to create a c++ project. Normally, I'd have to execute several commands to create various folders and files, to initialize git, setup cmake etc. But since I always want to use the same structure and tools for my C++ apps, this process can be very well automated.

With proji all I have to do is run a single command:

```
proji create cpp MyProjectName
```

Which results in a project structure like this:

![proji create result](assets/proji-create-result-cpp.png)

Now if I'd want to create a python project:

```
$ proji create py MyPythonProject
```

Which results in a project structure like this:

![proji create result](assets/proji-create-result-python.png)

As you can see the two results differ from another. The project relevant folder structure is different and some files are different too.

In the cpp example proji created a cpp specific `.vscode` folder which sets up the debugger in my editor of choice - `vscode`. Proji created a basic project specific `CMakeLists.txt` which enables me to build a debug or release version of my project with `cmake`.

In the python example proji setup a running `virtuelenv` with all packages of my choice already installed. The `.vscode` folder now contains python specific config files.

In both examples proji created barebone main-files, initialized `git` and checked out the `develop branch`.

## Built Using <a name = "built_using"></a>

- [go](https://golang.org/) - Main language
- [sqlite3](https://www.sqlite.org/index.html) - Database
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - Go Sqlite3 Driver
- [spf13/cobra](https://github.com/spf13/cobra) - CLI commands
- [spf13/viper](https://github.com/spf13/viper) - Manage config file
- [BurntSushi/toml](https://github.com/BurntSushi/toml) - Go toml parser
- [toml-lang/toml](https://github.com/toml-lang/toml) - Config Language
- [jedib0t/go-pretty](https://github.com/jedib0t/go-pretty) - CLI Styling
