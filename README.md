# USDT Monitor for Telegram

This is a simple daemon that monitors USDT transactions on the Bitcoin network and sends notifications to a Telegram bot.

## Features

- Monitors USDT transactions on the Bitcoin network
- Sends notifications to a Telegram bot
- Supports multiple addresses and custom thresholds

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/jixiaod/usdt-monitor.git
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build the binary:
   ```
   go build -o usdt-monitor
   ```

4. Run the binary:
   ```
   ./usdt-monitor
   ```

## Configuration    

Create a `config.json` file with the following structure:

```json
{
    "botToken": "5995992793:AAHuHApwCNShyNTZsSudrJm1uGjhpWXXXX",
    "chatID": -1002248200000,
    "usdtAPIBaseURL": "https://apilist.tronscanapi.com",
    "addresses": [
        {"name": "Address 1", "address": "TScqqrtVweHUozbx6v9H1Y9gJjvuXXXXXX"},
        {"name": "Address 2", "address": "TKzi5ymQUretLpBn5AeXyeQhdmoXXXXXXX"}
    ],
    "threshold": 0.01
}
```

