# TbotWorkflowStep - Workflow Step Optional Parameters
## ReplyTextFunc
Function for generating a custom formatted Text that should be displayed to the user for the given step.

Example
```go
step4 := tbotworkflow.NewWorkflowStep("Step4", "Confirmation", "", step4KB)
// Custom reply text func
step4.ReplyTextFunc = func(ui *tbotworkflow.UserInputs) string {
	return fmt.Sprintf(`Please confirm your details:
	Name: %s
	Email: %s
	Subscription: %s`, ui.Data["Name"], ui.Data["Email"], ui.Data["Plan"])
}
```

## ValidateInputFunc
Function for custom validation of the user input. Function should return a string and bool value.

If bool value is false, the contents of the string (error message) is displayed to the user and the same step will repeat again.

Example
```go
// Step2 of the workflow. Keyboard is nil so standard text keyboard will be displayed.
step2 := tbotworkflow.NewWorkflowStep("Step2", "Email", "Please enter your Email", nil)
step2.ValidateInputFunc = func(msg *tgbotapi.Message, kb *tgbotapi.ReplyKeyboardMarkup) (string, bool) {
	errMsg := ""
	ok := true

	_, err := mail.ParseAddress(msg.Text)
	if err != nil {
		errMsg = fmt.Sprintf("Invalid email: %s. Please enter a valid email address!", msg.Text)
		ok = false
	}
	return errMsg, ok
}
```

## CancelButtonConfig
CancelButtonConfig will tell the step if a particular user input should be considered as a workflow cancel request.

Example
```go
// Step 1
acKB := getACListKeyboard()
step1 := tbotworkflow.NewWorkflowStep("AC Name", "ACName", "Please select an AC to control", acKB)
// Register the RESET button as the "Cancel" button for the Step
cancelBtnConfig := tbotworkflow.NewCancelButtonConfig("RESET", "Clearing all input. Please start again")
step1.CancelButtonConfig = cancelBtnConfig
```

## ConditionFunc & ConditionalNext
Refer to the Conditional Workflow example.

[Contitional Workflow](https://github.com/hbbtekademy/tbotworkflow/tree/main/examples/ConditionalWorkflow)

# TBotWorkflow - Workflow Optional Paramters
## CancelButtonConfig
Use this parameter if all the steps in the workflow have the same Cancel Button.

Example
```go
wf := tbotworkflow.NewWorkflow("WF1", "ac_control", &step1)

// Register the RESET button as the "Cancel" button for the workflow
cancelBtnConfig := tbotworkflow.NewCancelButtonConfig("RESET", "Clearing all input. Please start again")
wf.CancelButtonConfig = cancelBtnConfig
```

# TBotWorkflowController - Workflow Controller Optional Parameters
## Logger
Go Standard Library logger. Logger is disabled by default. It can be enabled/disabled or completely overridden by user defined Std Lib logger

```go
wfc := tbotworkflow.NewWorkflowController("Workflows")
// Enable logging to stdout
wfc.EnableLogging(os.Stdout)

// Disable logging
wfc.DisableLogging()

// Override logger
wfc.SetLogger(log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags))
```

## WorkflowNotFoundReplyTextFunc
Error message to be displayed to the user if a workflow is not found for a given Telegram Bot Command.
```go
// Create new Workflow Controller
wfc := tbotworkflow.NewWorkflowController("WFC")
wfc.WorkflowNotFoundReplyTextFunc = func(msg *tgbotapi.Message) string {
	return fmt.Sprintf("Incorrect message: %s. Please select a valid Bot Command!", msg.Text)
}
```

## ValidateInputFunc
Define this function if the same Input Validation should be applied to all the steps of all the registered workflows.

