# miniDNS

miniDNS is a lightweight local DNS server written in Go that allows you to block access to specified websites by refusing DNS queries for them. It also forwards all other DNS queries to Google DNS (8.8.8.8). The project includes both the source code and a pre-built binary.

> **Warning:**
> Modifying your DNS settings can disrupt your internet connection if not restored properly. Use at your own risk. The current implementation is targeted for macOS (using `networksetup` commands) and may require administrator privileges.

---

## Features

- **Block Specific Sites:**
  Specify one or more domains to block. miniDNS will return a refusal for any DNS query matching those sites.

---

## Prerequisites

- **macOS:**
  This implementation uses the `networksetup` command available on macOS.

- **Go:**
  To build from source, you need Go installed.
  Download from [golang.org](https://golang.org/dl/).
  Alternatively, there is a pre-built binary available for macOS.

- **Administrator Privileges:**
  Changing DNS settings typically requires administrator privileges. Run the program with the necessary permissions.

---

## Installation

### Using the Pre-Built Binary

1. Download the pre-built `miniDNS` binary from the [GitHub Releases](https://github.com/lucastomic/miniDNS/releases) page.

2. Place the binary in your desired directory.

3. Make sure the binary is executable:
   ```bash
   chmod +x miniDNS
   ```

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/miniDNS.git
   cd miniDNS
   ```

2. Build the binary:
   ```bash
   go build -o miniDNS
   ```

---
### Creating an Alias

To be able to call the `miniDNS` binary from anywhere, you can create an alias. Add the following line to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
alias miniDNS='/path/to/your/miniDNS'
```

Replace `/path/to/your/miniDNS` with the actual path to the `miniDNS` binary. After adding the alias, reload your shell configuration:

```bash
source ~/.bashrc  # or ~/.zshrc, depending on your shell
```

---


## Usage

To run miniDNS, provide one or more sites (without `http://` or trailing `.`) as arguments. For example, to block `example.com` and `test.com`:

```bash
sudo ./miniDNS example.com test.com
```

### What Happens When You Run miniDNS

1. **Backup DNS Settings:**
   miniDNS fetches your current DNS servers for the Wi-Fi interface and stores them.

2. **Modify DNS Settings:**
   It sets your DNS server to `127.0.0.1` so that all DNS queries go to the local miniDNS server.

3. **Start Local DNS Server:**
   miniDNS listens on UDP port 53 for DNS queries.
   - Queries for forbidden sites are refused.
   - Other queries are forwarded to Google DNS.

4. **Graceful Shutdown:**
   On receiving an interrupt signal (Ctrl+C or termination), the program restores your original DNS settings before exiting.

---

## Troubleshooting

- **Port 53 Already in Use:**
  If you encounter an error indicating that port 53 is in use, ensure no other DNS service is running. You might need to stop conflicting services or run the binary with appropriate privileges.

- **Permission Issues:**
  Modifying DNS settings requires administrator rights. Use `sudo` or run as an administrator.

- **Platform Compatibility:**
  This version of miniDNS is designed for macOS. For other platforms, modifications might be needed, especially regarding DNS configuration commands.

---

## Contributing

Contributions, issues, and feature requests are welcome!
Feel free to check [Issues](https://github.com/lucastomic/miniDNS/issues) and [Pull Requests](https://github.com/lucastomic/miniDNS/pulls).

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Enjoy using miniDNS to manage your browsing experience!
