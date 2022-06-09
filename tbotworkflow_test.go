package tbotworkflow

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mockSendFunc func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	sentMsgs     []tgbotapi.Message
)

type botInteraction struct {
	stepNo        int
	botMsg        tgbotapi.Message
	wfDone        bool
	expectedKey   string
	expectedValue string
}

func init() {
	mockSendFunc = mockTelegramSendFunc()
}

func TestWorkflows(t *testing.T) {
	sentMsgs = []tgbotapi.Message{}

	wfc := NewWorkflowController("WFC")
	seqWF := newSeqWorkflow("CMD1")
	condWF := newCondWorkflow("CMD2")

	wfc.AddWorkflow(&seqWF)
	wfc.AddWorkflow(&condWF)

	botInteractions := getSeqBotInteractions()

	var userInput *UserInputs
	var done bool
	for _, bi := range botInteractions {
		t.Run(fmt.Sprintf("Seq Bot Interaction Test Step:%d/Input:%s", bi.stepNo, bi.botMsg.Text), func(t *testing.T) {
			userInput, done = wfc.Execute(&bi.botMsg, mockSendFunc)
			if done != bi.wfDone {
				t.Errorf("Workflow completed at step %d but should not have", bi.stepNo)
			}
			if done == true && userInput == nil {
				t.Errorf("Workflow completed at step %d but UserInputs are nil", bi.stepNo)
			}
			if done == true && userInput.Command != "CMD1" {
				t.Errorf("Expected User Inputs for Command: CMD1 but got %s", userInput.Command)
			}
			if done == true && len(userInput.Data) != 2 {
				t.Errorf("Expected 2 user inputs but got %d", len(userInput.Data))
			}
		})
	}

	for _, bi := range botInteractions {
		t.Run(fmt.Sprintf("Seq Bot Interaction UserInput Validation Step %d Expected=%s:%s", bi.stepNo, bi.expectedKey, bi.expectedValue),
			func(t *testing.T) {
				if done == true && bi.wfDone != true && userInput.Data[bi.expectedKey] != bi.expectedValue {
					t.Errorf("Expected Key: %s, Value: %s but got got Value: %s",
						bi.expectedKey, bi.expectedValue, userInput.Data[bi.expectedKey])
				}
			},
		)
	}

	botInteractions = getCondBotInteractions()
	for _, bi := range botInteractions {
		t.Run(fmt.Sprintf("Cond Bot Interaction Test Step:%d/Input:%s", bi.stepNo, bi.botMsg.Text), func(t *testing.T) {
			userInput, done = wfc.Execute(&bi.botMsg, mockSendFunc)
			if done != bi.wfDone {
				t.Errorf("Workflow completed at step %d but should not have", bi.stepNo)
			}
			if done == true && userInput == nil {
				t.Errorf("Workflow completed at step %d but UserInputs are nil", bi.stepNo)
			}
			if done == true && userInput.Command != "CMD2" {
				t.Errorf("Expected User Inputs for Command: CMD2 but got %s", userInput.Command)
			}
			if done == true && len(userInput.Data) != 3 {
				t.Errorf("Expected 3 user inputs but got %d", len(userInput.Data))
			}
		})
	}

	for _, bi := range botInteractions {
		t.Run(fmt.Sprintf("Cond Bot Interaction UserInput Validation Step %d Expected=%s:%s", bi.stepNo, bi.expectedKey, bi.expectedValue),
			func(t *testing.T) {
				if done == true && bi.wfDone != true && userInput.Data[bi.expectedKey] != bi.expectedValue {
					t.Errorf("Expected Key: %s, Value: %s but got got Value: %s",
						bi.expectedKey, bi.expectedValue, userInput.Data[bi.expectedKey])
				}
			},
		)
	}

	sentMsgs = []tgbotapi.Message{}
}

