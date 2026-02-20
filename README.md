# kubectl-plugin-arcane

A kubectl plugin for managing [Arcane](https://github.com/SneaksAndData/arcane-operator) streams.

## Features

- Start, stop, and backfill streams
- Declare and stop downtime for streams
- Integrates with Kubernetes via `kubectl`

## Installation

### Kubectl Plugin Manager (krew)

--TBD--

### Kubectl Plugin Manager (krew) with custom index (pre-release versions)

--TBD--

### Manual

1. **Download the latest release binary:**

   Go to the [releases page](https://github.com/sneaksAndData/kubectl-plugin-arcane/releases) and download the binary for your operating system and architecture.

2. **Move the binary to `~/.local/bin`:**

   ```sh
   mkdir -p ~/.local/bin
   mv path/to/downloaded/kubectl-arcane ~/.local/bin/
   chmod +x ~/.local/bin/kubectl-arcane
   ```

3. **Add `~/.local/bin` to your `PATH` (if not already):**

   Add the following line to your `~/.zshrc` or `~/.bashrc`:
   ```sh
   export PATH="$HOME/.local/bin:$PATH"
   ```
   Then reload your shell:
   ```sh
   source ~/.zshrc  # or source ~/.bashrc
   ```
>[!NOTE]
> The binary is named `kubectl-arcane` to be recognized as a kubectl plugin.
> After installation, you can use it with the `kubectl arcane` command.

4. **Verify installation:**

   ```sh
   kubectl arcane --help
   ```

5. **Unset Quarantine attribute on MacOS if you encounter permission issues:**

   ```sh
   xattr -d com.apple.quarantine ~/.local/bin/kubectl-arcane
   ```

## Usage

This plugin adds the `arcane` command to `kubectl` with the following subcommands:

### Stream Commands

- `kubectl arcane stream start <stream-class> <stream-id>`
Start a stream
 
- `kubectl arcane stream stop <stream-class> <stream-id>`
Stop a stream
 
- `kubectl arcane stream backfill <stream-class> <stream-id> [--wait]`
Run a stream in backfill mode
- `--wait`: Wait for backfill command to complete

### Downtime Commands

- `kubectl arcane downtime declare <stream-class> <prefix> <key>`
Stop the list of streams streams by the name prefix. The `<key>` parameter is used to identify the list of streams
that are in downtime, and will be used to resume the streams when downtime is stopped.

- `kubectl arcane downtime stop <stream-class> <key>`
Stop the downtime by waking up the list of streams that are in downtime by the `<key>` parameter.

## Help

For more information on a command, use:

```sh
kubectl arcane <command> --help
```

## License

See [LICENSE](LICENSE).
