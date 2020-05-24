package service

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/service/auth"
	"github.com/diamondburned/cchat-tui/tui/service/session"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Container struct {
	*tview.TreeView
	Root *tview.TreeNode
}

// TODO: deprecate customLists in favor of commands.

func NewContainer() *Container {
	root := tview.NewTreeNode(".")

	tree := tview.NewTreeView()
	tree.SetGraphics(true)
	tree.SetPrefixes([]string{"", "", "#", ">", ">", ">"})
	tree.SetRoot(root)
	tree.SetBackgroundColor(-1)
	tree.SetTopLevel(1)

	// Hack to set the tview foreground color back and forth on selection,
	// because tview sucks.
	var lastNode *tview.TreeNode
	tree.SetChangedFunc(func(node *tview.TreeNode) {
		// Undo the last node's color.
		if lastNode != nil {
			lastNode.SetColor(tcell.ColorWhite)
		}
		// Update the last node.
		lastNode = node
		lastNode.SetColor(tcell.ColorBlack)
	})

	return &Container{tree, root}
}

func (c *Container) FindService(fn func(cchat.Service) bool) cchat.Service {
	for _, child := range c.Root.GetChildren() {
		if sv := FromNode(child); sv != nil && fn(sv) {
			return sv
		}
	}
	return nil
}

// AddService is not thread-safe.
func (c *Container) AddService(sv cchat.Service) {
	c.Root.AddChild(NewService(sv).TreeNode)

	// Focus if this is the first child.
	if children := c.Root.GetChildren(); len(children) == 1 {
		c.SetCurrentNode(children[0])
	}
}

type Service struct {
	*tview.TreeNode
	Service cchat.Service
}

func FromNode(node *tview.TreeNode) cchat.Service {
	if sv, ok := node.GetReference().(cchat.Service); ok {
		return sv
	}
	return nil
}

func NewService(sv cchat.Service) *Service {
	service := &Service{
		TreeNode: tview.NewTreeNode(sv.Name()),
		Service:  sv,
	}
	service.SetReference(sv)
	service.SetExpanded(true)
	service.SetSelectedFunc(func() {
		auth.SpawnForm(sv, service.addSession)
	})

	return service
}

func (s *Service) addSession(ses cchat.Session) {
	node := session.NewNode(ses)
	s.TreeNode.AddChild(node.TreeNode)
}