func TestInvalidCommand(t *testing.T) {
	sentMsgs = []tgbotapi.Message{}
	wfc := NewWorkflowController("WFC")
	botInteractions := []*botInteraction{}

	botInteractions = append(botInteractions, &botInteraction{
		stepNo: 1,
		botMsg: mockBotCommand(1, "/XYZ"),
		wfDone: false,
	})

	var userInput *UserInputs
	var done bool
	userInput, done = wfc.Execute(&botInteractions[0].botMsg, mockSendFunc)
	if done == true || userInput != nil {
		t.Error("Workflow completed but should not have")
	}
	if len(sentMsgs) != 1 {
		t.Errorf("Expected only 1 message to be sent. But %d messages sent", len(sentMsgs))
	}
	expectedMsg := "Message \"XYZ\" cannot be processed. Please select valid command."
	if sentMsgs[0].Text != expectedMsg {
		t.Errorf("Expected \"%s\" message to be sent. But \"%s\" sent instead", expectedMsg, sentMsgs[0].Text)
	}

	wfc.WorkflowNotFoundReplyTextFunc = func(msg *tgbotapi.Message) string {
		return "Custom Error Message"
	}

	userInput, done = wfc.Execute(&botInteractions[0].botMsg, mockSendFunc)
	expectedMsg = "Custom Error Message"
	if sentMsgs[1].Text != expectedMsg {
		t.Errorf("Expected \"%s\" message to be sent. But \"%s\" sent instead", expectedMsg, sentMsgs[1].Text)
	}
	sentMsgs = []tgbotapi.Message{}
}

func TestBotCommand(t *testing.T) {
	botCmd := mockBotCommand(1, "/AC")
	if true != botCmd.IsCommand() {
		t.Error("Expected Bot Command true. But got false")
	}

	if "AC" != botCmd.Command() {
		t.Errorf("Expected Bot Command AC but got %s", botCmd.Command())
	}
}

func TestBotMessage(t *testing.T) {
	botMsg := mockBotMessage(1, "BotMessage")

	if true == botMsg.IsCommand() {
		t.Error("Expected Bot Command false. But got true")
	}

	if "BotMessage" != botMsg.Text {
		t.Errorf("Expected Bot Message BotMessage but got %s", botMsg.Text)
	}
}

func TestMockTelegramSendFunc(t *testing.T) {
	sentMsgs = []tgbotapi.Message{}
	msgConfig := tgbotapi.NewMessage(1, "TestMessage")
	msg, _ := mockSendFunc(msgConfig)

	if msg.Text != "TestMessage" {
		t.Errorf("Expected TestMessage but got %s", msg.Text)
	}
	if sentMsgs[0].Text != "TestMessage" {
		t.Errorf("Expected TestMessage but got %s", msg.Text)
	}
	sentMsgs = []tgbotapi.Message{}
}

func mockTelegramSendFunc() func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
		msgConfig := c.(tgbotapi.MessageConfig)
		msg := tgbotapi.Message{
			Text: msgConfig.Text,
		}

		sentMsgs = append(sentMsgs, msg)
		return msg, nil
	}
}

func mockBotCommand(chatID int64, text string) tgbotapi.Message {
	var entities []tgbotapi.MessageEntity
	entity := tgbotapi.MessageEntity{
		Type:   "bot_command",
		Length: len(text),
	}

	entities = append(entities, entity)
	msg := tgbotapi.Message{
		MessageID: 1,
		Text:      text,
		Entities:  entities,
		Chat: &tgbotapi.Chat{
			ID: chatID,
		},
		From: &tgbotapi.User{
			ID:       1234,
			UserName: "UnitTest",
		},
	}

	return msg
}

func mockBotMessage(chatID int64, text string) tgbotapi.Message {
	var entities []tgbotapi.MessageEntity
	entity := tgbotapi.MessageEntity{
		Type:   "",
		Length: len(text),
	}

	entities = append(entities, entity)
	msg := tgbotapi.Message{
		MessageID: 1,
		Text:      text,
		Entities:  entities,
		Chat: &tgbotapi.Chat{
			ID: chatID,
		},
		From: &tgbotapi.User{
			ID:       1234,
			UserName: "UnitTest",
		},
	}

	return msg
}

