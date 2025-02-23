# SSH Select

SSH Select is a simple command-line tool written in Go that allows users to easily select and connect to SSH servers using a standard SSH config format.

[![CI](https://github.com/mmarkk/sshselect/actions/workflows/ci.yml/badge.svg)](https://github.com/mmarkk/sshselect/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmarkk/sshselect)](https://goreportcard.com/report/github.com/mmarkk/sshselect)

## Features

- Interactive selection of SSH connections with server names
- Standard SSH config file format for easy configuration
- Fuzzy search through server names
- Arrow key and number-based selection
- Exit option in the selector
- Man page documentation

## Requirements

- Go 1.21.8 or higher
- Make (for installation)
- github.com/manifoldco/promptui v0.9.0 (automatically installed when using `go mod tidy`)

## Installation

### System-wide Installation (requires sudo)

```bash
git clone https://github.com/mmarkk/sshselect.git
cd sshselect
make
sudo make install
```

This installs:
- Binary to /usr/local/bin
- Man pages to /usr/local/share/man

### User-local Installation

```bash
git clone https://github.com/mmarkk/sshselect.git
cd sshselect
make
make install-user
```

This installs:
- Binary to ~/.local/bin
- Man pages to ~/.local/share/man

Ensure these directories are in your PATH/MANPATH:
```bash
# Add to your shell configuration file (e.g., .bashrc, .zshrc):
export PATH="$HOME/.local/bin:$PATH"
export MANPATH="$HOME/.local/share/man:$MANPATH"
```

## Configuration

SSH connections are configured in `~/.config/sshselect/config` using standard SSH config format:

```text
# Example configuration
Host webserver
    HostName 192.168.1.100
    User admin
    Port 2222  # Optional, defaults to 22

Host aws-instance
    HostName ec2-1-2-3-4.compute-1.amazonaws.com
    User ubuntu
```

The config file will be created automatically on first run with example entries.

## Usage

1. Run the program:
    ```bash
    sshselect
    ```

2. Use the interface:
    - Up/Down arrows to navigate
    - Type to search server names
    - Enter to connect
    - Select "Exit" or press Ctrl+C to quit

## Documentation

Manual pages are available after installation:
```bash
man sshselect
```

## Makefile Commands

- `make build`: Builds the sshselect binary
- `make install`: System-wide installation (requires sudo)
- `make install-user`: User-local installation
- `make uninstall`: Remove system-wide installation
- `make uninstall-user`: Remove user-local installation
- `make clean`: Remove built files

## Error Handling

- Invalid config entries will be skipped with warnings
- SSH connection errors will be reported
- Config file will be created if not found

## Security

When reporting issues or sharing configurations, be mindful not to expose sensitive information:
- Avoid sharing your actual SSH hostnames, IP addresses, or usernames
- Use example domains (like example.com) or placeholder IPs (like 192.0.2.x) in examples
- Review your config files for sensitive data before sharing
- The default config location (~/.config/sshselect/config) is local to your machine and not included in the repository

### Security Features
- Secure file permissions (0750 for directories, 0600 for config files)
- Path validation to prevent directory traversal attacks
- Strict SSH command validation to prevent command injection
- Detailed error messages for security-related issues

### Code Scanning

This repository uses GitHub's code scanning features to maintain security:
1. Go to repository Settings > Security > Code security and analysis
2. Enable "Code scanning"
3. Choose "CodeQL Analysis" as the scanning engine
4. Configure the default CodeQL workflow

This ensures automated security analysis on every push and pull request.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

Copyright © 2025 Mark McKenzie
