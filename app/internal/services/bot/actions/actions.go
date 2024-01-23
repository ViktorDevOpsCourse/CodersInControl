package actions

import (
	"fmt"
	"regexp"
)

type Action interface {
}

func CreateAction() Action {
	err := createAction("test action")
	if err != nil {
		return err
	}
	return nil
}

func createAction(action string) error {
	actionCommand := parseCommand(action)

	err := validateAction(actionCommand)
	if err != nil {
		return err
	}
	
	return nil
}

func validateAction(command ActionCommand) error {
	if command.Type == "" {
		return fmt.Errorf("unknown command")
	}

	if command.Branch == "" {
		return fmt.Errorf("servie tag not specified")
	}

	if command.Environment == "" {
		return fmt.Errorf("invalid environment")
	}

	return nil
}

type ActionCommand struct {
	Type        string
	Branch      string
	Environment string
}

var re = regexp.MustCompile(`<@(\S+)> (promote|list|diff|rollback) (\S+) to (stage|qa|prod)`)

func parseCommand(input string) ActionCommand {

	matches := re.FindStringSubmatch(input)
	return ActionCommand{
		Type:        matches[2],
		Branch:      matches[3],
		Environment: matches[4],
	}
}
