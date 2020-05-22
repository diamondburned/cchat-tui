package message

import (
	"github.com/rivo/tview"
)

type MessageView struct {
	*tview.Flex
	Container *Container
	Input     *tview.InputField
}

func NewMessageView() *MessageView {
	flex := tview.NewFlex()
	flex.SetBackgroundColor(-1)
	flex.SetDirection(tview.FlexRow)

	container := NewContainer()
	input := tview.NewInputField()
	// input.SetBackgroundColor(-1)
	input.SetFieldBackgroundColor(-1)
	input.SetPlaceholder("Message...")

	flex.AddItem(container, 0, 1, true)
	flex.AddItem(input, 1, 1, true)

	return &MessageView{
		flex,
		container,
		input,
	}
}

type Container struct {
	*tview.Flex
	Messages []Message
}

func NewContainer() *Container {
	flex := tview.NewFlex()
	flex.SetBackgroundColor(-1)
	flex.SetDirection(tview.FlexRow)

	return &Container{Flex: flex}
}

type Message struct {
	*tview.Flex
}
