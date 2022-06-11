# Conditional Workflow
## Defining a conditional step
```go
// Conditional Step 2
acControlKB := getACControlKeyboard()
step2 := tbotworkflow.NewWorkflowStep("AC Action", "ACAction", "Please select an option", acControlKB)
// Define the conditional function based on the user input message for that step
step2.ConditionFunc = func(msg *tgbotapi.Message) string {
  if msg.Text == "Quick Start" {
		return "QuickStart"
	}
	if msg.Text == "Turn OFF" {
		return "OFF"
	}
	if msg.Text == "Temperature" {
		return "Temp"
	}
	return ""
}
```

## Chaining the conditional step
```go
// Conditional chaining of steps
// The keys of ConditionalNext should match the values returned by the steps ConditionFunc
step2.ConditionalNext["QuickStart"] = &step3QS
step2.ConditionalNext["OFF"] = &step3Off
step2.ConditionalNext["Temp"] = &step3Temp
```

