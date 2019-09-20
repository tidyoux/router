package view

import (
	"github.com/tidyoux/router/server/web/control"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type View interface {
	Init()
	Reset()

	Render() vecty.ComponentOrHTML
}

func rerender() {
	control.DispatchEvent(control.NewEvent(control.ERerender))
}

type Base struct {
	note string
}

func (view *Base) Init() {}

func (view *Base) Reset() {
	view.SetNode("")
}

func (view *Base) SetNode(note string) {
	view.note = note
}

func (view *Base) renderNote(margin int) *vecty.HTML {
	if len(view.note) == 0 {
		return elem.Div()
	}

	return addCenter(margin, elem.Div(
		addClass("notification", "is-info"),
		elem.Button(
			addClass("delete"),
			onClick(func() {
				view.Reset()
				rerender()
			}),
		),
		addText(view.note),
	))
}

type Modal struct {
	active bool
}

func NewModal() *Modal {
	return &Modal{}
}

func (view *Modal) Init() {}

func (view *Modal) Reset() {
	view.active = false
}

func (view *Modal) Active() {
	view.active = true
}

func (view *Modal) Render(title string, body vecty.MarkupOrChild, onClose func()) vecty.ComponentOrHTML {
	modalCls := []string{"modal"}
	if view.active {
		modalCls = append(modalCls, "is-active")
	}

	return elem.Div(
		addClass(modalCls...),

		elem.Div(
			addClass("modal-background"),
			onClick(func() {
				if onClose != nil {
					onClose()
				}
				view.Reset()
				rerender()
			}),
		),

		elem.Div(
			addClass("modal-card"),

			elem.Header(
				addClass("modal-card-head"),

				elem.Paragraph(
					addClass("modal-card-title"),
					addText(title),
				),

				elem.Anchor(
					addClass("delete"),
					onClick(func() {
						if onClose != nil {
							onClose()
						}

						view.Reset()
						rerender()
					}),
				),
			),

			elem.Section(
				addClass("modal-card-body"),

				body,
			),
		),
	)
}

type Component struct {
	vecty.Core

	render func() vecty.ComponentOrHTML
}

func NewComponent(render func() vecty.ComponentOrHTML) *Component {
	return &Component{
		render: render,
	}
}

func (c *Component) Render() vecty.ComponentOrHTML {
	return c.render()
}
