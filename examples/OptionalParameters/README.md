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
Function for custom validation of the user input.

Example
```go
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

## WorkflowNotFoundReplyTextFunc

## ValidateInputFunc
