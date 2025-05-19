package commands

import (
	"flag"
	"fmt"
	"log"

	"github.com/usace/nsi-cli/api"
	"github.com/usace/nsi-cli/config"
)

type GroupsCommand struct {
	config          config.Config
	listGroupsFlags *flag.FlagSet
	findGroupsFlags *flag.FlagSet
	findGroupsId    *string
}

func NewGroupsCommand(config config.Config) GroupsCommand {
	listGroupsFlags := flag.NewFlagSet("list", flag.ExitOnError)
	findGroupsFlags := flag.NewFlagSet("find", flag.ExitOnError)

	return GroupsCommand{
		config:          config,
		listGroupsFlags: listGroupsFlags,
		findGroupsFlags: findGroupsFlags,
		findGroupsId:    findGroupsFlags.String("group-id", "", "The group ID"),
	}
}

func (c GroupsCommand) Parse(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You must specify an action.")
	}

	action := args[0]
	switch action {
	case "list":
		c.listGroupsFlags.Parse(args[1:])

	case "find":
		c.findGroupsFlags.Parse(args[1:])

	default:
		return fmt.Errorf("Invalid groups action: %s", action)
	}

	return nil
}

func (c GroupsCommand) Validate() error {
	return nil
}

func (c GroupsCommand) Exec() {
	if c.listGroupsFlags.Parsed() {
		c.findAllGroups()
	}

	if c.findGroupsFlags.Parsed() {
		c.findGroup()
	}

	return
}

func (c GroupsCommand) findAllGroups() {
	groups, err := api.Get(c.config, "groups")
	if err != nil {
		log.Fatalf("Error fetching group records: %v", err)
	}

	fmt.Printf(api.ResponseStr(groups))
}

func (c GroupsCommand) findGroup() {
	group, err := api.Get(c.config, "groups", *c.findGroupsId)
	if err != nil {
		log.Fatalf("Error fetching group record: %v", err)
	}

	fmt.Printf(api.ResponseStr(group))
}
