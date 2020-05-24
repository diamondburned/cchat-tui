package session

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat-tui/tui/service/server"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Node struct {
	*tview.TreeNode
	Session cchat.Session
}

func NewNode(ses cchat.Session) *Node {
	n, err := ses.Name()
	if err != nil {
		log.Error(errors.Wrap(err, "Error getting session name"))
	}

	node := &Node{
		tview.NewTreeNode(n),
		ses,
	}
	node.TreeNode.SetReference(node)

	return node
}

// SetServers is thread-safe.
func (n *Node) SetServers(cservers []cchat.Server) {
	var children = make([]*tview.TreeNode, len(cservers))
	for i, cserver := range cservers {
		children[i] = server.NewNode(cserver).TreeNode
	}

	app.QueueUpdateDraw(func() {
		n.TreeNode.SetChildren(children)
	})
}
