package main

import (
	"bytes"
	"fmt"
	"os"
	"prcrastinate/internal/command"
	"prcrastinate/internal/platform"
)

// TODO: overall, consider printing help output to STDIN instead of STDERR so that things like
// less and other pagers can handle it reasonably

func main() {
	if len(os.Args) < 2 {
		failMsg :=
			"Must specify a valid subcommand.\n" +
			"Run \"prcr help\" to see valid commands and usage information."
		platform.FailOut(failMsg)
	}

	switch os.Args[1] {
		case "-h", "--help":
			printTopLevelHelp()
		case "help":
			if len(os.Args) == 2 {
				fmt.Fprintln(os.Stderr, printTopLevelHelp())
				os.Exit(1)
			}

			cmd, cmdFound := commands()[os.Args[2]]
			if !cmdFound {
				failMsg :=
					"Expected a valid Prcrastinate command, but was given \"%s\"\n" +
					"Run \"prcr help\" to see valid objects and usage information."
				platform.FailOut(fmt.Sprintf(failMsg, os.Args[1]))
			}

			fmt.Fprintln(os.Stderr, command.FmtHelp(cmd))
			os.Exit(1)
		default: {
			cmd, cmdFound := commands()[os.Args[1]]
			if !cmdFound {
				failMsg :=
					"Expected a valid Prcrastinate command, but was given \"%s\"\n" +
					"Run \"prcr help\" to see valid objects and usage information."
				platform.FailOut(fmt.Sprintf(failMsg, os.Args[1]))
			}

			cmd.Run(os.Args[2:])
		}
	}
}

func commands() map[string]command.CmdData {
	return command.MapNames([]command.CmdData{
		new(command.GenerateConfig),
		new(command.Pull),
	})
}

func printTopLevelHelp() string {
	var outputBuf bytes.Buffer

	topHelpHeader :=
		"Prcrastinate [prcr] is a CLI to aid in Github-centric PR review activities.\n" +
		"Prcrastinate uses subcommands\n" +
		"\n" +
		"You can see this message by running \"prcr help\".\n" +
		"\n" +
		"SUBCOMMANDS IN PRCRASTINATE:\n"
	fmt.Fprint(&outputBuf, topHelpHeader)

	commands := commands()

	nameSize := 0
	for _, cmd := range commands {
		curSize := len([]rune(cmd.Name()))
		if curSize > nameSize {
			nameSize = curSize
		}
	}
	nameSize++
	fmtString := fmt.Sprintf("  %%-%ds %%s\n", nameSize)

	for _, cmd := range commands {
		fmt.Fprintf(&outputBuf, fmtString, cmd.Name(), cmd.ShortDescr())
	}

	fmt.Fprint(&outputBuf, "\nSee help about any subcommand by running \"prcr help [command]\".")
	fmt.Fprint(&outputBuf, "\n\nCOMMON OPTIONS:\n")
	flags := command.InitBaseFlags(new(command.BaseArgs))
	fmt.Fprint(&outputBuf, flags.FlagUsages())

	return outputBuf.String()
}
