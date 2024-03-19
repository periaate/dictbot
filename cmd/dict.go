package cmd

import (
	"fmt"
	"log"

	_ "embed"

	"github.com/bwmarrin/discordgo"
	"github.com/periaate/dict"
)

//go:embed dict.json
var rawDict []byte

var dictMap dict.DictMap

var cmds = []discordgo.ApplicationCommand{
	{
		Name:        "dict",
		Description: "Wiktionary search command.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "The word to search for.",
				Required:    true,
			},
		},
	},
	{
		Name:        "copyright",
		Description: "Prints copyright and legal information",
	},
}

const notice = `This bot serves data which originates from [Wiktionary](<https://en.wiktionary.org/>). The dataset used can be found on [kaikki.org](<https://kaikki.org/dictionary/English/index.html>), a machine-readable Wiktionary fork. The data is made available under the same licenses as Wiktionary — both [CC-BY-SA](<https://en.wiktionary.org/wiki/Wiktionary:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License>) and [GFDL](<https://en.wiktionary.org/wiki/Wiktionary:Text_of_the_GNU_Free_Documentation_License>).`

func dictSlash(bot *discordgo.Session, i *discordgo.InteractionCreate) {
	var res string
	cmdata := i.ApplicationCommandData()
	switch cmdata.Name {
	case "copyright":
		res = notice
	case "dict":
		searchTerm := cmdata.Options[0].StringValue()
		res = dictMap.Query(searchTerm)
	default:
		res = "Unknown command"
	}

	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res,
		},
	}
	fmt.Println("responding to interaction")
	err := bot.InteractionRespond(i.Interaction, &resp)
	if err != nil {
		log.Println("error reacting to interaction", err)
	}
}

func init() {
	for _, v := range cmds {
		_, cmdFound := Commands[v.Name]
		_, handlerFound := Handlers[v.Name]
		if !cmdFound && !handlerFound {
			Commands[v.Name] = &v
			Handlers[v.Name] = dictSlash
		}
	}
	var err error
	dictMap, err = dict.ParseDict(rawDict, nil)
	if err != nil {
		log.Fatalf("Cannot parse dictionary: %v", err)
	}
}
