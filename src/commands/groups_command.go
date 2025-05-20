package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/usace/nsi-cli/api"
	"github.com/usace/nsi-cli/config"
)

type GroupsCommand struct {
	config            config.Config
	listGroupsFlags   *flag.FlagSet
	findGroupsFlags   *flag.FlagSet
	findGroupsId      *string
	addUserFlags      *flag.FlagSet
	addUserGroupId    *string
	addUserUserId     *string
	addUserRole       *string
	removeUserFlags   *flag.FlagSet
	removeUserGroupId *string
	removeUserUserId  *string
}

func NewGroupsCommand(config config.Config) GroupsCommand {
	listGroupsFlags := flag.NewFlagSet("list", flag.ExitOnError)
	findGroupsFlags := flag.NewFlagSet("find", flag.ExitOnError)
	addUserFlags := flag.NewFlagSet("add-user", flag.ExitOnError)
	removeUserFlags := flag.NewFlagSet("remove-user", flag.ExitOnError)

	return GroupsCommand{
		config:            config,
		listGroupsFlags:   listGroupsFlags,
		findGroupsFlags:   findGroupsFlags,
		findGroupsId:      findGroupsFlags.String("group-id", "", "The group ID"),
		addUserFlags:      addUserFlags,
		addUserGroupId:    addUserFlags.String("group-id", "", "The group ID"),
		addUserUserId:     addUserFlags.String("user-id", "", "The user ID"),
		addUserRole:       addUserFlags.String("role", "", "The user role"),
		removeUserFlags:   removeUserFlags,
		removeUserGroupId: removeUserFlags.String("group-id", "", "The group ID"),
		removeUserUserId:  removeUserFlags.String("user-id", "", "The user ID"),
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

	case "add-user":
		c.addUserFlags.Parse(args[1:])

	case "remove-user":
		c.removeUserFlags.Parse(args[1:])

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

	if c.addUserFlags.Parsed() {
		c.addUserToGroup()
	}

	if c.removeUserFlags.Parsed() {
		c.removeUserFromGroup()
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

func (c GroupsCommand) addUserToGroup() {
	userId := *c.removeUserUserId
	groupId := *c.removeUserGroupId
	role := *c.addUserRole

	data := struct {
		GroupId string `json:"group_id"`
		Role    string `json:"role"`
	}{
		GroupId: groupId,
		Role:    role,
	}
	payload, _ := json.Marshal(data)

	_, err := api.PostJson(c.config, payload, "users", userId, "groups")
	if err != nil {
		log.Fatalf("Error adding user to group: %v", err)
	}

	fmt.Printf("User added to group")
}

func (c GroupsCommand) removeUserFromGroup() {
	userId := *c.removeUserUserId
	groupId := *c.removeUserGroupId

	_, err := api.Delete(c.config, "users", userId, "groups", groupId)
	if err != nil {
		log.Fatalf("Error removing user from group: %v", err)
	}

	fmt.Printf("User removed from group")
}
