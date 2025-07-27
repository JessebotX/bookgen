package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type ProgramInfo struct {
	Name          string
	UsageSynopsis string
	Description   string
	Version       string
}

type FlagHelpInfo struct {
	Short       string
	Long        string
	Description string
	Kind        reflect.Kind
}

type CommandHelpInfo struct {
	Name        string
	Description string
	Kind        reflect.Kind
}

type HelpExample struct {
	Usage       string
	Description string
}

const (
	DescriptionWrapWidth  = 62
	DescriptionWrapIndent = 10
)

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

			subcommand, isSubcommand := field.Tag.Lookup("command")

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

func OptsWriteHelp(w io.Writer, opts any, programName, cmdName string, examples ...HelpExample) error {
	return optsWriteHelp(w, opts, programName, cmdName, "", false, examples...)
}

func optsWriteHelp(w io.Writer, opts any, programName, cmdName, cmdDescription string, insideSubcommand bool, examples ...HelpExample) error {
	// ---
	// Read all commands and flags
	// ---
	var commands []CommandHelpInfo
	var flags []FlagHelpInfo

	reflectType := reflect.TypeOf(opts).Elem()
	reflectValue := reflect.ValueOf(opts).Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		fieldValue := reflectValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field '%s' cannot be given a value", field.Name)
		}

		// Commands
		subcommand, ok := field.Tag.Lookup("command")
		if ok {
			description, ok := field.Tag.Lookup("desc")
			if !ok {
				description = ""
			}

			if strings.EqualFold(subcommand, cmdName) {
				// ---
				// Print subcommand help
				// ---
				subcommandField := fieldValue.Addr().Interface()

				return optsWriteHelp(w, subcommandField, programName, cmdName, description, true)
			}

			optsCommand := CommandHelpInfo{
				Name:        subcommand,
				Description: strings.TrimSpace(wrapString(description, DescriptionWrapWidth, DescriptionWrapIndent)),
				Kind:        fieldValue.Kind(),
			}

			commands = append(commands, optsCommand)

			continue
		}

		// Flags
		long, longExists := field.Tag.Lookup("long")
		short, shortExists := field.Tag.Lookup("short")

		if !longExists && !shortExists {
			continue
		}

		description, ok := field.Tag.Lookup("desc")
		if !ok {
			description = ""
		}

		optsFlag := FlagHelpInfo{
			Long:        long,
			Short:       short,
			Description: strings.TrimSpace(wrapString(description, DescriptionWrapWidth, DescriptionWrapIndent)),
			Kind:        fieldValue.Kind(),
		}
		flags = append(flags, optsFlag)
	}

	// ---
	// Error for printing help for subcommand that does not exist
	// ---
	if !insideSubcommand && cmdName != "" {
		cmdExists := false
		for _, command := range commands {
			if strings.EqualFold(cmdName, command.Name) {
				cmdExists = true
			}

			if !cmdExists {
				return fmt.Errorf("no help/usage information for command '%s'", cmdName)
			}
		}
	}

	// ---
	// Begin printing
	// ---
	synopsis := programName
	if cmdName != "" {
		synopsis += " " + cmdName
	}
	if len(commands) > 0 {
		synopsis += " COMMAND"
	}
	if len(flags) > 0 {
		synopsis += " FLAGS..."
	}

	fmt.Fprintf(w, "USAGE\n")
	fmt.Fprintf(w, "    %s\n", synopsis)

	if strings.TrimSpace(cmdDescription) != "" {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "DESCRIPTION\n")
		fmt.Fprintf(w, "    %s\n", cmdDescription)
	}

	if len(examples) > 0 {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "EXAMPLES\n")

		for _, eg := range examples {
			fmt.Fprintf(w, "    $ %s\n", eg.Usage)
			fmt.Fprintf(w, "          %s\n", wrapString(eg.Description, DescriptionWrapWidth, DescriptionWrapIndent))
		}
	}

	if len(commands) > 0 {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "COMMANDS\n")
	}

	slices.SortFunc(commands, func(a, b CommandHelpInfo) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, command := range commands {
		fmt.Fprintf(w, "    %s\n", command.Name)
		if command.Description != "" {
			fmt.Fprintf(w, "          %s\n", wrapString(command.Description, DescriptionWrapWidth, DescriptionWrapIndent))
		}
	}

	if len(flags) > 0 {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "FLAGS\n")
	}

	slices.SortFunc(flags, func(a, b FlagHelpInfo) int {
		if a.Short != "" && b.Short != "" {
			if n := strings.Compare(a.Short, b.Short); n != 0 {
				return n
			}
		} else if a.Short != "" && b.Short == "" {
			if n := strings.Compare(a.Short, b.Long); n != 0 {
				return n
			}
		} else if a.Short == "" && b.Short != "" {
			if n := strings.Compare(a.Long, b.Short); n != 0 {
				return n
			}
		}

		return strings.Compare(a.Long, b.Long)
	})

	for _, f := range flags {
		if f.Long != "" && f.Short != "" {
			fmt.Fprintf(w, "    -%s, --%s ", f.Short, f.Long)
		} else if f.Long != "" {
			fmt.Fprintf(w, "    --%s ", f.Long)
		} else if f.Short != "" {
			fmt.Fprintf(w, "    -%s ", f.Short)
		}

		switch f.Kind {
		case reflect.String:
			fmt.Fprintf(w, "<TEXT-VALUE>")
		case reflect.Int:
			fmt.Fprintf(w, "<INTEGER-VALUE>")
		default:
			break
		}

		fmt.Fprintf(w, "\n")

		if f.Description != "" {
			fmt.Fprintf(w, "          %s\n", wrapString(f.Description, DescriptionWrapWidth, DescriptionWrapIndent))
		}
	}

	return nil
}

// Credit: <https://github.com/mitchellh/go-wordwrap>
// License: MIT
// <https://github.com/mitchellh/go-wordwrap/blob/master/LICENSE.md>
// (modified by adding indentation param used after newline)
//
// wrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func wrapString(s string, lim, indentation uint) string {
	nbsp := rune(0xA0)

	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	var wordBufLen, spaceBufLen uint

	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+spaceBufLen > lim {
					current = 0
				} else {
					current += spaceBufLen
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
				spaceBufLen = 0
			} else {
				current += spaceBufLen + wordBufLen
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) && char != nbsp {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBufLen + wordBufLen
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}

			spaceBuf.WriteRune(char)
			spaceBufLen++
		} else {
			wordBuf.WriteRune(char)
			wordBufLen++

			if current+wordBufLen+spaceBufLen > lim && wordBufLen < lim {
				buf.WriteByte(byte('\n'))
				for range indentation {
					buf.WriteByte(byte(' '))
				}

				current = 0
				spaceBuf.Reset()
				spaceBufLen = 0
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBufLen <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}
