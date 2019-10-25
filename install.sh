#!/usr/bin/env sh

# Create config directory
CONF_DIR="${HOME}/.config/proji/"

mkdir -p "$CONF_DIR"
mkdir -p "${CONF_DIR}db/"
mkdir -p "${CONF_DIR}examples/"
mkdir -p "${CONF_DIR}scripts/"
mkdir -p "${CONF_DIR}templates/"

# Download config files
if ! [ -f "${CONF_DIR}config.toml" ]; then
    curl --silent -o "${CONF_DIR}config.toml" https://raw.githubusercontent.com/nikoksr/proji/master/assets/examples/example-config.toml
fi

if ! [ -f "${CONF_DIR}examples/proji-class.toml" ]; then
    curl --silent -o "${CONF_DIR}examples/proji-class.toml" https://raw.githubusercontent.com/nikoksr/proji/master/assets/examples/example-class-export.toml
fi

# Add shell completion
SHELL_IN_USE=$(basename $(echo "$SHELL"))

if [ "$SHELL_IN_USE" = "bash" ]; then
    ./proji completion bash
    # mv proji-bash-completion destination
    echo "Move the created completion file to your shell's default completion folder."
elif [ "$SHELL_IN_USE" = "zsh" ]; then
    ./proji completion zsh
    # mv proji-zsh-completion destination
    echo "Move the created completion file to your shell's default completion folder."
else
    echo "Shell $SHELL_IN_USE is not supoorted for completion yet."
fi

# Install the binary
sudo install proji /usr/local/bin
