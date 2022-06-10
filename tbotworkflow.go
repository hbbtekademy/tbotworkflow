package tbotworkflow

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// Default Text sent to the user for an incorrect input.
	defaultWFNotFoundReplyText string = "Message \"%s\" cannot be processed. Please select valid command."

	parseModeHTML     string = "HTML"
	parseModeMarkDown string = "MarkdownV2"
)

// CancelButtonConfig will tell the workflow if a particular user input
// should be considered as a workflow cancel request.
type CancelButtonConfig struct {
	cancelButtonExists bool
	cancelButtonText   string
	cancelButtonReply  string
}

// NewCancelButtonConfig returns a pointer to cancel button config.
// buttonText: Tells the workflow step what is the cancel buttons text. E.g. "Reset", "Cancel", "Restart" etc.
// buttonReply: Tells the workflow step what reply should be sent to the user when the cancel button is pressed.
// When user presses the Cancel button, the workflow will exit with the "buttonReply" Text.
func NewCancelButtonConfig(buttonText string, buttonReply string) *CancelButtonConfig {
	return &CancelButtonConfig{
		cancelButtonExists: true,
		cancelButtonText:   buttonText,
		cancelButtonReply:  buttonReply,
	}
}

// TBotWorkflowStep is the baisc unit of the workflow.
// A workflow consists of a number of TBotWorkflowStep's chained together.
type TBotWorkflowStep struct {
	// Name of the workflow step.
	Name string
	// Text that should be sent to the user at start of the step.
	ReplyText string
	// Key for the user input that will be available at the end of the workflow.
	Key string
	// Keyboard that should be presented to the user.
	// Set to nil to display the default text input keyboard.
	KB *tgbotapi.ReplyKeyboardMarkup
	// Next step to be executed after this step.
	// Set to nil for the last step.
	Next *TBotWorkflowStep
	// Function to be evaluated to determine the next step for conditional workflows.
	// If ConditionFunc is set, Step defined in "Next" will be ignored.
	ConditionFunc func(msg *tgbotapi.Message) string
	// A map of "ConditionFunc" outputs and the TBotWorkflowStep that should be executed for each of those outputs.
	ConditionalNext map[string]*TBotWorkflowStep
	// Function to generate the Text that should be sent to the user at start of the step.
	// If ReplyTextFunc is set, value defined in "ReplyText" is ignored.
	ReplyTextFunc func(ui *UserInputs) string
	// Function to validate the users input.
	// If the validation fails (function returns false), the string returned by this function is sent to the user.
	ValidateInputFunc func(msg *tgbotapi.Message, kb *tgbotapi.ReplyKeyboardMarkup) (string, bool)
	// Cancel button config for the step. Overrides the config set in the TBotWorkflow.
	CancelButtonConfig *CancelButtonConfig
}

// NewWorkflowStep returns a pointer to TBotWorkflowStep for given
// name, key, replyText and Telegram Reply Markup Keyboard.
// key is the "Key" in the UserInputs available at the end of the workflow.
// replyText is the text that should be sent to the user at start of the step.
func NewWorkflowStep(name string, key string, replyText string,
	kb *tgbotapi.ReplyKeyboardMarkup) TBotWorkflowStep {

	step := TBotWorkflowStep{
		Name:            name,
		ReplyText:       replyText,
		Key:             key,
		KB:              kb,
		Next:            nil,
		ConditionFunc:   nil,
		ConditionalNext: make(map[string]*TBotWorkflowStep),
	}

	return step
}

func (s *TBotWorkflowStep) isLastStep() bool {
	return s.Next == nil && s.ConditionFunc == nil
}

// UserInputs captures the user inputs for each step of the workflow
type UserInputs struct {
	// Telegram User ID
	UID int64
	// Telegram Command
	Command string
	// Data map to store the user inputs.
	// Map key is the "Key" defined in the TBotWorkflowStep
	// Map value is the Text entered by the user.
	Data map[string]string
}

