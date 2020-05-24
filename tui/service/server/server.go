package server

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Node struct {
	*tview.TreeNode // use GetReference
	cchat.Server
}

func FromTreeNode(node *tview.TreeNode) *Node {
	// don't panic
	svnode, _ := node.GetReference().(*Node)
	return svnode
}

func NewNode(server cchat.Server) *Node {
	name, err := server.Name()
	if err != nil {
		log.Error(errors.Wrap(err, "Failed to make a server node"))
	}

	if name == "" {
		name = "no name"
	}

	var node = &Node{
		TreeNode: tview.NewTreeNode(name),
		Server:   server,
	}
	node.TreeNode.SetReference(node)

	list, ok := node.Server.(cchat.ServerList)
	if ok {
		if err := list.Servers(node); err != nil {
			log.Error(errors.Wrap(err, "Failed to list server "+node.GetText()))
		}
	}

	return node
}

func (node *Node) SetServers(servers []cchat.Server) {
	var children = make([]*tview.TreeNode, len(servers))
	for i, server := range servers {
		// We can reference TreeNode right away here, as we've already set a
		// reference in the NewNode constructor.
		children[i] = NewNode(server).TreeNode
	}

	app.QueueUpdateDraw(func() {
		node.TreeNode.SetChildren(children)
	})
}
