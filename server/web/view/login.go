package view

import (
	"github.com/tidyoux/router/server/web/control"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Login struct {
	Base

	username string
	password string
}

func NewLogin() *Login {
	return &Login{}
}

func (view *Login) Reset() {
	view.username = ""
	view.password = ""
	view.Base.Reset()
}

func (view *Login) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Section(
			addClass("section"),

			addCenter(4,
				elem.Div(
					addClass("card"),

					elem.Header(
						addClass("card-header"),
						elem.Paragraph(
							addClass("card-header-title"),
							vecty.Text("Login"),
						),
					),

					elem.Div(
						addClass("card-content"),

						elem.Div(
							addClass("field"),

							elem.Div(
								addClass("control", "has-icons-left"),

								addInput([]vecty.MarkupOrChild{
									addClass("input", "is-primary"),
									addProprety("type", "text"),
									addProprety("placeholder", "Username"),
								}, func(s string) {
									view.username = s
								}),

								elem.Span(
									addClass("icon", "is-small", "is-left"),
									addIcon("user"),
								),
							),
						),

						elem.Div(
							addClass("field"),

							elem.Div(
								addClass("control", "has-icons-left"),

								addInput([]vecty.MarkupOrChild{
									addClass("input", "is-primary"),
									addProprety("type", "password"),
									addProprety("placeholder", "Password"),
								}, func(s string) {
									view.password = s
								}),

								elem.Span(
									addClass("icon", "is-small", "is-left"),
									addIcon("lock"),
								),
							),
						),

						elem.Div(
							addClass("field"),
							elem.Div(
								addClass("control"),
								elem.Anchor(
									addClass("button", "is-success"),
									addText("Login"),
									onClick(func() {
										view.SetNode("")

										control.DispatchEvent(
											control.NewEvent(control.ELogin).
												Set("username", view.username).
												Set("password", view.password))
									}),
								),
							),
						),
					),
				),
			),

			view.renderNote(4),
		),
	)
}
