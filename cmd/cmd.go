// package cmd contains all bindings to message commands, slash commands, etc.
// New commands are added by creating a file in this package, defining the functions, commands
// and creating an init() function which should then add the functions and/or commands to these maps.
// These maps are called in server/server.go.
package cmd

import "github.com/bwmarrin/discordgo"

var (
	// Commands is a map of metadata for slash commands
	// metadata includes things like name, description, etc.
	Commands = map[string]*discordgo.ApplicationCommand{}
	// Handlers is a map of functions for slash commands
	Handlers = map[string]func(bot *discordgo.Session, i *discordgo.InteractionCreate){}
	// MessageCreate is a map of functions for MessageCreate functions
	// MessageCreate refers to commands which are triggered by normal discord messages
	MessageCreate = map[string]func(bot *discordgo.Session, m *discordgo.MessageCreate){}
)
