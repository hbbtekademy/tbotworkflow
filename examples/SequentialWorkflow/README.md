# Sequential Telegram Bot Workflow

## Step 1
User initiates a Telegram Bot Command. *e.g. /subscribe*

<table>
  <tr>
    <td> Telegram App </td> <td> Code Snippet </td>
  </tr>
  <tr>
    <td>
      
![Step 1](https://raw.githubusercontent.com/hbbtekademy/images-repo/main/tbotworkflow/examples/SequentialWorkflow/SeqStep1.jpg)
    </td>
    <td>
      <pre>
      
```go
// Step1 of the workflow. 
// Keyboard is nil so standard text keyboard will be displayed.
step1 := tbotworkflow.NewWorkflowStep("Step1", "Name", "Please enter your Name", nil)
```

</pre>
    </td>
  </tr>
</table>

## Step 2
Request user for some more text input

<table>
  <tr>
    <td> Telegram App </td> <td> Code Snippet </td>
  </tr>
  <tr>
    <td>
      
![Step 2](https://raw.githubusercontent.com/hbbtekademy/images-repo/main/tbotworkflow/examples/SequentialWorkflow/SeqStep2.jpg)
    </td>
    <td>
      <pre>
      
```go
// Step2 of the workflow. 
// Keyboard is nil so standard text keyboard will be displayed.
step2 := tbotworkflow.NewWorkflowStep("Step2", "Email", "Please enter your Email", nil)
```

</pre>
    </td>
  </tr>
</table>

## Step 3
Request user input with custom Keyboard

<table>
  <tr>
    <td> Telegram App </td> <td> Code Snippet </td>
  </tr>
  <tr>
    <td>
      
![Step 3](https://raw.githubusercontent.com/hbbtekademy/images-repo/main/tbotworkflow/examples/SequentialWorkflow/SeqStep3.jpg)
    </td>
    <td>
      <pre>
      
```go
// Step3 Keyboard.
step3KB := getStep3Keyboard()
// Step3 of the workflow with custom keyboard.
step3 := tbotworkflow.NewWorkflowStep("Step3", "Plan", 
  "Please select your Subscription Plan", step3KB)

// getStep3Keyboard returns keyboard to be displayed in step3 of the workflow
// This is usual Go Telegram Bot API code.
func getStep3Keyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Silver"),
			tgbotapi.NewKeyboardButton("Gold"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel"),
		),
	)
	return &kb
}
```

</pre>
    </td>
  </tr>
</table>

## Step 4
Request user confirmation using a custom keyboard and custom Reply Text Function

<table>
  <tr>
    <td> Telegram App </td> <td> Code Snippet </td>
  </tr>
  <tr>
    <td>
      
![Step 4](https://raw.githubusercontent.com/hbbtekademy/images-repo/main/tbotworkflow/examples/SequentialWorkflow/SeqStep4.jpg)
    </td>
    <td>
      <pre>
      
```go
// Step4 keyboard.
step4KB := getStep4Keyboard()
// Step4 with custom keyboard.
step4 := tbotworkflow.NewWorkflowStep("Step4", "Confirmation", "", step4KB)
// Custom reply text func
step4.ReplyTextFunc = func(ui *tbotworkflow.UserInputs) string {
	return fmt.Sprintf(`Please confirm your details:
	Name: %s
	Email: %s
	Subscription: %s`, ui.Data["Name"], ui.Data["Email"], ui.Data["Plan"])
}

// getStep4Keyboard returns keyboard to be displayed in step4 of the workflow
// This is usual Go Telegram Bot API code.
func getStep4Keyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Proceed"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel"),
		),
	)
	return &kb
}
```

</pre>
    </td>
  </tr>
</table>

## Step 5
Complete the workflow and return with the last step

<table>
  <tr>
    <td> Telegram App </td> <td> Code Snippet </td>
  </tr>
  <tr>
    <td>
      
![Step 5](https://raw.githubusercontent.com/hbbtekademy/images-repo/main/tbotworkflow/examples/SequentialWorkflow/SeqStep5.jpg)
    </td>
    <td>
      <pre>
      
```go
// Step5 (Last) of the workflow
step5 := tbotworkflow.NewWorkflowStep("Step5", "", "Thanks Proceeding with registration", nil)
```

</pre>
    </td>
  </tr>
</table>

## Chaining the steps together
```go
// Sequential Workflow. Chaining the steps.
step1.Next = &step2
step2.Next = &step3
step3.Next = &step4
step4.Next = &step5
```

## Creating the Workflow and Workflow Controller
```go
// Create a new workflow for command "/subscribe" with Step1 as at the root/starting step.
wf := tbotworkflow.NewWorkflow("WF1", "subscribe", &step1)

// Cancel button config for the entire workflow.
cancelBtnConfig := tbotworkflow.NewCancelButtonConfig("Cancel", "Canceling registeration.")
wf.CancelButtonConfig = cancelBtnConfig

// Create new Workflow Controller
wfc := tbotworkflow.NewWorkflowController("WFC")
// Add the Workflow to the controller. Any number of workflows can be added to a Workflow Controller.
wfc.AddWorkflow(&wf)
```

## Executing the Workflow and processing the user inputs
```go
// Process the Telegram Bot Updates
for update := range updates {
	// Do any pre-processing before invoking the workflows here.

	// Execute the workflow
	userInputs, done := wfc.Execute(update.Message, botAPI.Send)
	if !done {
		continue
	}

	log.Printf("UserID: %d, Command: %s, Data: %v", userInputs.UID, userInputs.Command, userInputs.Data)

	// Handle the user inputs as required.
}
```

<pre>
2022/06/09 20:28:09 UserID: 5104523246, Command: SUBSCRIBE, Data: map[Confirmation:Proceed Email:Hbb@hbb.com Name:HBB HBB Plan:Gold]
</pre>
