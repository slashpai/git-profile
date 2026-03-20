# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.1.0] - 2026-03-20

### Added

- CLI commands: `add`, `remove` (`rm`), `list` (`ls`), `use`, `show`, `update`, `verify`
- YAML-based profile storage (`~/.git-profiles.yaml`)
- Per-repo (`--local`) and global (`--global`) profile switching
- `show --remotes` flag to display git remotes
- `verify --show-identity` flag to display full git identity alongside match result
- Interactive prompts with input validation for required fields
- Helpful error messages with profile suggestions when a profile is not found
- Clean profile switching: stale signing config is removed when switching profiles
- `--version` / `-v` flag
- Makefile with `build`, `test`, `fmt`, `vet`, `local-install`, and `help` targets
