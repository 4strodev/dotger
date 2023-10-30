# Dotger
Dotger is a dotfiles manager inspired by stow.

## Why?
Since I started managging my dotfiles I used different methods. I discovered stow and it works but
I found it unintuitive.

## Usage (under development)
Put your dotfiles in a monorepo organized as you want.
Lets take a look to my neovim files.

Create a folder named `.dotfiles` in your home path. If you dont know what is your home path put this in a terminal.
```sh
echo $HOME
```

Inside this folder I will create another folder named `neovim`.
In this folder I will move all my neovim files. Now create a file named `.dotger.toml`

```toml
[destination]
path = "~/.config/nvim"
mkdir = true # create destination folder if not exists
recursive = true # create parent folders if it doesn't exists (only valid if mkdir is true)
```

```sh
# link entry
# under $HOME/.dotfiles/neovim
dotger link neovim
```

