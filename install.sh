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
    ./proji completion bash >~/.config/proji/completion.bash.inc
    printf "
        # Proji shell completion
        source '$HOME/.config/proji/completion.bash.inc'
        " >>$HOME/.bash_profile
    source $HOME/.bash_profile
elif [ "$SHELL_IN_USE" = "zsh" ]; then
    echo $FPATH
    # ./proji completion zsh >"${fpath[1]}/_proji"
else
    echo "Shell $SHELL_IN_USE is not supported for completion."
fi

# Install the binary
sudo install proji /usr/local/bin
