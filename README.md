# git-profile

Manage multiple git profiles easily. Switch between different git identities (name, email, GPG signing key) across repositories with a single command.

## Install

```bash
go install github.com/slashpai/git-profile@latest
```

Or build from source:

```bash
git clone https://github.com/slashpai/git-profile.git
cd git-profile
make local-install
```

> [!NOTE]
> `make local-install` builds the binary and copies it to `/usr/local/bin` (requires sudo).

## Usage

### Add a profile

```bash
git-profile add personal
```

You'll be prompted for:

```
  user.name: John Doe
  user.email: john@personal.com
  user.signingkey - GPG key ID (optional, Enter to skip): ABC123DEF
  commit.gpgsign [y/N]: y
  SSH key path, e.g. ~/.ssh/id_ed25519 (optional, Enter to skip): ~/.ssh/id_ed25519_personal
```

> [!TIP]
> To find your GPG key ID, run `gpg --list-secret-keys --keyid-format long`. The key ID is the hex string after the algorithm (e.g. `ed25519/ABC123DEF`).

### List profiles

```bash
git-profile list
```

```shell
PROFILE   NAME            EMAIL              SIGNING KEY  GPG SIGN  SSH KEY
-------   ----            -----              -----------  --------  -------
personal  John Doe   john@personal.com  ABC123DEF    yes       ~/.ssh/id_ed25519_personal
work      John Doe   john@company.com   XYZ789       yes       ~/.ssh/id_ed25519_work
```

### Switch profile (per-repo)

```bash
git-profile use personal
```

### Switch profile (global)

```bash
git-profile use work --global
```

> [!WARNING]
> Using `--global` will overwrite your global git config for `user.name`, `user.email`, `user.signingkey`, and `commit.gpgsign`.

### Show current identity

```bash
git-profile show
```

```shell
Current git identity:
  user.name        = John Doe
  user.email       = john@personal.com
  user.signingkey  = ABC123DEF
  commit.gpgsign   = true
```

Use `--remotes` / `-r` to also display git remotes:

```bash
git-profile show --remotes
```

### Remove a profile

```bash
git-profile remove personal
# or
git-profile rm personal
```

## Config File

Profiles are stored in `~/.git-profiles.yaml`. The file is created automatically the first time you add a profile.

> [!NOTE]
> No config file is needed to get started. If the file doesn't exist, the tool treats it as an empty config. Commands like `list` will show no profiles, and `show` reads directly from `git config` so it works regardless of whether any profiles have been saved.

You can override the path with `--config`:

```bash
git-profile --config /path/to/profiles.yaml list
```

Example config:

```yaml
profiles:
  personal:
    name: John Doe
    email: john@personal.com
    signingkey: ABC123DEF
    gpgsign: true
    sshkey: ~/.ssh/id_ed25519_personal
  work:
    name: John Doe
    email: john@company.com
    signingkey: XYZ789
    gpgsign: true
    sshkey: ~/.ssh/id_ed25519_work
```

## How It Works

- `git-profile use <name>` applies the profile by running `git config --local` (or `--global` with the flag) for `user.name`, `user.email`, `user.signingkey`, and `commit.gpgsign`.

> [!TIP]
> The SSH key field is stored for your reference only and is not applied via git config. It helps you remember which key goes with which profile.
