# SSH Select

SSH Select is a simple command-line tool written in Go that allows users to easily select and connect to SSH servers from a predefined list.

## Features

-   Presents a numbered list of available SSH connections to the user
-   Allows the user to select a connection by entering its number
-   Executes the selected SSH command

## Requirements

-   Go 1.x or higher
-   Make (for easy installation)

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
3. Build and install the executable:

    ```
    make install
    ```

    This will compile the program and install it to `~/bin/sshselect`.

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
3. Build the executable:
    ```
    go build -o sshselect
    ```
4. (Optional) Move the executable to a directory in your PATH:
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

2. You will see a numbered list of available SSH connections. Enter the number of the connection you want to use.

3. The program will execute the selected SSH command, connecting you to the chosen server.

## Customization

To modify the list of available SSH connections, edit the `logins` slice in the `main` function of the `sshselect.go` file. Add or remove SSH commands as needed, then rebuild and reinstall the program using `make install`.

## Makefile Commands

-   `make build`: Builds the sshselect binary
-   `make install`: Builds and installs the binary to ~/bin
-   `make clean`: Removes the built binary

## Error Handling

-   Invalid user input (e.g., selecting a number that doesn't correspond to a connection) will result in an error message.
-   If there's an error executing the SSH command, the program will exit with an error message.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
