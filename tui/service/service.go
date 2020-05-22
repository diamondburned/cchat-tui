package service

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/log"
	"github.com/diamondburned/cchat-tui/tui/service/server"
	"github.com/diamondburned/cchat-tui/tui/ti"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Container struct {
	*tview.TreeView
	Root *tview.TreeNode
}

func NewContainer() *Container {
	root := tview.NewTreeNode(".")

	tree := tview.NewTreeView()
	tree.SetGraphics(true)
	tree.SetPrefixes([]string{"", "", "#", ">", ">", ">"})
	tree.SetRoot(root)
	tree.SetBackgroundColor(-1)
	tree.SetTopLevel(1)
	tree.SetDoneFunc(func(tcell.Key) {})

	root.AddChild(tview.NewTreeNode("test 1"))
	root.AddChild(tview.NewTreeNode("test 1"))
	root.AddChild(tview.NewTreeNode("test 1"))
	root.AddChild(tview.NewTreeNode("test 1"))
	root.AddChild(tview.NewTreeNode("test 1"))
	root.AddChild(tview.NewTreeNode("test 1"))

	return &Container{tree, root}
}

func (c *Container) AddService(sv cchat.Service, d ti.Drawer) {
	c.Root.AddChild(NewService(sv, d).TreeNode)
}

type Service struct {
	*tview.TreeNode
	Service cchat.Service
	drawer  ti.Drawer
}

func NewService(sv cchat.Service, d ti.Drawer) *Service {
	service := &Service{
		TreeNode: tview.NewTreeNode(sv.Name()),
		Service:  sv,
		drawer:   d,
	}
	service.TreeNode.SetReference(sv)

	if err := sv.Servers(service); err != nil {
		log.Error(errors.Wrap(err, "Failed to list service "+service.GetText()))
	}

	return service
}

func (sv *Service) SetServers(cservers []cchat.Server) {
	var children = make([]*tview.TreeNode, len(cservers))
	for i, cserver := range cservers {
		children[i] = server.NewNode(cserver, sv.drawer).TreeNode
	}

	sv.drawer.QueueUpdateDraw(func() {
		sv.TreeNode.SetChildren(children)
	})
}
