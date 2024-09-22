# SSH Select

SSH Select is a simple command-line tool written in Go that allows users to easily select and connect to SSH servers from a predefined list.

## Features

-   Displays a list of available SSH connections with server names
-   Supports selection using both arrow keys and number input
-   Allows searching through the list of servers
-   Executes the selected SSH command

## Requirements

-   Go 1.21.8 or higher
-   Make (for easy installation)
-   github.com/manifoldco/promptui v0.9.0 (automatically installed when using `go mod tidy`)

## Installation

### Using Make (Recommended)

1. Clone this repository:
    ```
    git clone https://github.com/yourusername/sshselect.git
    ```
2. Change to the project directory:
    ```
    cd sshselect
    ```
3. Download dependencies and build the executable:

    ```
    make install
    ```

    This will run `go mod tidy`, compile the program, and install it to `~/bin/sshselect`.

4. Ensure `~/bin` is in your PATH. Add the following line to your shell configuration file (e.g., `.bashrc`, `.zshrc`):
    ```
    export PATH="$HOME/bin:$PATH"
    ```
    Then, reload your shell configuration or restart your terminal.

### Manual Installation

1. Clone this repository:
    ```
    git clone https://github.com/yourusername/sshselect.git
    ```
2. Change to the project directory:
    ```
    cd sshselect
    ```
3. Download dependencies:
    ```
    go mod tidy
    ```
4. Build the executable:
    ```
    go build -o sshselect
    ```
5. (Optional) Move the executable to a directory in your PATH:
    ```
    mv sshselect ~/bin/
    ```

## Usage

1. Run the `sshselect` program:

    ```
    sshselect
    ```

    If you didn't install it to a directory in your PATH, run it using:

    ```
    ./sshselect
    ```

2. You will see a list of available SSH connections with server names. You can:

    - Use the up and down arrow keys to navigate through the list
    - Type a number to quickly select a connection
    - Start typing to search for a specific server name

3. Press Enter to confirm your selection.

4. The program will execute the selected SSH command, connecting you to the chosen server.

## Customization

To modify the list of available SSH connections, edit the `logins` slice in the `main` function of the `sshselect.go` file. Each entry in the slice should be a `sshLogin` struct with a `Name` and `Command` field. For example:

```go
logins := []sshLogin{
    {Name: "Web Server", Command: "ssh user@webserver.example.com"},
    {Name: "Database", Command: "ssh dbuser@db.example.com -p 2222"},
    // Add more entries as needed
}
```

After making changes, rebuild and reinstall the program using `make install`.

## Makefile Commands

-   `make build`: Builds the sshselect binary
-   `make install`: Runs `go mod tidy`, builds and installs the binary to ~/bin
-   `make clean`: Removes the built binary

## Error Handling

-   Invalid user input will result in an error message.
-   If there's an error executing the SSH command, the program will exit with an error message.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
