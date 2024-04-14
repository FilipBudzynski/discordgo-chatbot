# Discord music bot written in golang

A personal project made for the purpose of using YouTube API to play songs directly from the links

---
## Motivation

The Motivation for the project comes from the lack of discord music bots that play songs from YouTube
This is not a deployed discord-bot due to the fact that bots which are using YT api offten get banned and are no more to be used.

---
## Purpose

This bot is ment to be used privatly only. You can download the code and run the executable or use the docker image provided. # not implemented yet

---
## Usage

Discord token is read from the .env file via godotenv by the name of DISCORD_TOKEN.
After starting the program you can add it to your discord server via this step by step guide: https://discord.com/developers/docs/quick-start/getting-started

---
Part of the code responsible for streaming the audio to the voice channel of discord was based on and slightly tweeked on the dca code found here: https://github.com/bwmarrin/dca.
Therefore special thanks to the bwmarrin for creating the discordgo package and providing good example of the dca streaming.

