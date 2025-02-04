/* For license and copyright information please see the LEGAL file in the code repository */

package cmd

import (
	"fmt"
	"os"

	errs "memar/operation/command/errors"
	command_p "memar/operation/command/protocol"
	"memar/protocol"
)

func ServeCLA(c protocol.Command, arguments []string) (err protocol.Error) {
	var serviceName string
	if len(arguments) > 0 {
		serviceName = arguments[0]
	} else {
		serviceName = "help"
	}

	// Also check for finding help command to check any custom help command
	var command protocol.Command = c.SubCommand(serviceName)
	if command == nil {
		// We don't find any related command even custom help, so print auto generated help.
		if serviceName == "help" || serviceName == "-h" || serviceName == "--help" {
			// TODO:::
			fmt.Fprintf(os.Stdout, "We must print help for You, But it is not implement yet. Sorry!\n")
			err = &errs.ErrServiceNotFound
			// helpMessage()
			// Accept 'go mod help' and 'go mod help foo' for 'go help mod' and 'go help mod foo'.
			// help.Help(os.Stdout, append(strings.Split(cfg.CmdName, " "), args[1:]...))
		} else {
			fmt.Fprintf(os.Stderr, "unknown command\nRun '%s help' for usage.\n", CommandPath(c))
			err = &errs.ErrServiceNotFound
		}
		return
	}

	var cmdName = command.Name()
	var cmdAbbr = command.Abbreviation()
	if serviceName != cmdName && serviceName != cmdAbbr {
		fmt.Fprintf(os.Stderr, "	Do you mean '%s %s'?\n", CommandPath(c), cmdName)
		err = &errs.ErrServiceCallByAlias
		return
	}

	err = command.ServeCLA(arguments[1:])
	return
}

// Root finds root command. or return nil if it is the root
func Root(c command_p.Command) (root command_p.Command) {
	for {
		if root.Parent() != nil {
			root = root.Parent()
		} else {
			break
		}
	}
	return
}

// CommandPath returns the full path to this command exclude itself.
func CommandPath(command command_p.Command) (fullName string) {
	for {
		fullName = command.Name() + " " + fullName
		command = command.Parent()
		if command == nil {
			break
		}
	}
	// remove trailing space
	fullName = fullName[:len(fullName)-1]
	return
}