// workflowTracker tracks at which step each user is in a given Workflow
type workflowTracker struct {
	UID                int64
	WorkflowName       string
	Command            string
	CurrentStep        *TBotWorkflowStep
	userInputs         UserInputs
	cancelButtonConfig *CancelButtonConfig
}

// userWfTracker keeps track of the workflow progress for all the users.
type userWfTracker struct {
	tracker map[int64]*workflowTracker
	m       sync.Mutex
}

func (u *userWfTracker) Add(uid int64, wfTracker *workflowTracker) {
	u.m.Lock()
	defer u.m.Unlock()
	u.tracker[uid] = wfTracker
}

func (u *userWfTracker) Delete(uid int64) {
	u.m.Lock()
	defer u.m.Unlock()
	delete(u.tracker, uid)
}

func (u *userWfTracker) Get(uid int64) (*workflowTracker, bool) {
	u.m.Lock()
	defer u.m.Unlock()
	wft, ok := u.tracker[uid]
	return wft, ok
}

// TBotWorkflow is the workflow to be triggered for a particular Bot Command
type TBotWorkflow struct {
	// Name of the workflow
	Name string
	// Command for which this workflow should be triggered
	Command string
	// The first step of the Workflow
	RootStep *TBotWorkflowStep
	// Cancel button config for this workflow. Can be overridden by the config set in the TBotWorkflowStep.
	CancelButtonConfig *CancelButtonConfig
}

// NewWorkflow returns a TBotWorkflow
func NewWorkflow(name string, command string, rootStep *TBotWorkflowStep) TBotWorkflow {
	wf := TBotWorkflow{
		Name:     name,
		Command:  strings.ToUpper(command),
		RootStep: rootStep,
	}

	return wf
}

// TBotWorkflowController controls all the workflows and their execution
type TBotWorkflowController struct {
	// Name of the Workflow Controller.
	Name string
	// A map of all the workflows controlled.
	workflows map[string]*TBotWorkflow
	// For tracking the workflow progress of all the users.
	userWFTracker userWfTracker
	// For logging useful information.
	// Logger is disabled by default but can be overriden/enabled/disabled
	// using the methods provided on the WorkflowController.
	Logger *log.Logger
	// Function to override the default Text sent to the users in case
	// this controller cannot handle the command sent by the user.
	WorkflowNotFoundReplyTextFunc func(msg *tgbotapi.Message) string
	// Global function to validate the user inputs.
	// ValidateInputFunc on TBotWorkflowStep takes priority over this function.
	ValidateInputFunc func(msg *tgbotapi.Message, kb *tgbotapi.ReplyKeyboardMarkup) (string, bool)
	// Telegram text parse mode. HTML or MarkdownV2.
	// Default value is HTML
	parseMode string
}

// NewWorkflowController returns a pointer to a new WorkflowController
func NewWorkflowController(name string) *TBotWorkflowController {
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	logger.SetOutput(ioutil.Discard)
	wfs := TBotWorkflowController{
		Name:          name,
		workflows:     make(map[string]*TBotWorkflow),
		userWFTracker: userWfTracker{tracker: make(map[int64]*workflowTracker)},
		Logger:        logger,
		parseMode:     parseModeHTML,
	}

	return &wfs
}

// SetLogger can be used to override the default logger
func (w *TBotWorkflowController) SetLogger(logger *log.Logger) {
	w.Logger = logger
}

// EnableLogging can be used to enable logging to the required io writer
func (w *TBotWorkflowController) EnableLogging(writer io.Writer) {
	w.Logger.SetOutput(writer)
}

// DisableLogging can be used to disable the logger
func (w *TBotWorkflowController) DisableLogging() {
	w.Logger.SetOutput(ioutil.Discard)
}

// SetMsgParseMode can be used to set the Telegram message parseMode (HTML or MarkdownV2).
// Default value is HTML
func (w *TBotWorkflowController) SetMsgParseMode(parseMode string) {
	w.parseMode = parseMode
}