func getCondBotInteractions() []*botInteraction {
	botInteractions := []*botInteraction{}
	botInteractions = append(botInteractions, &botInteraction{
		stepNo:        1,
		botMsg:        mockBotCommand(1, "/CMD2"),
		wfDone:        false,
		expectedKey:   "K1",
		expectedValue: "Step1Option2",
	})
	botInteractions = append(botInteractions, &botInteraction{
		stepNo:        2,
		botMsg:        mockBotMessage(1, "Step1Option2"),
		wfDone:        false,
		expectedKey:   "CondK2",
		expectedValue: "Step2Condition1",
	})
	botInteractions = append(botInteractions, &botInteraction{
		stepNo:        3,
		botMsg:        mockBotMessage(1, "Step2Condition1"),
		wfDone:        false,
		expectedKey:   "C1K3",
		expectedValue: "C1Step3Option1",
	})
	botInteractions = append(botInteractions, &botInteraction{
		stepNo: 4,
		botMsg: mockBotMessage(1, "C1Step3Option1"),
		wfDone: true,
	})

	return botInteractions
}

func getSeqBotInteractions() []*botInteraction {
	botInteractions := []*botInteraction{}
	botInteractions = append(botInteractions, &botInteraction{
		stepNo:        1,
		botMsg:        mockBotCommand(1, "/CMD1"),
		wfDone:        false,
		expectedKey:   "K1",
		expectedValue: "Step1Option1",
	})
	botInteractions = append(botInteractions, &botInteraction{
		stepNo:        2,
		botMsg:        mockBotMessage(1, "Step1Option1"),
		wfDone:        false,
		expectedKey:   "K2",
		expectedValue: "Step2Option3",
	})
	botInteractions = append(botInteractions, &botInteraction{
		stepNo: 3,
		botMsg: mockBotMessage(1, "Step2Option3"),
		wfDone: true,
	})

	return botInteractions
}

func newSeqWorkflow(cmd string) TBotWorkflow {

	resetButton := NewCancelButtonConfig("RESET", "Clearing input. Please start again")
	step1KB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step1Option1"),
			tgbotapi.NewKeyboardButton("Step1Option2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step1Option3"),
			tgbotapi.NewKeyboardButton("Step1Option4"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	step1 := NewWorkflowStep("Step1", "K1", "Please select an option", &step1KB)

	step2KB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step2Option1"),
			tgbotapi.NewKeyboardButton("Step2Option2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step2Option3"),
			tgbotapi.NewKeyboardButton("Step2Option4"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	step2 := NewWorkflowStep("Step2", "K2", "Please select another option", &step2KB)
	step3 := NewWorkflowStep("Step3", "", "Please verify the selected Options", nil)

	step1.Next = &step2
	step2.Next = &step3
	step3.Next = nil

	wf := NewWorkflow("WF", cmd, &step1)
	wf.CancelButtonConfig = resetButton
	return wf
}

func newCondWorkflow(cmd string) TBotWorkflow {
	resetButton := NewCancelButtonConfig("RESET", "Clearing input. Please start again")
	step1KB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step1Option1"),
			tgbotapi.NewKeyboardButton("Step1Option2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step1Option3"),
			tgbotapi.NewKeyboardButton("Step1Option4"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	step1 := NewWorkflowStep("Step1", "K1", "Please select an option", &step1KB)

	step2KB := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Step2Condition1"),
			tgbotapi.NewKeyboardButton("Step2Condition2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
	condStep2 := NewWorkflowStep("CondStep2", "CondK2", "Please select a condition", &step2KB)
	condStep2.ConditionFunc = func(msg *tgbotapi.Message) string {
		switch msg.Text {
		case "Step2Condition1":
			return "C1"
		case "Step2Condition2":
			return "C2"
		default:
			return ""
		}
	}
	c1Step3KB := getSingleButtonKeyboard("C1Step3Option1")
	c2Step3KB := getSingleButtonKeyboard("C2Step3Option1")
	c1Step3 := NewWorkflowStep("C1Step3", "C1K3", "Please select an option", &c1Step3KB)
	c2Step3 := NewWorkflowStep("C2Step3", "C2K3", "Please select an option", &c2Step3KB)

	step4 := NewWorkflowStep("Step4", "", "Please verify the selected Options", nil)

	step1.Next = &condStep2
	condStep2.ConditionalNext["C1"] = &c1Step3
	condStep2.ConditionalNext["C2"] = &c2Step3
	c1Step3.Next = &step4
	c2Step3.Next = &step4
	step4.Next = nil

	wf := NewWorkflow("ConditionalWF", cmd, &step1)
	wf.CancelButtonConfig = resetButton
	return wf
}

func getSingleButtonKeyboard(btnTxt string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(btnTxt),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("RESET"),
		),
	)
}
