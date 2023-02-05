# Peary
Peary is the bot that issues temporary roles for a time period to help players in arrangement of gaming sessions. Especially for games that are meant to played with friends or have small communities.

# What it does
After sending `/play` command player receives role. By this role players can be pinged by those who are looking to play with someone. The role itself is displayed in seperate section on a server and has a stand out color to make it noticable when user writes something in channel. 

Meanwhile user is free to do anything they want to on Discord. No need to wait in LFG voice channel or monitor dedicated LFG channel on the server. Gone are days of asking "are you still up?".

# How to
`peary` is designed to be self-hosted and used on a single server (you can't have the same bot running on multiple servers simultaniously). However, I still tried to make deployment as easy as possible even for non tech savvy people via using `docker`. Sropusly, don't be intimidated by the amount of steps, they are all fairly simple.

**NOTE:** Currently this tutorial assumes you will host the bot under Linux. Mac and BSD *might* work but it's not guaranteed. Windows will not work due to me not knowing how to use docker-related thing called *volumes* there.

1. Create a discord bot and add it to your server. In Scope panel click: "bot", "application.commands". In permission pannel click: "Manage Roles" (or give this permission to the bot after adding it to the server)

2. Install Docker by proceeding to https://docs.docker.com/desktop/, proceed by link respective to your platform (mind the note above) and follow the instructions.

3. Clone this repository.

4. Create file named `.env` inside repository's folder and paste there the following contents:

```
BOT_APP_ID=replace with your bot app id
BOT_TOKEN=replace with your bot oauth2 token
```

5. Run the following commands from the console one-by-one:

`sudo make build` -- builds docker image for the bot

`sudo make run` -- starts newly built bot or restarts if it's already built 

Now your bot should be up and running. If it is so, congratulations!

To stop the bot type `^C` (`Ctrl+C`).

To remove type `sudo make clear`

# How to send bug report
Message me on Discord: Niki W.#4040
