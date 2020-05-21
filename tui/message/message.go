package message

import "github.com/rivo/tview"

type MessageView struct {
	*tview.Flex
	Container *Container
	Input     *tview.InputField
}

type Container struct {
	*tview.Flex
}

type Message struct {
	*tview.Flex
}
