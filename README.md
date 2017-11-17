## Discord Broodwar Bot

### Build

This assumes you already have a working Go environment setup and that
DiscordGo is correctly installed on your system.


From within the pingpong example folder, run the below command to compile the
example.

```sh
go build
```

### Usage

Get your token from
[Discord Apps Page](https://discordapp.com/developers/applications/me).

```
cat TOKEN | ./discordbw
```

The below example shows how to start the bot

```sh
./discordbw  -t YOUR_BOT_TOKEN
Bot is now running.  Press CTRL-C to exit.
```
