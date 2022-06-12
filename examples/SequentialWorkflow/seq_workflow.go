package main

import (
	"fmt"
	"log"
	"net/mail"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tbotworkflow "github.com/hbbtekademy/tbotworkflow"
)

var (
	botToken string = "Put your Telegram Bot TOKEN Here"
)

// getSeqWorkflow will return a sequential TBotWorkflow
func getSeqWorkflow() *tbotworkflow.TBotWorkflow {
	// Define all the steps & keyboards required for the workflow

	// Step1 of the workflow. Keyboard is nil so standard text keyboard will be displayed.
	step1 := tbotworkflow.NewWorkflowStep("Step1", "Name", "Please enter your Name", nil)

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

	// Step3 Keyboard.
	step3KB := getStep3Keyboard()
	// Step3 of the workflow with custom keyboard.
	step3 := tbotworkflow.NewWorkflowStep("Step3", "Plan", "Please select your Subscription Plan", step3KB)

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

	// Step5 (Last) of the workflow
	step5 := tbotworkflow.NewWorkflowStep("Step5", "", "Thanks Proceeding with registration", nil)

	// Sequential Workflow. Chaining the steps.
	step1.Next = &step2
	step2.Next = &step3
	step3.Next = &step4
	step4.Next = &step5

	// Create a new workflow for command "/subscribe" with Step1 as at the root/starting step.
	wf := tbotworkflow.NewWorkflow("WF1", "subscribe", &step1)

	// Cancel button config for the entire workflow.
	cancelBtnConfig := tbotworkflow.NewCancelButtonConfig("Cancel", "Canceling registeration.")
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
	wf := getSeqWorkflow()

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

		log.Printf("UserID: %d, Command: %s, Data: %v", userInputs.UID, userInputs.Command, userInputs.Data)

		// Handle the user inputs as required.
	}
}

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
