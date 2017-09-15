package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/UCCNetworkingSociety/Netsoc-Discord-Bot/config"
	"github.com/UCCNetworkingSociety/Netsoc-Discord-Bot/logging"
	"github.com/bwmarrin/discordgo"
)

// commMap maps the name of a command to the function which executes the command
var commMap map[string]*command

// HelpCommand is the name of the command which lists the commands available
// and can give information about a specific command.
const HelpCommand = "help"

// command is a function which executes the given command and arguments on
// the provided discord session.
type command struct {
	help string
	exec func(context.Context, *discordgo.Session, *discordgo.MessageCreate, []string) error
}

func init() {
	commMap = map[string]*command{
		"ping": &command{
			help: "Responds 'Pong!' to and 'ping'.",
			exec: pingCommand,
		},
		"setShortcut": &command{
			help: "Sets a shortcut command. Usage: !set keyword url_link_to_resource",
			exec: setShortcutCommand,
		},
	}
}

// pingCommand is a basic command which will responds "Pong!" to any ping.
func pingCommand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	l, ok := logging.FromContext(ctx)
	s.ChannelMessageSend(m.ChannelID, "Pong!")
	if ok {
		l.Infof("Responding 'Pong!' to ping command")
	}
	return nil
}

// setShortcutCommand sets string => string shortcut that can be called later to print a value
func setShortcutCommand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, c []string) error {
	l, ok := logging.FromContext(ctx)
	conf := config.GetConfig()

	if len(c) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Too few arguments supplied. Refer to !help for usage.")
		return fmt.Errorf("Too few arguments supplied for set command")
	} else if len(c) > 3 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments supplied. Refer to !help for usage.")
		return fmt.Errorf("Too many arguments supplied for set command")
	}

	member, _ := s.GuildMember(conf.GuildID, m.Author.ID)
	// Check the user has a role defined in the config for this command
	isAllowed := false
	for _, role := range member.Roles {
		state := s.State

		roleInfo, err := state.Role(conf.GuildID, role)
		if err != nil {
			return fmt.Errorf("failed to retrieve role information: %s", err)
		}

		if stringInSlice(roleInfo.Name, conf.Permissions.Set) {
			isAllowed = true
			break
		}
	}

	if isAllowed {
		commMap[c[1]] = &command{
			help: c[2],
			exec: printShortcut,
		}

		if ok {
			l.Infof("%q is setting a shortcut for [%q] => [%q]", m.Author, c[1], c[2])
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You do not have permissions to use this command.")
		if ok {
			l.Infof("%q is not allowed to execute the set command", m.Author)
		}
	}

	return nil
}

// showHelpCommand lists all of the commands available and explains what they do.
func showHelpCommand(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	if l, ok := logging.FromContext(ctx); ok {
		l.Infof("Responding to help command")
	}
	var out string
	for name, c := range commMap {
		out += fmt.Sprintf("%s: %s\n", name, c.help)
	}
	s.ChannelMessageSend(m.ChannelID, out)
	return nil
}

// isHelpCommand tells you whether the given message arguments are calling the help command.
func isHelpCommand(msgArgs []string) bool {
	if len(msgArgs) < 1 {
		return false
	}
	return msgArgs[0] == HelpCommand
}

// Execute parses a msg and executes the command, if it exists.
func Execute(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, msg string) error {
	args := strings.Fields(msg)
	// the help command is a special case because the help command must loop though
	// the map of all other commands.
	if isHelpCommand(args) {
		if err := showHelpCommand(ctx, s, m, args); err != nil {
			return fmt.Errorf("failed to execute help command: %s", err)
		}
		return nil
	}
	if c, ok := commMap[args[0]]; ok {
		if err := c.exec(ctx, s, m, args); err != nil {
			return fmt.Errorf("failed to execute command: %s", err)
		}
		return nil
	}
	return fmt.Errorf("Failed to recognise the command %q", args[0])
}

// stringInSlice searches for a given value in a flat slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// printShortcut uses the help text of the command to print the shortcut's value
func printShortcut(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	s.ChannelMessageSend(m.ChannelID, commMap[args[0]].help)
	return nil
}
