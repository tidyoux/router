package view

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

func addText(text string, m ...vecty.MarkupOrChild) *vecty.HTML {
	return vecty.Text(text, m...)
}

func addIcon(name string) *vecty.HTML {
	return elem.Italic(
		addClass("fas", fmt.Sprintf("fa-%s", name)),
	)
}

func addInput(markups []vecty.MarkupOrChild, h func(string)) *vecty.HTML {
	if h != nil {
		markups = append(markups, onInput(h))
	}
	return elem.Input(markups...)
}

func addSelect(selectValue string, options []string, h func(string)) *vecty.HTML {
	markups := make([]vecty.MarkupOrChild, 0, 2+len(options))
	markups = append(markups, addProprety("value", selectValue))
	for _, op := range options {
		markups = append(markups, elem.Option(
			addText(op),
		))
	}

	if h != nil {
		markups = append(markups, onSelectChange(func(value string) {
			h(value)
		}))
	}

	return elem.Div(
		addClass("select"),
		elem.Select(markups...),
	)
}

func addClass(class ...string) vecty.MarkupList {
	return vecty.Markup(
		vecty.Class(class...),
	)
}

func addProprety(key string, value interface{}) vecty.MarkupList {
	return vecty.Markup(
		vecty.Property(key, value),
	)
}

func addAttribute(key string, value interface{}) vecty.MarkupList {
	return vecty.Markup(
		vecty.Attribute(key, value),
	)
}

func onInput(h func(s string)) vecty.MarkupList {
	return vecty.Markup(
		event.Input(func(event *vecty.Event) {
			h(event.Target.Get("value").String())
		}),
	)
}

func onClick(h func()) vecty.MarkupList {
	return vecty.Markup(
		event.Click(func(event *vecty.Event) {
			h()
		}),
	)
}

func onCheckChange(h func(bool)) vecty.MarkupList {
	return vecty.Markup(
		event.Change(func(event *vecty.Event) {
			h(event.Target.Get("checked").Bool())
		}),
	)
}

func onSelectChange(h func(value string)) vecty.MarkupList {
	return vecty.Markup(
		event.Change(func(event *vecty.Event) {
			h(event.Target.Get("value").String())
		}),
	)
}

func addColumns(margin int, nodes []vecty.MarkupOrChild, weights []int) *vecty.HTML {
	items := make([]vecty.MarkupOrChild, 0, 3+len(nodes))
	items = append(items, addClass("columns"))
	if margin > 0 {
		items = append(items, addEmptyColumn(margin))
	}

	for i, node := range nodes {
		w := -1
		if i < len(weights) {
			w = weights[i]
		}
		items = append(items, addColumn(w, node))
	}

	if margin > 0 {
		items = append(items, addEmptyColumn(margin))
	}

	return elem.Div(items...)
}

func addCenter(margin int, nodes ...vecty.MarkupOrChild) *vecty.HTML {
	return addColumns(margin, nodes, nil)
}

func addColumn(width int, node vecty.MarkupOrChild) *vecty.HTML {
	cls := []string{"column"}
	if width > 0 {
		cls = append(cls, fmt.Sprintf("is-%d", width))
	}

	if node == nil {
		return elem.Div(
			addClass(cls...),
		)
	}

	return elem.Div(
		addClass(cls...),
		node,
	)
}

func addEmptyColumn(width int) *vecty.HTML {
	return addColumn(width, nil)
}