// AddWorkflow is used to add a single workflow to the controller.
func (w *TBotWorkflowController) AddWorkflow(wf *TBotWorkflow) {
	w.workflows[wf.Command] = wf
}

// Execute runs one of the registered workflows given a Message from the user.
// At the end of the workflow, this method returns the user inputs captured at each step.
// bool = true means workflow has ended. UserInputs pointer will be "nil" till the workflow ends.
// This method takes the Message from the user and the Send function of the Telegram Bot API as inputs.
func (w *TBotWorkflowController) Execute(msg *tgbotapi.Message,
	sendFunc func(c tgbotapi.Chattable) (tgbotapi.Message, error)) (*UserInputs, bool) {
	if w.userWFTracker.tracker == nil {
		w.userWFTracker.tracker = make(map[int64]*workflowTracker)
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "")
	reply.ReplyToMessageID = msg.MessageID
	reply.ParseMode = w.parseMode

	userId := msg.From.ID
	userName := msg.From.UserName
	msgText := msg.Text
	w.Logger.Printf("Received message: %s from User: %d/%s", msgText, userId, userName)

	var wf *TBotWorkflow
	var found bool

	if msg.IsCommand() {
		cmd := strings.ToUpper(msg.Command())
		if wf, found = w.workflows[cmd]; !found {
			reply.Text = w.getWFNotFoundReplyText(msg, cmd)
			if _, err := sendFunc(reply); err != nil {
				w.Logger.Printf("Failed sending message. Error: %v", err)
			}
			return nil, false
		}
	}

	userWfTracker, found := w.userWFTracker.Get(userId)
	if msg.IsCommand() {
		w.Logger.Printf("User not found. Adding entry to tracker")

		cmd := strings.ToUpper(msg.Command())
		wfTracker := workflowTracker{
			UID:                userId,
			WorkflowName:       wf.Name,
			Command:            cmd,
			CurrentStep:        wf.RootStep,
			userInputs:         UserInputs{UID: userId, Command: cmd, Data: make(map[string]string)},
			cancelButtonConfig: wf.CancelButtonConfig,
		}
		w.userWFTracker.Add(userId, &wfTracker)
		userWfTracker = &wfTracker
		found = true
	}

	if !found {
		reply.Text = w.getWFNotFoundReplyText(msg, msg.Text)
		reply.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}

		if _, err := sendFunc(reply); err != nil {
			w.Logger.Printf("Failed sending message. Error: %v", err)
		}
		return nil, false
	}

	cancelBtnConfig := w.getCancelBtnConfig(userWfTracker)
	if cancelBtnConfig.cancelButtonExists && msgText == cancelBtnConfig.cancelButtonText {
		w.userWFTracker.Delete(userId)
		reply.Text = cancelBtnConfig.cancelButtonReply
		reply.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}

		if _, err := sendFunc(reply); err != nil {
			w.Logger.Printf("Failed sending message. Error: %v", err)
		}
		return nil, false
	}

	if !msg.IsCommand() {
		invalidReplyText, ok := w.validateInput(msg, userWfTracker)
		if ok {
			userWfTracker.userInputs.Data[userWfTracker.CurrentStep.Key] = msg.Text
		} else {
			reply.Text = invalidReplyText
			reply.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}

			if _, err := sendFunc(reply); err != nil {
				w.Logger.Printf("Failed sending message. Error: %v", err)
			}
		}

		if !userWfTracker.CurrentStep.isLastStep() && ok {
			if userWfTracker.CurrentStep.ConditionFunc != nil {
				nextStep := userWfTracker.CurrentStep.ConditionalNext[userWfTracker.CurrentStep.ConditionFunc(msg)]
				if nextStep == nil {
					reply.Text = fmt.Sprintf("Workflow %s broken. Cannot determine next step for CurrentStep: %s",
						userWfTracker.WorkflowName, userWfTracker.CurrentStep.Name)
					reply.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
					if _, err := sendFunc(reply); err != nil {
						w.Logger.Printf("Failed sending message. Error: %v", err)
					}
					w.userWFTracker.Delete(userId)
					return nil, false
				}
				w.Logger.Printf("Current Step: %s, Next Step: %s", userWfTracker.CurrentStep.Name, nextStep.Name)
				userWfTracker.CurrentStep = nextStep
			} else {
				if userWfTracker.CurrentStep.Next != nil {
					w.Logger.Printf("Current Step: %s, Next Step: %s", userWfTracker.CurrentStep.Name, userWfTracker.CurrentStep.Next.Name)
				}
				userWfTracker.CurrentStep = userWfTracker.CurrentStep.Next
			}
		}
	}

	reply.Text = userWfTracker.CurrentStep.ReplyText
	reply.ReplyMarkup = userWfTracker.CurrentStep.KB
	if userWfTracker.CurrentStep.KB == nil {
		reply.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
	}
	if userWfTracker.CurrentStep.ReplyTextFunc != nil {
		reply.Text = userWfTracker.CurrentStep.ReplyTextFunc(&userWfTracker.userInputs)
	}
	if userWfTracker.CurrentStep.isLastStep() {
		reply.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
	}

	if _, err := sendFunc(reply); err != nil {
		w.Logger.Printf("Failed sending message. Error: %v", err)
	}

	if userWfTracker.CurrentStep.isLastStep() {
		w.Logger.Println("WF ended. Return all the collected user inputs...")
		w.userWFTracker.Delete(userId)
		return &userWfTracker.userInputs, true
	}

	return nil, false
}

