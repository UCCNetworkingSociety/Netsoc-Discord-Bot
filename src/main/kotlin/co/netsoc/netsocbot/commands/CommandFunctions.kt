package co.netsoc.netsocbot.commands

import co.netsoc.netsocbot.*
import co.netsoc.netsocbot.utils.*
import com.jessecorbett.diskord.api.model.Message
import java.io.File

fun help(message: Message): String {
    var out = "```"
    for (command in helpStrings.keys) {
        out += "${PREFIX}${command}: ${helpStrings[command]}\n"
    }
    return out + "```"
}

fun ping(message: Message): String {
    return "pong!"
}

suspend fun register(message: Message): Unit {
    val author = message.author
    if (message.guildId != null) {
        if (guilds.containsKey(author.id)) {
            guilds[author.id] += "," + message.guildId!!
        } else {
            guilds[author.id] = message.guildId!!
        }
        messageUser(author, "Please message me your UCC email address so I can verify you as a member of UCC")
    }
}

fun brownie(message: Message): String {
    val receiver = message.usersMentioned.firstOrNull()
    val author = message.author

//    gonna use a file cause i dont know if im allowed to use a databse
    var leaderboardFile = File("brownieLeaderboard.txt")
    var yes = leaderboardFile.readLines()
    if (receiver != null && author.id != receiver.id)  {
        var fileOut = ""
        var triggered = false

//        checks file for the id of te user
        for (id in yes) {
            var id2 = id.split(" ")
            var count = id2.elementAt(1).toInt()
            if (id2.elementAt(0) == receiver.id){
                count += 1
                triggered = true
            }
            fileOut += id2.elementAt(0) + " " + count.toString() + "\n"
        }

//        this triggers if the users id isnt in the file and writes to anew line in the file
        if (!triggered){
            fileOut += receiver.id + " 1\n"
        }
        leaderboardFile.writeText(fileOut)

        return receiver.username + " has been given a brownie point"
    }
//    this part triggers if someone hasnt been brownied and only calls the command brownie, they get to see their brownie score
    var count = 0
    for (id in yes) {
        var id2 = id.split(" ")
        if (id2.elementAt(0) == author.id){
            count = id2.elementAt(1).toInt()
            break
        }
    }
    return author.username + " has " + count.toString() + " brownie points"
}