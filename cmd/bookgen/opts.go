package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func OptsParse(opts any, args []string) (string, []string, error) {
	return OptsParseHelper(opts, args, false)
}

func OptsParseHelper(opts any, args []string, returnConsumedArgs bool) (string, []string, error) {
	var posArgs []string
	var consumedArgs []string
	var command string

	reflectValue := reflect.ValueOf(opts).Elem()
	reflectType := reflect.TypeOf(opts).Elem()

	for i := 1; i < len(args); i++ {
		currentArg := args[i]
		isSet := false

		for j := 0; j < reflectType.NumField(); j++ {
			field := reflectType.Field(j)

			subcommand, ok := field.Tag.Lookup("subcommand")
			if ok && currentArg == subcommand {
				command = subcommand
				subcommandField := reflectValue.FieldByName(field.Name).Addr().Interface()

				_, consumed, err := OptsParseHelper(subcommandField, args[i:], true)
				if err != nil {
					return command, posArgs, err
				}

				consumedArgs = consumed
				isSet = true
				continue
			}

			long, ok := field.Tag.Lookup("long")
			if !ok {
				continue
			}

			short, shortExists := field.Tag.Lookup("short")

			if strings.EqualFold(currentArg, "--"+long) || (shortExists && strings.EqualFold(currentArg, "-"+short)) {
				fieldValue := reflectValue.FieldByName(field.Name)
				if !fieldValue.CanSet() {
					return command, posArgs, fmt.Errorf("arg '%v': field %v cannot be set", currentArg, field.Name)
				}

				switch fieldValue.Kind() {
				case reflect.Bool:
					fieldValue.SetBool(true)
					if returnConsumedArgs {
						consumedArgs = append(consumedArgs, currentArg)
					}
				case reflect.String:
					if (i + 1) >= len(args) {
						return command, posArgs, fmt.Errorf("arg '%v': missing value argument", currentArg)
					}

					fieldValue.SetString(args[i+1])
					if returnConsumedArgs {
						consumedArgs = append(consumedArgs, currentArg, args[i+1])
					}
					i++
				case reflect.Int:
					if (i + 1) >= len(args) {
						return command, posArgs, fmt.Errorf("arg '%v': missing value argument", currentArg)
					}

					intArg, err := strconv.Atoi(args[i+1])
					if err != nil {
						return command, posArgs, fmt.Errorf("arg '%v': %w", currentArg, err)
					}

					fieldValue.SetInt(int64(intArg))
					if returnConsumedArgs {
						consumedArgs = append(consumedArgs, currentArg, args[i+1])
					}
					i++
				default:
					return command, posArgs, fmt.Errorf("arg '%v': unsupported field type %v", currentArg, fieldValue.Type())
				}

				isSet = true
			}
		}

		if !isSet {
			posArgs = append(posArgs, currentArg)
		}
	}

	if returnConsumedArgs {
		return command, consumedArgs, nil
	} else {
		// essentially do posArgs - consumedArgs to get the actual positional arguments after reading from subcommands
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

// func optsPrintHelp(opts any, commands []CommandInfo) {
// 	indentSize := 4
// 	indentLevel1 := 1
// 	indentLevel2 := 3

// 	fmt.Println("USAGE")

// 	for range indentSize * indentLevel1 {
// 		fmt.Printf(" ")
// 	}

// 	fmt.Println("bookgen [command] [flags...]")

// 	fmt.Println()

// 	fmt.Println("COMMANDS")

// 	for _, cmd := range commands {
// 		for range indentSize * indentLevel1 {
// 			fmt.Printf(" ")
// 		}

// 		fmt.Println(cmd.Name)

// 		for range indentSize * indentLevel2 {
// 			fmt.Printf(" ")
// 		}

// 		fmt.Println(cmd.Description)
// 	}

// 	fmt.Println()

// 	fmt.Println("FLAGS")

// 	optsStructValue := reflect.ValueOf(opts).Elem()
// 	optsStructType := reflect.TypeOf(opts).Elem()

// 	for i := 0; i < optsStructType.NumField(); i++ {
// 		field := optsStructType.Field(i)
// 		long, ok := field.Tag.Lookup("long")
// 		if !ok {
// 			continue
// 		}

// 		for range indentSize * indentLevel1 {
// 			fmt.Printf(" ")
// 		}

// 		short, ok := field.Tag.Lookup("short")
// 		if ok {
// 			fmt.Printf("-%v, ", short)
// 		} else {
// 			fmt.Printf("    ")
// 		}

// 		fieldKind := optsStructValue.FieldByName(field.Name).Kind()
// 		switch fieldKind {
// 		case reflect.Bool:
// 			fmt.Printf("--%v\n", long)
// 		default:
// 			fmt.Printf("--%v <value>\n", long)
// 		}

// 		desc, ok := field.Tag.Lookup("desc")
// 		if !ok {
// 			continue
// 		}

// 		for range indentSize * indentLevel2 {
// 			fmt.Printf(" ")
// 		}

// 		fmt.Printf("%v\n", desc)
// 	}
// }
