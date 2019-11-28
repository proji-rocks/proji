<div align="center">
<img width=300px height=300px src="assets/images/proji-logo.svg" alt="proji logo">

<h2 align="center">proji</h2>

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nikoksr/proji?sort=semver)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/643b7cce9fd2491e9fde38de6e1c58ad)](https://www.codacy.com/manual/nikoksr/proji?utm_source=github.com&utm_medium=referral&utm_content=nikoksr/proji&utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikoksr/proji)](https://goreportcard.com/report/github.com/nikoksr/proji)
[![CircleCI](https://circleci.com/gh/nikoksr/proji/tree/master.svg?style=svg&circle-token=437a39b49c4fbc9656f7aed86aea369d584ecb87)](https://circleci.com/gh/nikoksr/proji/tree/master)

</div>

* * *

<p align="center">Proji is a simple and fast project creator and manager.</p><br />

## Table Of Contents

-   [About](#about)
-   [Demo](#demo)
-   [Getting Started](#getting_started)
    -   [Dependencies](#dependencies)
    -   [Installation](#installation)
    -   [Running the Tests](#running_the_tests)
    -   [Tab Completion](#tab_completion)
-   [Basic Usage](#basic_usage)
    -   [Setting up a Class](#setting_up_a_class)
    -   [Creating our first projects](#creating_our_first_projects)
-   [Advanced Usage](#advanced_usage)
    -   [Class](#au_class)
    -   [Project](#au_project)
    -   [Status](#au_status)
-   [Contribute](#contribute)
-   [Credits](#credits)
-   [License](#license)

## About <a name = "about"></a>

Proji helps you start new projects quickly and easily. With a single command, it creates and configures complex project directories in seconds, which would take minutes to do manually. Proji creates directories and files for you that are either completely new or copied from a template. For example, why write a completely new ReadMe for each project when you can use a template that you only need to adapt to your project? Furthermore, after creating the project directory, proji can also execute scripts for you that initialize the tools that are relevant to your work and your project (e.g. git).

Proji is not just for software developers, but for anyone who regularly creates new projects and wants to save time and effort. Developers, authors, artists or just for your general organization.... it doesn't matter, because you define what you need.

Proji increases your efficiency, simplifies your workflow and improves the structure of your project directories.

## Demo <a name = "demo"></a>

Quick comparison between proji and a common way of project creation. In the demo we create a python project which consists of three folders and two files. We create a `virtualenv`, install three python packages in it and finally we initialize the project as a `git repo`, make the first commit and add a develop branch.

Common way:

<p align="center">
  <a href="" rel="noopener">
    <img src="assets/gifs/regular-project-demo.gif" alt="Create a go example project">
  </a>
</p>

Proji:

<p align="center">
  <a href="" rel="noopener">
    <img src="assets/gifs/proji-project-demo.gif" alt="Create a go project">
  </a>
</p>

As you can see **proji would have saved you about 1:15 minutes and 230 keystrokes**. And remember, this would have been only for one project and it would not have been not a very complex one either.

## Getting Started <a name = "getting_started"></a>

Proji is currently only supported under linux and a work in progress. You can either download a pre-compiled binary from the latest [release](https://github.com/nikoksr/proji/releases) or install it from source.

Might work under Mac but it's not tested yet.

### Dependencies <a name = "dependencies"></a>

-   [go](https://golang.org/) - Main language
-   [sqlite3](https://www.sqlite.org/index.html) - Database
-   [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - Go Sqlite3 Driver
-   [spf13/cobra](https://github.com/spf13/cobra) - CLI commands
-   [spf13/viper](https://github.com/spf13/viper) - Manage config file
-   [BurntSushi/toml](https://github.com/BurntSushi/toml) - Go toml parser
-   [toml-lang/toml](https://github.com/toml-lang/toml) - Config Language
-   [jedib0t/go-pretty](https://github.com/jedib0t/go-pretty) - CLI Styling
-   [stretchr/testify](https://github.com/stretchr/testify) - Test Framework
-   [gocolly/colly](https://github.com/gocolly/colly) - Web Scraping
-   [briandowns/spinner](https://github.com/briandowns/spinner) - Terminal Spinner

### Installation <a name = "installation"></a>

#### Binary Distributions

1.  Download the latest [release](https://github.com/nikoksr/proji/releases) for your system
2.  Extract the tar with: `$ tar -xf proji-XXXX-XXX.tar.gz`
3.  Run the installer: `$ ./install.sh`

#### Install From Source

1.  `$ go get -u github.com/nikoksr/proji`
2.  `$ go get -v -t -d ./...`
3.  `$ go install ./cmd/proji/` or `go build -o proji ./cmd/proji`
4.  `$ ./install.sh`

Validate the success of your installation by executing `$ proji`. The help text for proji should be printed to the cli.

### Running The Tests <a name = "running_the_tests"></a>

-   `$ go vet ./...`
-   `$ go test -v ./...`

### Tab Completion <a name = "tab_completion"></a>

Proji does support tab completion but at the moment you have to set it up yourself. The following instructions were inspired by [kubernetes completion](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enabling-shell-autocompletion).

#### Bash

For tab-completion under bash you first need to install and enable [bash_completion](https://github.com/scop/bash-completion#installation).

You now need to ensure that the proji completion script gets sourced in all your shell sessions.

    # Create the completion file
    ./proji completion bash > ~/.config/proji/completion.bash.inc

    # Make your bash_profile source it
    printf "
      # Proji shell completion
      source '$HOME/.config/proji/completion.bash.inc'
      " >> $HOME/.bash_profile

    # Source it once for immediate completion
    source $HOME/.bash_profile

#### Zsh

This command will create a zsh completion file in your current users default zsh completion folder:

    ./proji completion zsh > "${fpath[1]}/_proji"

## Basic Usage <a name="basic_usage"></a>

Suppose I create python projects on a regular basis and want to have the same directory structure for each of these projects. I would therefore have to execute every command necessary to create the appropriate directories and files and would then have to run tools like git and virtualenv to fully get my usual development environment up and running.

That would not be too bad if you only create a new project every few weeks or months. However, if you want to create new projects more regularly, be it to test something quickly, learn something new, or quickly create an environment to reproduce and potentially solve a problem found on stackoverflow, then this process quickly becomes very tiring.

### Setting up a Class <a name = "setting_up_a_class"></a>

To solve this problem with proji, we first have to create a so-called class. A class in proji defines the structure and behavior for projects of a particular topic (python in this example). It serves as a template through which proji will create new projects for you in the future.

In our case, we want to have the same basic structure for our python projects in the future. So we'll create a class for python. This class will determine which directories and files we always want to get created by proji and which scripts proji should execute after project generation, for example a script for git which automatically initializes the project, creates a develop branch and makes a first commit.

Note that folders and files can either be created new and empty or be copied from a so-called template. In the config folder you can find a template folder (`~/.config/proji/templates/`) in which you can store folders and files that you want to use as templates. In our example we could put a template python file into this folder. The file could contain a very basic python script something like a 'hello world' program. We can tell proji to always copy this file into our newly created python projects. The same goes for folders. The goal of the templates is to save you some more time.

In addition, we can assign scripts to a proji class which will be executed in a desired order after the project directory has been created. Scripts must be saved under `~/.config/proji/scripts/` and can then be referenced by name in the class config.

#### Structure of a Class

-   **Name:** A name that describes the type/topic of the class (e.g. `python`)
-   **Label:** A label that serves as an abbreviation for easily calling the class (e.g. `py`)
-   **Folders:** A list of folders to be created
-   **Files:** A list of files to be created
-   **Scripts:** A list of scripts to run after the project directory has been created

#### Create a Class

There are four ways to create a new class:

##### 1. Config file

Proji offers the possibility to export and import classes through config files. The easiest way to create a new class would be to export the proji sample config and then adapt it to the needs of the class you want to create. To do so execute the command `$ proji class export --example .`.

Proji creates the file [proji-class-example.toml](assets/examples/example-class-export.toml) in the current working directory. If you open this file in a text editor, you will find a richly annotated configuration of an example class. This config should then be adapted according to your needs.

Once the config has been edited and saved, it can be imported using the `$ proji class import --config proji-class-example.toml` (or whatever you named the file) command. Proji then creates a new class based on the imported config.

##### 2. Local directory

Proji can import the structure of a local directory as a class. It will analyze the directory and create a class config based on all sub-directories and files. 

To do so, first create the config with `$ proji class import --directory path/to/dir`. Then edit the resulting config as needed in your editor of choice. For example, remove folders and files that you do not want to add to your class, or add files, folders, and scripts that are missing. When you are done, import the configuration as already seen above with `$ proji class import --config proji-dir.toml`.

##### 3. Remote repository

Similar behaviour as the `directory` example but instead of importing the structure of a local directory you can import the structure of a remote repository. Currently tested and supported code hosting platforms are [github](https://github.com/) and [gitlab](https://gitlab.com).

As an example, let's create a class based on this repository. First, run the command `$ proji class import --remote-repo https://github.com/nikoksr/proji`. Similar to the `directory` flag this will export a class config based on the remote repository. You can then edit the config to your liking and needs and import it afterwards with `$ proji class import --config proji-proji.toml`.

##### 4. Class add command

The third option is to use the `$ proji class add CLASS-NAME [CLASS-NAME...]` command to create one or more classes in an interactive CLI. Proji will query the necessary data for the new class from you and then create the new class based on that data.

The advantage of the config file is that incorrect information can easily be corrected. For example, if you entered a script that does not exist or whose name was simply misspelled, you can easily change the name in the configuration file. This is not possible in the CLI menu. If the entry is incorrect, the creation process must be restarted.

After the class has been created or imported, we can use the command `$ proji class ls` to display a list of our available classes. The command `$ proji class show LABEL [LABEL...]` allows us to display a detailed view of one or more classes.

### Creating our first projects <a name = "creating_our_first_projects"></a>

Now that we have created our python class in proji, we can use it to easily create new projects. A class is created once and is then reused by proji over and over again, and although the process of creating a class might initially seem a bit complex, you will very soon start saving a lot of time and keystrokes and will improve the general structure of your projects.

Assuming our class has been assigned the label `py`, we can create one or more projects with the command `$ proji create py my-py-project-1 my-py-project-2 my-py-project-3`.

<p align="left">
  <a href="" rel="noopener">
 <img src="assets/gifs/create-three-projects.gif" alt="Create projects example"></a>
</p>

And voilà, proji has created three new project directories where you can start your work immediately. The project directories are all built identically, have the same subdirectories and files, and all ran the same scripts.

Take a look at the [python class config](assets/examples/proji-python.toml) and the [git](assets/examples/init_git.sh) and [virtualenv](assets/examples/init_virtualenv.sh) scripts that were used in this example.

## Advanced Usage <a name="advanced_usage"></a>

Help for all commands is also available with `$ proji help`.

### Class <a name = "au_class"></a>

-   Add a class: `$ proji class add NAME`

-   Remove one or more classes: `$ proji class rm LABEL [LABEL...]`

-   Import one or more classes from configs: `$ proji class import --config FILE [--config FILE...]`

-   Import one or more classes from directories: `$ proji class import --directory DIR [--directory DIR...]`

-   Import one or more classes from remote repositories: `$ proji class import --remote-repo URL [--remote-repo URL...]`

-   Export one or more classes: `$ proji class export LABEL [LABEL...]`

-   List all classes: `$ proji class ls`

-   Show details of one or more classes: `$ proji class show LABEL [LABEL...]`

### Project <a name = "au_project"></a>

-   Create one or more projects: `$ proji create LABEL NAME [NAME...]`

-   Add a project: `$ proji add LABEL PATH STATUS`

-   Remove one or more projects: `$ proji rm ID [ID...]`

-   Set new project path: `$ proji set path PATH PROJECT-ID`

-   Set new project status: `$ proji set status STATUS PROJECT-ID`

-   List all projects: `$ proji ls`

-   Clean up project database: `$ proji clean`

### Status <a name = "au_status"></a>

-   Add one or more statuses: `$ proji status add STATUS [STATUS...]`

-   Remove one or more statuses: `$ proji status rm ID [ID...]`

-   List all statuses: `$ proji status ls`

## Contribute <a name = "contribute"></a>

Contributions to the project are highly appreciated. Take a look at the [code of conduct](./CODE_OF_CONDUCT.md) and at the [contribution instructions](./CONTRIBUTING.md) and you should be good to go.

## Credits <a name = "credits"></a>

-   Logo by Alex Potterson (alexeyrdoxi@gmail.com)

## License <a name = "license"></a>

Proji is released under the MIT license. See [LICENSE](LICENSE)
