package c2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var helpMenu string = `
CSC C2 v0.0.1 (2024-03-01)
-------------------------------------------
Main menu commands:
Command                  Description
-------------------------------------------
help                     Show the menu
clear                    Clear the console
exit                     Exit C2
-------------------------------------------

`

func StartCLI() {
	autoCompleter := readline.NewPrefixCompleter(
		readline.PcItem("clear"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "C2 > ",
		AutoComplete:    autoCompleter,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err) // Fatal Case
	}

	defer rl.Close() // Defer closing the readline unit we exit func

	historyFile, err := os.Open(".history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for {
		command, err := rl.Readlne()
		if err != nil {
			log.Fatal(err)
		}

		if _, err = historyFile.WriteString(command + "\n"); err != nil {
			fmt.Println(err)
		}

		command = strings.TrimSpace(command) // Format our string command

		switch (command) {
		case "exit": // Exit cmd
			{
				fmt.Println("[C2] Exit Command, Shutting down")
				os.Exit(0)
			}
		case "help": // Help cmd
			{
				fmt.Print(helpMenu)
			}
		case "clear": // Clear cmd
			{
				if err := clearConsole(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func clearConsole() error {
	var cmd *exec.cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	// Redirect stdout
	cmd.Stdout = os.Stdout
	return cmd.Run()
}