package main

import (
	"reflect"
	"strings"
)

type Program struct {
	Name          string
	UsageSynopsis string
	Description   string
}

func OptsParse(opts any, args []string) (string, []string, error) {
	var command string
	var posArgs []string

	commandsRegistered := 0

	reflectValue := reflect.ValueOf(opts).Elem()
	reflectType := reflect.TypeOf(opts).Elem()

	for i := 1; i < len(args); i++ {
		isSet := false

		if args[i] == "--" {
			break
		}

		// Ignored if no commands are registered in the end
		if strings.EqualFold(args[i], "help") && command == "" {
			command = "help"
			continue
		}

		// Ignored if no commands are registered in the end
		if strings.EqualFold(args[i], "version") && command == "" {
			command = "version"
			continue
		}

		for j := 0; j < reflectType.NumField(); j++ {
			field := reflectType.Field(j)
			fieldReflect := reflectValue.FieldByName(field.Name)

			subcommand, isSubcommand := field.Tag.Lookup("subcommand")
			if isSubcommand {
				commandsRegistered++
			}

			if isSubcommand && command == "" && strings.EqualFold(subcommand, args[i]) {
				command = subcommand
				subcommandField := fieldReflect.Addr().Interface()

				_, _, err := OptsParse(subcommandField, args[i:])
				if err != nil {
					return command, posArgs, err
				}

				isSet = true
			}
		}

		if !isSet {
			posArgs = append(posArgs, args[i])
		}
	}

	if commandsRegistered == 0 {
		command = ""
	}

	return command, posArgs, nil
}
