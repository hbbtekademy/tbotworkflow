package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hbbtekademy/tbotworkflow"
)

var (
	botToken string = "Put your Bot TOKEN here"
)

func getCondWorkflow() *tbotworkflow.TBotWorkflow {

	// Step 1
	acKB := getACListKeyboard()
	step1 := tbotworkflow.NewWorkflowStep("AC Name", "ACName", "Please select an AC to control", acKB)

	// Conditional Step 2
	acControlKB := getACControlKeyboard()
	step2 := tbotworkflow.NewWorkflowStep("AC Action", "ACAction", "Please select an option", acControlKB)
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

	// Step 3 - Quick Start Branch (Last step of the branch)
	step3QS := tbotworkflow.NewWorkflowStep("Quick Start", "", "Quick Starting AC", nil)
	step3QS.ReplyTextFunc = func(ui *tbotworkflow.UserInputs) string {
		acName := ui.Data["ACName"]
		return fmt.Sprintf("Starting AC <b>%s</b> with Temp 27 C at Medium fan speed", acName)
	}

	// Step 3 - Turn OFF Branch (Last step of the branch)
	step3Off := tbotworkflow.NewWorkflowStep("Turn OFF", "", "Turning AC OFF", nil)
	step3Off.ReplyTextFunc = func(ui *tbotworkflow.UserInputs) string {
		acName := ui.Data["ACName"]
		return fmt.Sprintf("Turning OFF AC <b>%s</b>", acName)
	}

	// Step 3 - Temperature Branch
	tempKB := getTemperatureKeyboard()
	step3Temp := tbotworkflow.NewWorkflowStep("AC Temperature", "ACTemp", "Please select AC Temperature", tempKB)

	// Step 4 - Fan Speed
	fsKB := getFanSpeedKeyboard()
	step4 := tbotworkflow.NewWorkflowStep("AC Fan Speed", "ACFanSpeed", "Please select the Fan Speed", fsKB)

	// Step 5 - Confirmation
	powerKB := getACPowerKeyboard()
	step5 := tbotworkflow.NewWorkflowStep("AC Power", "ACPower", "Please confirm AC can be turned ON", powerKB)

	// Step 6 - Last step
	step6 := tbotworkflow.NewWorkflowStep("ACLastStep", "", "Starting AC with following Parameters:", nil)
	step6.ReplyTextFunc = func(ui *tbotworkflow.UserInputs) string {
		acName := ui.Data["ACName"]
		temp := ui.Data["ACTemp"]
		fs := ui.Data["ACFanSpeed"]
		return fmt.Sprintf("Starting AC <b>%s</b> with Temp %s and Fan Speed %s", acName, temp, fs)
	}

	// Chain the steps
	step1.Next = &step2
	// Conditional chaining of steps
	step2.ConditionalNext["QuickStart"] = &step3QS
	step2.ConditionalNext["OFF"] = &step3Off
	step2.ConditionalNext["Temp"] = &step3Temp
	// Seq steps continue for step3Temp
	step3Temp.Next = &step4
	step4.Next = &step5
	step5.Next = &step6

	wf := tbotworkflow.NewWorkflow("WF1", "ac_control", &step1)

	// Register the RESET button as the "Cancel" button for the workflow
	cancelBtnConfig := tbotworkflow.NewCancelButtonConfig("RESET", "Clearing all input. Please start again")
	wf.CancelButtonConfig = cancelBtnConfig

	return &wf
}

func main() {
	// Create your telegram Bot API client
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed creating BotAPI. Error: %v", err)
	}

	// Subscribe to the Bot updates
	u := tgbotapi.NewUpdate(-1)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	// Create a TBotWorkflow
	wf := getCondWorkflow()

	// Create new Workflow Controller
	wfc := tbotworkflow.NewWorkflowController("WFC")
	// Add the Workflow to the controller. Any number of workflows can be added to a Workflow Controller.
	wfc.AddWorkflow(wf)

	// Process the Telegram Bot Updates
	for update := range updates {
		// Do any pre-processing before invoking the workflows here.

		// Execute the workflow
		userInputs, done := wfc.Execute(update.Message, botAPI.Send)
		if !done {
			continue
		}

		// Handle the user inputs as required.
		cmd := userInputs.Command
		uid := userInputs.UID
		if cmd == "AC_CONTROL" {
			action := userInputs.Data["ACAction"] // ACAction is the "Key" field in the step
			acName := userInputs.Data["ACName"]
			temp := userInputs.Data["ACTemp"]
			fan := userInputs.Data["ACFanSpeed"]

			log.Printf("UID: %d, Command: %s, Params[Action: %s, Name: %s, Temp: %s, Fan: %s]",
				uid, cmd, action, acName, temp, fan)
		}
	}
}

// getACListKeyboard returns the list of AirCons that can be controlled
func getACListKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Main Hall"),
			tgbotapi.NewKeyboardButton("Bedroom 1"),
			tgbotapi.NewKeyboardButton("Bedroom 2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	kb.Selective = true
	return &kb
}

// getACControlKeyboard returns the AirCon Controls
func getACControlKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Quick Start"),
			tgbotapi.NewKeyboardButton("Turn OFF"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Temperature"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	kb.Selective = true
	return &kb
}

// getTemperatureKeyboard returns the Temperature Control Keyboard
func getTemperatureKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("19 C"),
			tgbotapi.NewKeyboardButton("20 C"),
			tgbotapi.NewKeyboardButton("21 C"),
			tgbotapi.NewKeyboardButton("22 C"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("23 C"),
			tgbotapi.NewKeyboardButton("24 C"),
			tgbotapi.NewKeyboardButton("25 C"),
			tgbotapi.NewKeyboardButton("26 C"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	kb.Selective = true
	return &kb
}

// getFanSpeedKeyboard returns the Fan Speed Controk Keyboard
func getFanSpeedKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Min"),
			tgbotapi.NewKeyboardButton("Med"),
			tgbotapi.NewKeyboardButton("Max"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Auto"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	kb.Selective = true
	return &kb
}

func getACPowerKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Turn ON"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	kb.Selective = true
	return &kb
}
