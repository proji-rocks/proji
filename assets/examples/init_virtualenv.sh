#!/bin/sh

# Create virtualenv
echo "> Creating virtualenv"
virtualenv --quiet .env
echo .env >>.gitignore
source .env/bin/activate

# Install python packages
echo "> Installing python packages"
pip install --quiet pylint pep8 black pytest
