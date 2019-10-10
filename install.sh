#!/usr/bin/env sh

# Create config directory
CONF_DIR="${HOME}/.config/proji/"

mkdir -p "$CONF_DIR"
mkdir -p "${CONF_DIR}db/"
mkdir -p "${CONF_DIR}examples/"
mkdir -p "${CONF_DIR}scripts/"
mkdir -p "${CONF_DIR}templates/"

# Download config files
curl --silent -o "${CONF_DIR}config.toml" https://raw.githubusercontent.com/nikoksr/proji/master/configs/example-config.toml
curl --silent -o "${CONF_DIR}examples/proji-class.toml" https://raw.githubusercontent.com/nikoksr/proji/master/configs/example-class-export.toml

# Add shell completion
SHELL=$(basename $(echo $SHELL))

if [ "$SHELL" = "bash" ]; then
    ./proji completion bash
    # mv proji-bash-completion destination
    echo "Move the created completion file to your shell's default completion folder."
elif [ "$SHELL" = "zsh" ]; then
    ./proji completion zsh
    # mv proji-zsh-completion destination
    echo "Move the created completion file to your shell's default completion folder."
else
    echo "Shell $SHELL is not supoorted for completion yet."
fi

# Install the binary
sudo install proji /usr/local/bin
