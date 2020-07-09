package static

const CompletionHelpMessage = `To load completions:

Bash:

$ source <(proji completion bash)

# To load completions for each activeSession, execute once:
Linux:
  $ proji completion bash > /etc/bash_completion.d/proji
MacOS:
  $ proji completion bash > /usr/local/etc/bash_completion.d/proji

Zsh:

$ source <(proji completion zsh)

# To load completions for each activeSession, execute once:
$ proji completion zsh > "${fpath[1]}/_proji"

Fish:

$ proji completion fish | source

# To load completions for each activeSession, execute once:
$ proji completion fish > ~/.config/fish/completions/proji.fish
`