// defaultValidateInput is the default input validation method.
// This method will compare the user input with the Keyboard Button Text.
func (w *TBotWorkflowController) defaultValidateInput(msg *tgbotapi.Message, kb *tgbotapi.ReplyKeyboardMarkup) (string, bool) {
	if kb == nil {
		return "", true
	}

	buttons := kb.Keyboard
	validated := false
	replyText := ""

	if len(buttons) == 0 {
		validated = true
	}

	for i := 0; i < len(buttons); i++ {
		for j := 0; j < len(buttons[i]); j++ {
			if msg.Text == buttons[i][j].Text {
				validated = true
				break
			}
		}
	}

	if !validated {
		replyText = fmt.Sprintf("Invalid input %s. Please try again", msg.Text)
	}

	w.Logger.Printf("User input: %s validated: %v", msg.Text, validated)
	return replyText, validated
}

func (w *TBotWorkflowController) getWFNotFoundReplyText(msg *tgbotapi.Message, text string) string {
	replyText := ""
	if w.WorkflowNotFoundReplyTextFunc != nil {
		replyText = w.WorkflowNotFoundReplyTextFunc(msg)
	} else {
		replyText = fmt.Sprintf(defaultWFNotFoundReplyText, text)
	}
	return replyText
}

func (w *TBotWorkflowController) getCancelBtnConfig(userWfTracker *workflowTracker) *CancelButtonConfig {
	if userWfTracker.CurrentStep.CancelButtonConfig != nil {
		return userWfTracker.CurrentStep.CancelButtonConfig
	} else if userWfTracker.cancelButtonConfig != nil {
		return userWfTracker.cancelButtonConfig
	}
	return &CancelButtonConfig{
		cancelButtonExists: false,
		cancelButtonText:   "",
		cancelButtonReply:  "",
	}
}

func (w *TBotWorkflowController) validateInput(msg *tgbotapi.Message, userWfTracker *workflowTracker) (string, bool) {
	invalidReplyText := ""
	ok := true
	if userWfTracker.CurrentStep.ValidateInputFunc != nil {
		invalidReplyText, ok = userWfTracker.CurrentStep.ValidateInputFunc(msg, userWfTracker.CurrentStep.KB)
	} else if w.ValidateInputFunc != nil {
		invalidReplyText, ok = w.ValidateInputFunc(msg, userWfTracker.CurrentStep.KB)
	} else {
		invalidReplyText, ok = w.defaultValidateInput(msg, userWfTracker.CurrentStep.KB)
	}
	return invalidReplyText, ok
}
