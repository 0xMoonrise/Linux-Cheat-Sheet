# Linux Notes

*This guide has been written for 37 1a 1c 50 19 00 11 1d 4f 11 0c 12 16 06 10 16 01 03 42 08 01 17 49 1a 0d 06 1c 1c 18 06 05 0c 01 07 49 04 0c 1f 18 1e 54 26 42 01 0e 05 0c 53 06 04 1c 02 54 02 07 1d 41*

## How to setup an ssh config file

sometimes and specially in fresh linux instalation the directory `.ssh` does not exist, and you must create it with `mkdir ~/.ssh`.

**Step 1:** Setup a key file with `key-gen` command.
```sh
ssh-keygen -f ~/.ssh/KeyName -C "name@host" -N ""
```
<details>

<summary>Notes</summary>

- `-f` Refers to the output file and specifies its path and name.
- `-C` Adds a comment to identify the key. It's commonly used with `name@host` to indicate which server the key belongs to.
- `-N` Specifies the passphrase for the key. Using `-N ""` sets an empty passphrase, so the private key is not protected by a passphrase.
</details>

Output:
```text
Generating public/private ed25519 key pair.
Your identification has been saved in /home/moonrise/.ssh/KeyName
Your public key has been saved in /home/moonrise/.ssh/KeyName.pub
The key fingerprint is:
SHA256:36mx/p/MZzw22dirCDKuyQxIJMZtkUqG0KF6y+gM7Qk name@host
The key's randomart image is:
+--[ED25519 256]--+
|+..o.            |
|o+o..            |
|=+.o             |
|=..              |
|...     S        |
|.=..     . . .   |
|E.+.   o .o o  +o|
|+o .+ o o .+.o.**|
| oo  =.. .+o.oB++|
+----[SHA256]-----+
```

**Step 2:** Get your public key
```
cat ~/.ssh/KeyName.pub
```
output
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIEyCwDnnJSph2i4Nl9vy3kXwzOCOoLPgNZuFJw2RmU/ name@host
```
copy your public key `ssh-ed25519 AAAA... name@host`

**Step 3:** In your remote server dump your public key into `~/.ssh/authorized_keys`
```
echo "ssh-ed25519 AAAAC3... name@host" >> ~/.ssh/authorized_keys
```
If you use `echo` to write your public key into `authorized_keys` file, don't forget to use `>>` or you are going to overwrite the file.

**Step 4:** At this point you can login into your server with `-i` flag.
```
ssh -i ~/.ssh/KeyName user@remote
```
<details>

<summary>Notes</summary>

- `-i` Selects  a  file  from  which  the  identity  (private key) for public key authentication is read.
</details>

**Step 5:** Create a config file to easily login into the server
```
Host remote-server
  HostName 123.123.123.123
  User user
  port 22
  IdentityFile ~/.ssh/KeyName
  IdentitiesOnly yes
```

use a file editor like `nano` or `vim` to dump this file in `~/.ssh/config`

<details>

<summary>Notes</summary>

- `Host` it's the server name
- `HostName` ip remote server
- `User` user name
- `port` port number where server
- `IdentityFile` where you locate the key file
- `IdentitiesOnly` Optional to ensure key file only
</details>

## How to copy text with keyboard
Install `xclip` command, sometimes it's not installed by default

Basic use:
```bash
echo "Hello World" | xclip -sel clip
```
Now you can paste it.

copy an entire file:
```bash
xclip -sel clip /etc/passwd
```

you can also paste with `xclip`

```bash
xclip -o -sel clip
```
dump the clipboard content with:

```bash
xclip -o -sel clip > notes.txt
```

`tmux` it's a multiplexer terminal, it allows to do a lot of tasks and features to improve your capabilities in the terminal. it's a huge tool, do not overwhelm what you might do. Now we are just interest about how to copy a section in a previous prompt.

```bash
tmux
```

Now it should appear a green bar in the bottom, that's meaning a fresh `tmux` instalation and now you can command with `ctrl+b` then `[`. Tmux uses a secuential keybinds. That must dettach your cursor from your command prompt, and now you can move it freely. 

Put your cursor where you want to start to copy and press `scape`, that's the begening of copy, a yellow bar should be apeear. You can also select in a rectangle shape if you press the key `r`. when you are ready to copy just press `enter` and it's done.

This is the most tiny use about `tmux`, but you can have, sessions, panes, windows, in a single terminal, that allows you to enhance your work in the terminal in a better way.

## Basic use git commands

show the state of a git repository
```
git status
```

Here you can see which files have been modified or added. 
```text
On branch main

No commits yet

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	README.md
	go.mod
	go.sum
	main.go

nothing added to commit but untracked files present (use "git add" to track)
```

When you are read to commit your changes in the repository, you can use to add in the stash:

```
git add .
```

And finally commit your changes and make a "check point" in the state of your repository.

```
git commit
```

It will open your default text editor that it store in `$EDITOR`, usually it's `vim` or `nano`. Now you can write a a commit message to document your changes.

### Pushing your changes in your repository server

at this point you have source code and a change tracker, this should be enough to work in local, but if you want to save your work in a remote server or share your project publicly or private, you must your a service provider such as: `github.com`, `gitlab.com`, `bitbucket.com` etc.

**step 1:** setting up your service provider

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

**step 2:** set an ssh key like we did in `ssh-keygen`, get your public key and upload into the service provider like `github.com`, follow to steps:

1. Setting
2. SSH and GPG keys
3. New SSH key
4. Title: foo, Key type: Authentication Key, and paste your public key.

done.

**step 3:**
test connection
```
ssh -T git@github.com
```

ssh github config file
```
Host github.com
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519
    IdentitiesOnly yes
```

(optional): force always over ssh

```
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

troubleshooting and testing connection
```
ssh -vT git@github.com
```
