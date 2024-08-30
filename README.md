# Electrum-Go-Server-Flashing

## Overview
Electrum-Go-Server-Flashing is a tool designed for Bitcoin flashing on wallets that support Electrum servers. This utility integrates with a lightweight Bitcoin node server, specifically `miniBTCD`, to execute transactions that simulate Bitcoin transfers. The flashing method works on wallets that communicate with an Electrum server, allowing for rapid interactions and testing on the Bitcoin network.

## Prerequisites

### Software Dependencies
To use Electrum-Go-Server-Flashing, you must have the following software installed on your system:

- **Go (Golang)**: Ensure you have Go 1.16 or higher installed.
- **miniBTCD**: Your local mock Bitcoin node server.
- **Electrum Wallet**: The wallet should support Electrum server connectivity.
- **Redis**: For handling data and caching (optional).

### Setting Up miniBTCD
`miniBTCD` is the lightweight Bitcoin node server that will interact with your Electrum-compatible wallet. Ensure that you have your local instance of `miniBTCD` running on your machine.

You can run the Bitcoin node server using the following This Link:

[miniBTCD](https://github.com/Thankgod20/miniBTCD)

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/Electrum-Go-Server-Flashing.git
   cd Electrum-Go-Server-Flashing
   ```

2. **Build the project:**

   ```bash
   go build
   ```

3. **Configure Electrum Server:**
   Set up your Electrum wallet to communicate with the Electrum server and point to the `miniBTCD` node that is running locally.

4. **Run the Electrum Flashing Tool:**

   ```bash
   ./ElectrumGoServerFlashing --rpc --url="127.0.0.1:18885" 
   ```

   This command starts the flashing tool, connecting it to your local `miniBTCD` node on port `18885`.

## Features

- **Bitcoin Flashing**: Send simulated Bitcoin transactions to any wallet connected to an Electrum server.
- **Electrum Compatibility**: Works with any wallet supporting Electrum servers.
- **MiniBTCD Integration**: Leverages `miniBTCD` for a fast and lightweight mock node implementation.

## Usage

Once you've started the `miniBTCD` server and your Electrum wallet is connected, you can run the flashing commands. The tool will communicate with the wallet via the Electrum server, sending simulated Bitcoin transactions.

## Example Commands

```bash
./ElectrumGoServerFlashing --rpc --url="127.0.0.1:18885"
```

This command will flash 1.5 Bitcoin to the specified wallet address.

## Notes

- **Test Environment**: This tool is meant for testing and educational purposes. Do **NOT** use this for illegal activities.
- **MiniBTCD Node**: Ensure that your node is synced and configured properly to interact with the Electrum server.
- **Transaction Visibility**: Flashing transactions may appear on connected wallets, but these transactions are simulated.

## Contributing
Contributions are welcome! Feel free to fork the repository and submit pull requests to enhance the functionality.

## License
This project is licensed under the MIT License.
