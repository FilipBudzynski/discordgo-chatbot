# Discord music bot written in golang

A personal project made for the purpose of using YouTube API to play songs directly from the links

[Motivation](https://github.com/FilipBudzynski/discordgo-chatbot?tab=readme-ov-file#motivation)  
[Purpose](https://github.com/FilipBudzynski/discordgo-chatbot?tab=readme-ov-file#purpose)  
[How to run](https://github.com/FilipBudzynski/discordgo-chatbot?tab=readme-ov-file#how-to-run)  
[How to use](https://github.com/FilipBudzynski/discordgo-chatbot?tab=readme-ov-file#how-to-use)  
[Thanks and reagrds](https://github.com/FilipBudzynski/discordgo-chatbot?tab=readme-ov-file#special-thanks)  

---
## Motivation

The Motivation for the project comes from the lack of discord music bots that play songs from YouTube
This is not a deployed discord-bot due to the fact that bots which are using YT api offten get banned and are no more to be used.

---
## Purpose

This bot is ment to be used privatly only. You can download the code and run the executable or use the docker image provided. # not implemented yet

---
## How to run

This project uses yt-dlp and opus so make sure you have them installed.
* yt-dlp: https://github.com/yt-dlp/yt-dlp
* opus:  ```sudo apt-get install -y libopus-dev```

Discord token is read from the .env file with godotenv by the key of `DISCORD_TOKEN`.
* First create a discord bot following this step by step guide: https://discord.com/developers/docs/quick-start/getting-started
* Compile the program from the root repository with ```go build -o discordgo-music-bot .```
* Run the executable ```./discordgo-music-bot``` and play some music!

*(I'm currentlly working on the dockerfile that will let you quickly setup the bot to play your favourite songs.)*

## How to use

* `!play <yt-link>` plays the song you want
* `!queue` shows the current queue
* `!skip` skips the current played song
* `!ping` will sand you back a pong!

## Special thanks
---
Part of the code responsible for streaming the audio to the voice channel of discord was based on and slightly tweeked on the dca code found here: https://github.com/bwmarrin/dca.
Therefore special thanks to the bwmarrin for creating the discordgo package and providing good example of the dca streaming.

