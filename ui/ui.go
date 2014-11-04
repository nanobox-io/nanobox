package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/mitchellh/colorstring"
)

// Prompt will prompt for input from the shell and return a trimmed response
func Prompt(p string) string {
	reader := bufio.NewReader(os.Stdin)

	//
	fmt.Print(p)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Unable to read input. See ~/.pagodabox/log.txt for details")
		Error("ui.Prompt", err)
	}

	return strings.TrimSpace(input)
}

// CPrint wraps a print message in 'colorstring' and passes it to fmt.Print
func CPrint(msg string) { fmt.Print(colorstring.Color(msg)) }

// CPrintln wraps a print message in 'colorstring' and passes it to fmt.Println
func CPrintln(msg string) { fmt.Println(colorstring.Color(msg)) }

// CPrompt wraps a message in 'colorstring' and passes it to prompt
func CPrompt(msg string) string { return Prompt(colorstring.Color(msg)) }

// Error will log a critical error to the .pagodabox/log.txt file and then call
// os.Exit(1). If no log.txt file can be found it will create one.
func Error(cmd string, msg error) {

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println("Unable to determine a home directory.", err)
		os.Exit(1)
	}

	logFile := homeDir + "/.pagodabox/log.txt"

	// check to see if they have a log file. we wont handle an error here because
	// there is a chance there wont be a log file. Instead we'll handle the file.
	// if we find a log file we can write errors to it
	l, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND, 0666)
	if l != nil {

		if _, err := l.WriteString("" +
			"\r\n================================== ERROR =======================================" +
			"\r\n" +
			"\r\nCommand : " + cmd +
			"\r\nTime    : " + time.Now().Format(time.RFC822) +
			"\r\nMessage :" +
			"\r\n" +
			"\r\n" + msg.Error() +
			"\r\n" +
			"\r\n================================================================================" +
			"\r\n"); err != nil {
			panic(err)
		}

		os.Exit(1)

		// if no log file is found we'll create it.
	} else {

		//
		l, err := os.Create(logFile)
		if err != nil {
			fmt.Println("Unable to create a .pagodabox/log.txt", err)
			os.Exit(1)
		}

		if _, err := l.WriteString("" +
			"Pagoda Box log.txt file - logging since " + time.Now().Format(time.RFC822) +
			"\r\n" +
			"\r\nAnytime the Pagoda Box CLI encounters an error it will dump the output here. If" +
			"\r\nyou are encountering an error with a command and you believe the error to be the" +
			"\r\nCLI's fault, you can find that error here and send it to us to take a look at." +
			"\r\n"); err != nil {
			panic(err)
		}

		// if this was run because someone deleted their log.txt, there is probably
		// an error that needs to be handled.
		if msg != nil {
			Error(cmd, msg)
		}

	}

	defer l.Close()
}
