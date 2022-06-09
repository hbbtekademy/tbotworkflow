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

## Step 5

## Output after Step 5 (Workflow end)
