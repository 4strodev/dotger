# Dotger
Dotger is a dotfiles manager inspired by stow.

## Why?
Since I started managging my dotfiles I used different methods. I discovered stow and it works but
I found it unintuitive.

## Usage (under development)

### Assumptions
- I asume that you are working on a linux workstation.
- To make the examples easier take the username of your username as `user`.

Put your dotfiles in a monorepo organized as you want.
Lets take a look to my neovim files.

Create a folder named `.dotfiles` in your home path. If you dont know what is your home path put this in a terminal.
```sh
echo $HOME
```

### Config file
Inside this folder I will create another folder named `neovim`.
In this folder I will move all my neovim files under a `nvim` folder. Now create a file named `.dotger.toml`

```toml
# .dotger.toml
[destination]
path = "/home/user/.config"
mkdir = true # create destination folder if not exists
```

Dotger will take all the files and folders of this *entry* and move it to the destination folder using a symlink.
After linking the config entry in your `/home/user/.config/` will be a symblink named nvim pointing to
`$HOME/.dotfiles/neovim/nvim`

### 

```sh
# link entry
# under $HOME/.dotfiles
dotger link neovim
```

### Environment variables
Dotger is able to parse the config files as go text templates. For the moment you have a helper function
to get environment variables. Look this example

```toml
# .dotger.toml
[destination]
path = "{{ getenv "HOME" }}/.config" # this will result on /home/user/.config
mkdir = true # create destination folder if not exists
```
