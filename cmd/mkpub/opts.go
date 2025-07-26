package main

import (
	"fmt"
	"reflect"
	"strconv"
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
			fieldValue := reflectValue.FieldByName(field.Name)

			subcommand, isSubcommand := field.Tag.Lookup("subcommand")
			if isSubcommand {
				commandsRegistered++
			}

			if isSubcommand && command == "" && strings.EqualFold(subcommand, args[i]) {
				command = subcommand
				subcommandField := fieldValue.Addr().Interface()

				_, _, err := OptsParse(subcommandField, args[i:])
				if err != nil {
					return command, posArgs, err
				}

				isSet = true
				continue
			}

			long, longExists := field.Tag.Lookup("long")
			short, shortExists := field.Tag.Lookup("short")

			if !longExists && !shortExists { // not valid
				continue
			}

			if (longExists && strings.EqualFold(args[i], "--"+long)) || (shortExists && strings.EqualFold(args[i], "-"+short)) {
				if !fieldValue.CanSet() {
					return command, posArgs, fmt.Errorf("flag '%s': opts field '%s' cannot be given a value", args[i], field.Name)
				}

				switch fieldValue.Kind() {
				case reflect.Bool:
					fieldValue.SetBool(true)
				case reflect.String:
					if (i + 1) >= len(args) {
						return command, posArgs, fmt.Errorf("flag '%s': missing value argument of type 'string' for flag", args[i])
					}

					fieldValue.SetString(args[i+1])
					i++
				case reflect.Int:
					if (i + 1) >= len(args) {
						return command, posArgs, fmt.Errorf("flag '%s': missing value argument of type 'int' for flag", args[i])
					}

					intArg, err := strconv.Atoi(args[i+1])
					if err != nil {
						return command, posArgs, fmt.Errorf("flag '%s': %w", args[i], err)
					}

					fieldValue.SetInt(int64(intArg))
					i++
				default:
					return command, posArgs, fmt.Errorf("flag '%s': unsupported field type %v", args[i], fieldValue.Type())
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
