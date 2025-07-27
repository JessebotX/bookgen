package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ProgramInfo struct {
	Name          string
	UsageSynopsis string
	Description   string
	Version       string
}

func OptsParse(opts any, args []string) (string, []string, error) {
	return optsParse(opts, args, false, true)
}

func optsParse(opts any, args []string, getArgsConsumed bool, parseEnv bool) (string, []string, error) {
	var command string
	var consumedArgs []string
	var posArgs []string

	reflectValue := reflect.ValueOf(opts).Elem()
	reflectType := reflect.TypeOf(opts).Elem()

	if parseEnv {
		if err := optsParseEnv(opts); err != nil {
			return command, posArgs, err
		}
	}

	for i := 1; i < len(args); i++ {
		isSet := false

		if args[i] == "--" {
			break
		}

		for j := 0; j < reflectType.NumField(); j++ {
			field := reflectType.Field(j)
			fieldValue := reflectValue.FieldByName(field.Name)

			subcommand, isSubcommand := field.Tag.Lookup("subcommand")

			if isSubcommand && command == "" && strings.EqualFold(subcommand, args[i]) {
				command = subcommand
				subcommandField := fieldValue.Addr().Interface()

				_, consumed, err := optsParse(subcommandField, args[i:], true, false)
				if err != nil {
					return command, posArgs, err
				}

				consumedArgs = consumed
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
					if getArgsConsumed {
						consumedArgs = append(consumedArgs, args[i])
					}
				case reflect.String:
					if (i + 1) >= len(args) {
						return command, posArgs, fmt.Errorf("flag '%s': missing value argument of type 'string' for flag", args[i])
					}

					fieldValue.SetString(args[i+1])
					if getArgsConsumed {
						consumedArgs = append(consumedArgs, args[i], args[i+1])
					}
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
					if getArgsConsumed {
						consumedArgs = append(consumedArgs, args[i], args[i+1])
					}
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

	if getArgsConsumed {
		return command, consumedArgs, nil
	} else {
		// essentially do posArgs - consumedArgs to get the actual
		// positional arguments after reading from subcommands
		var newPosArgs []string

		for _, posArg := range posArgs {
			found := false
			for _, consumedArg := range consumedArgs {
				if posArg == consumedArg {
					found = true
				}
			}

			if !found {
				newPosArgs = append(newPosArgs, posArg)
			}
		}

		return command, newPosArgs, nil
	}
}

func optsParseEnv(opts any) error {
	reflectValue := reflect.ValueOf(opts).Elem()
	reflectType := reflect.TypeOf(opts).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)

		envTag, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}

		envs := strings.Split(envTag, ",")

		for _, v := range envs {
			parsed := strings.SplitN(v, "==", 2)

			fieldValue := reflectValue.FieldByName(field.Name)
			if !fieldValue.CanSet() {
				return fmt.Errorf("field '%s' cannot be given a value", field.Name)
			}

			envValue, envValueExists := os.LookupEnv(parsed[0])
			if len(parsed) == 2 {
				if fieldValue.Kind() != reflect.Bool {
					return fmt.Errorf("field '%s' tag 'env' operator '==' can only be used for fields of type bool (read: '%s')", field.Name, envTag)
				}

				if envValueExists && envValue == parsed[1] {
					fieldValue.SetBool(true)
				} else {
					fieldValue.SetBool(false)
				}
			} else if len(parsed) == 1 {
				if !envValueExists {
					continue
				}

				switch fieldValue.Kind() {
				case reflect.Bool:
					if envValue == "" {
						fieldValue.SetBool(false)
					} else {
						fieldValue.SetBool(true)
					}
				case reflect.String:
					fieldValue.SetString(envValue)
				case reflect.Int:
					intArg, err := strconv.Atoi(parsed[1])
					if err != nil {
						return fmt.Errorf("field '%s': %w", field.Name, err)
					}

					fieldValue.SetInt(int64(intArg))
				default:
					return fmt.Errorf("field '%s' unsupported type '%v'", field.Name, fieldValue.Type())
				}
			}
		}
	}

	return nil
}

func OptsWriteHelp(w io.Writer, opts any, prog ProgramInfo) {
	_ = opts

	fmt.Fprintf(w, "USAGE\n")
	fmt.Fprintf(w, "    %s\n", prog.UsageSynopsis)
}
