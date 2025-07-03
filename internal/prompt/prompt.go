package prompt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// Input prompts for user input with validation
func Input(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   defaultValue,
		Templates: getInputTemplates(),
		Validate: func(input string) error {
			if input == "" && defaultValue == "" {
				return fmt.Errorf("Value cannot be empty")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	if result == "" {
		return defaultValue, nil
	}
	return result, nil
}

// Select prompts user to select from a list of options
func Select(label string, options []string) (int, string, error) {
	prompt := promptui.Select{
		Label:     label,
		Items:     options,
		Templates: getSelectTemplates(),
		Size:      10, // Show 10 items at a time
	}

	index, result, err := prompt.Run()
	if err != nil {
		return 0, "", fmt.Errorf("prompt failed: %v", err)
	}

	return index, result, nil
}

// InputWithValidation prompts for user input with custom validation
func InputWithValidation(label, defaultValue string, validate func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   defaultValue,
		Templates: getInputTemplates(),
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	if result == "" {
		return defaultValue, nil
	}
	return result, nil
}

// InputNumber prompts for numeric input with range validation
func InputNumber(label string, defaultValue, min, max int) (int, error) {
	validate := func(input string) error {
		if input == "" {
			return nil
		}
		num, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("Please enter a valid number")
		}
		if num < min || num > max {
			return fmt.Errorf("Please enter a number between %d and %d", min, max)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     label,
		Default:   strconv.Itoa(defaultValue),
		Templates: getInputTemplates(),
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("prompt failed: %v", err)
	}

	if result == "" {
		return defaultValue, nil
	}

	num, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %v", err)
	}

	return num, nil
}

// getInputTemplates returns custom templates for input prompts
func getInputTemplates() *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Prompt:          "{{ . }} ❯ ",
		Valid:           "{{ . }} ✔ ",
		Invalid:         "{{ . }} ✗ ",
		Success:         "{{ . | green | bold }} ✔ ",
		ValidationError: "{{ \"✗\" | red }} {{ . | red }}",
	}
}

// getSelectTemplates returns custom templates for select prompts
func getSelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . | bold }}?",
		Active:   "{{ \"›\" | cyan }} {{ . | cyan | bold }}",
		Inactive: "  {{ . | white }}",
		Selected: "{{ \"✔\" | green }} {{ . | green | bold }}",
		Details: `
{{ "──────────────────" | faint }}
{{ "Selected:" | faint }}  {{ . }}
{{ "──────────────────" | faint }}`,
	}
}

// Success prints a success message with enhanced formatting
func Success(message string) {
	prefix := color.New(color.FgBlack, color.BgGreen, color.Bold).Sprint(" SUCCESS ")
	content := color.New(color.FgGreen, color.Bold).Sprint(message)
	fmt.Printf("%s %s\n", prefix, content)
}

// Error prints an error message with enhanced formatting
func Error(message string) {
	prefix := color.New(color.FgWhite, color.BgRed, color.Bold).Sprint(" ERROR ")
	content := color.New(color.FgRed, color.Bold).Sprint(message)
	fmt.Printf("%s %s\n", prefix, content)
}

// Info prints an info message with enhanced formatting
func Info(message string) {
	prefix := color.New(color.FgBlack, color.BgCyan, color.Bold).Sprint(" INFO ")
	content := color.New(color.FgCyan).Sprint(message)
	fmt.Printf("%s %s\n", prefix, content)
}

// Warning prints a warning message with enhanced formatting
func Warning(message string) {
	prefix := color.New(color.FgBlack, color.BgYellow, color.Bold).Sprint(" WARNING ")
	content := color.New(color.FgYellow, color.Bold).Sprint(message)
	fmt.Printf("%s %s\n", prefix, content)
}

// StartSpinner starts a spinner with the given message
// Returns a spinner that should be stopped with spinner.Stop()
func StartSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Color("cyan")
	s.Start()
	return s
}
