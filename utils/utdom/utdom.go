package utdom

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/koinworks/asgard-heimdal/libs/serror"
)

type DOMNode struct {
	Nodes []*html.Node
}

type DOMCondition struct {
	Handle func(*html.Node) bool
}

/**
 * Extensible function
 **/

func CIsTagName(n string) DOMCondition {
	return DOMCondition{
		Handle: func(e *html.Node) bool {
			return e.Type == html.ElementNode && e.Data == n
		},
	}
}

func CIsAttrExists(n string) DOMCondition {
	return DOMCondition{
		Handle: func(e *html.Node) bool {
			for _, v := range e.Attr {
				if v.Key == n {
					return true
				}
			}
			return false
		},
	}
}

func CIsAttrValue(v string) DOMCondition {
	return DOMCondition{
		Handle: func(e *html.Node) bool {
			for _, c := range e.Attr {
				if c.Val == v {
					return true
				}
			}
			return false
		},
	}
}

func CIsAttrMatch(n string, v string) DOMCondition {
	return DOMCondition{
		Handle: func(e *html.Node) bool {
			for _, c := range e.Attr {
				if c.Key == n && c.Val == v {
					return true
				}
			}
			return false
		},
	}
}

func CIsID(n string) DOMCondition {
	return CIsAttrMatch("id", n)
}

func CIsClassName(n string) DOMCondition {
	cls := strings.Split(n, ",")
	for k, v := range cls {
		cls[k] = strings.Trim(v, " ")
	}

	return DOMCondition{
		Handle: func(e *html.Node) bool {
			for _, c := range e.Attr {
				if c.Key == "class" {
					for _, v := range cls {
						if !strings.Contains(c.Val, v) {
							return false
						}
					}
					return true
				}
			}
			return false
		},
	}
}

/**
 * Public function
 **/

func Construct(v string) (*DOMNode, serror.SError) {
	node, err := html.Parse(strings.NewReader(v))
	if err != nil {
		return nil, serror.NewFromError(err)
	}

	return &DOMNode{
		Nodes: []*html.Node{node},
	}, nil
}

func (ox *DOMNode) IsExists() bool {
	return ox.Nodes != nil && len(ox.Nodes) > 0
}

func (ox *DOMNode) AdvancedFind(fns []DOMCondition) *DOMNode {
	if len(fns) > 0 {
		nodes := []*html.Node{}

		var f func(*html.Node)
		f = func(e *html.Node) {
			ok := true
			for _, v := range fns {
				if v.Handle != nil {
					if isOk := v.Handle(e); !isOk {
						ok = false
						break
					}
				}
			}

			if ok {
				nodes = append(nodes, e)
			}
			for c := e.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		for _, v := range ox.Nodes {
			f(v)
		}

		return &DOMNode{Nodes: nodes}
	}
	return &DOMNode{Nodes: ox.Nodes}
}

func (ox *DOMNode) AdvancedFilter(fns []DOMCondition) {
	if ox.IsExists() && fns != nil && len(fns) > 0 {
		newNodes := []*html.Node{}
		for _, e := range ox.Nodes {
			ok := true
			for _, v := range fns {
				if v.Handle != nil {
					if isOk := v.Handle(e); !isOk {
						ok = false
						break
					}
				}
			}

			if ok {
				newNodes = append(newNodes, e)
			}
		}
		ox.Nodes = newNodes
	}
}

func (ox *DOMNode) FindByTagName(n string) *DOMNode {
	return ox.AdvancedFind([]DOMCondition{CIsTagName(n)})
}

func (ox *DOMNode) FilterByTagName(n string) {
	ox.AdvancedFilter([]DOMCondition{CIsTagName(n)})
}

func (ox *DOMNode) FindByAttrName(n string) *DOMNode {
	return ox.AdvancedFind([]DOMCondition{CIsAttrExists(n)})
}

func (ox *DOMNode) FilterByAttrName(n string) {
	ox.AdvancedFilter([]DOMCondition{CIsAttrExists(n)})
}

func (ox *DOMNode) FindByAttrValue(v string) *DOMNode {
	return ox.AdvancedFind([]DOMCondition{CIsAttrValue(v)})
}

func (ox *DOMNode) FilterByAttrValue(v string) {
	ox.AdvancedFilter([]DOMCondition{CIsAttrValue(v)})
}

func (ox *DOMNode) FindByAttrMatch(n string, v string) *DOMNode {
	return ox.AdvancedFind([]DOMCondition{CIsAttrMatch(n, v)})
}

func (ox *DOMNode) FilterByAttrMatch(n string, v string) {
	ox.AdvancedFilter([]DOMCondition{CIsAttrMatch(n, v)})
}

func (ox *DOMNode) FindByID(n string) *DOMNode {
	return ox.AdvancedFind([]DOMCondition{CIsID(n)})
}

func (ox *DOMNode) FilterByID(n string) {
	ox.AdvancedFilter([]DOMCondition{CIsID(n)})
}

func (ox *DOMNode) FindByClassName(n string) *DOMNode {
	return ox.AdvancedFind([]DOMCondition{CIsClassName(n)})
}

func (ox *DOMNode) FilterByClassName(n string) {
	ox.AdvancedFilter([]DOMCondition{CIsClassName(n)})
}

func (ox *DOMNode) GetAttributes() (map[string]string, bool) {
	if ox.IsExists() {
		res := make(map[string]string)
		for _, v2 := range ox.Nodes[0].Attr {
			res[v2.Key] = v2.Val
		}
		return res, true
	}
	return nil, false
}

func (ox *DOMNode) GetAttribute(nm string) (string, bool) {
	if ox.IsExists() {
		for _, v := range ox.Nodes[0].Attr {
			if v.Key == nm {
				return v.Val, true
			}
		}
	}
	return "", false
}

func (ox *DOMNode) Each(c func(int, *DOMNode)) {
	if ox.IsExists() {
		for i, v := range ox.Nodes {
			if c != nil {
				c(i, &DOMNode{Nodes: []*html.Node{v}})
			}
		}
	}
}
