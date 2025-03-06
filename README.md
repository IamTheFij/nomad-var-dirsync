# Nomad Var DirSync

Nomad Var DirSync is a command-line tool designed to synchronize directories with Nomad variables. It allows you to write the contents of a directory to Nomad variables and read Nomad variables back into a directory.

## Installation

To install Nomad Var DirSync, you need to have Go installed on your machine. You can download and install Go from [the official website](https://golang.org/dl/).

Once you have Go installed, you can install Nomad Var DirSync by running the following command:

```sh
go install git.iamthefij.com/iamthefij/nomad-var-dirsync@latest
```

## Usage

Nomad Var DirSync provides two main actions: `write` and `read`.

### Write Directory to Nomad Variables

To write the contents of a directory to Nomad variables, use the `write` action:

```sh
nomad-var-dirsync -root-var=<root-variable-path> write <source-directory>
```

- `-root-var`: The root path for the Nomad variable.
- `<destination-directory>`: The path to the directory you want to write to Nomad variables.

### Read Nomad Variables to Directory

To read Nomad variables back into a directory, use the `read` action:

```sh
nomad-var-dirsync -root-var=<root-variable-path> -dir-perms=<permissions> read <target-directory>
```

- `-root-var`: The root path for the Nomad variable.
- `-dir-perms`: (Optional) Default permissions for new directories (default: `0o777`).
- `<target-directory>`: The path to the directory where you want to read the Nomad variables.


## Environment Variables

- `NOMAD_ADDR`: The address of the Nomad server (default: `http://localhost:4646`).
- `NOMAD_TOKEN`: The secret ID token for authenticating with the Nomad API.

This should also support the same environment variables as the official [Nomad CLI](https://www.nomadproject.io/docs/commands/cli).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
