package main

import (
	"fmt"
	"github.com/logrusorgru/aurora/v4"
)

// PrintMessage Prints a message
func PrintMessage(message string) {
	fmt.Printf("%v\n", styleBold(message))
}

// PrintQuestion Prints a question-styled message
func PrintQuestion(message string) {
	fmt.Printf("%v ", styleQuestion(message))
}

// PrintSuccessMessage Prints a successful-styled message
func PrintSuccessMessage(messages ...any) {
	finalMessage := ""

	for i := 0; i < len(messages); i++ {
		finalMessage += fmt.Sprintf("%v, ", messages[i])
	}

	finalMessage = finalMessage[:len(finalMessage)-2]

	fmt.Printf("%v \n", styleSuccess(finalMessage))
}

// PrintErrorMessage Prints an error-styled message
func PrintErrorMessage(errorLabel string, messages ...any) {
	finalErrorLabel := styleBold(styleError("[ERROR]"))
	finalMessage := ""

	if errorLabel != "" {
		finalErrorLabel += fmt.Sprintf(" %v ", styleError(errorLabel))
	}

	for i := 0; i < len(messages); i++ {
		finalMessage += fmt.Sprintf("%v, ", messages[i])
	}

	finalMessage = finalMessage[:len(finalMessage)-2]

	fmt.Printf("%v %v\n", finalErrorLabel, styleError(finalMessage))
}

// ***********************************************************************************
// ** Styles *************************************************************************
// ***********************************************************************************

// styleBold returns a bold-styled message
func styleBold(message any) string {
	return aurora.Bold(message).String()
}

func styleQuestion(message ...any) string {
	return aurora.BgBlack(message).String()
}

// styleError returns a error-styled message
func styleError(message any) string {
	return aurora.Red(message).String()
}

// styleSuccessful returns a successful-style message
func styleSuccess(message any) string {
	return aurora.Green(message).String()
}
