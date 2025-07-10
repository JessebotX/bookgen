package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Command struct {
	Name        string
	Description string

	// TODO: support flagsets
	Flags any
}

func flagParse(args []string, opts any) ([]string, error) {
	posArgs := make([]string, 0)

	optsStructValue := reflect.ValueOf(opts).Elem()
	optsStructType := reflect.TypeOf(opts).Elem()

	for i := 1; i < len(args); i++ {
		isSet := false

	outOpts:
		for j := 0; j < optsStructType.NumField(); j++ {
			field := optsStructType.Field(j)
			flag, ok := field.Tag.Lookup("long")
			if !ok {
				return []string{}, fmt.Errorf("field %v missing/empty tag name (required)", field.Name)
			}

			short, shortExists := field.Tag.Lookup("short")

			// Tries to match the following:
			//
			// --flag, -f, --flag=value, -f=value
			if strings.EqualFold(args[i], "--"+flag) || strings.EqualFold(args[i], "-"+short) {
				fieldValue := optsStructValue.FieldByName(field.Name)
				if !fieldValue.CanSet() {
					return []string{}, fmt.Errorf("arg '%v': field %v cannot be set", args[i], field.Name)
				}

				switch fieldValue.Kind() {
				case reflect.Bool:
					fieldValue.SetBool(true)
				case reflect.String:
					if (i + 1) >= len(args) {
						return []string{}, fmt.Errorf("arg '%v': missing value argument", args[i])
					}

					fieldValue.SetString(args[i+1])
				case reflect.Int:
					if (i + 1) >= len(args) {
						return []string{}, fmt.Errorf("arg '%v': missing value argument", args[i])
					}

					intArg, err := strconv.Atoi(args[i+1])
					if err != nil {
						return []string{}, fmt.Errorf("arg '%v': %w", args[i], err)
					}

					fieldValue.SetInt(int64(intArg))
					i++
				default:
					return []string{}, fmt.Errorf("arg '%v': unsupported field type %v", args[i], fieldValue.Type())
				}
				isSet = true
				break outOpts
			} else if strings.HasPrefix(args[i], "--"+flag+"=") || (shortExists && strings.HasPrefix(args[i], "-"+short+"=")) {
				argSplit := strings.SplitN(args[i], "=", 2)
				if len(argSplit) != 2 {
					return []string{}, fmt.Errorf("arg '%v': missing value after '='", args[i])
				}

				fieldValue := optsStructValue.FieldByName(field.Name)
				if !fieldValue.CanSet() {
					return []string{}, fmt.Errorf("arg '%v': field %v cannot be set", args[i], field.Name)
				}

				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(argSplit[1])
				case reflect.Int:
					intArg, err := strconv.Atoi(argSplit[1])
					if err != nil {
						return []string{}, fmt.Errorf("arg '%v': %w", args[i], err)
					}

					fieldValue.SetInt(int64(intArg))
				default:
					return []string{}, fmt.Errorf("arg '%v': unsupported field type %v", args[i], fieldValue.Type())
				}

				isSet = true
				break outOpts
			}
		}

		if !isSet {
			posArgs = append(posArgs, args[i])
		}
	}

	return posArgs, nil
}

func flagPrintHelp(opts any, commands []Command) {
	indentSize := 4
	indentLevel1 := 1
	indentLevel2 := 3

	fmt.Println("USAGE")

	for range indentSize * indentLevel1 {
		fmt.Printf(" ")
	}

	fmt.Println("bookgen [command] [flags...]")

	fmt.Println()

	fmt.Println("COMMANDS")

	for _, cmd := range commands {
		for range indentSize * indentLevel1 {
			fmt.Printf(" ")
		}

		fmt.Println(cmd.Name)

		for range indentSize * indentLevel2 {
			fmt.Printf(" ")
		}

		fmt.Println(cmd.Description)
	}

	fmt.Println()

	fmt.Println("FLAGS")

	optsStructValue := reflect.ValueOf(opts).Elem()
	optsStructType := reflect.TypeOf(opts).Elem()

	for i := 0; i < optsStructType.NumField(); i++ {
		field := optsStructType.Field(i)
		long, ok := field.Tag.Lookup("long")
		if !ok {
			continue
		}

		for range indentSize * indentLevel1 {
			fmt.Printf(" ")
		}

		short, ok := field.Tag.Lookup("short")
		if ok {
			fmt.Printf("-%v, ", short)
		} else {
			fmt.Printf("    ")
		}

		fieldKind := optsStructValue.FieldByName(field.Name).Kind()
		switch fieldKind {
		case reflect.Bool:
			fmt.Printf("--%v\n", long)
		default:
			fmt.Printf("--%v <value>\n", long)
		}

		desc, ok := field.Tag.Lookup("desc")
		if !ok {
			continue
		}

		for range indentSize * indentLevel2 {
			fmt.Printf(" ")
		}

		fmt.Printf("%v\n", desc)
	}
}
