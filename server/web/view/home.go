package view

import (
	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server/web/cache"
	"github.com/tidyoux/router/server/web/control"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Home struct {
	Base
	updatePassword *UpdatePassword
	admin          *Admin
	worker         *Worker

	activeView vecty.Component
}

func NewHome() *Home {
	return &Home{
		updatePassword: NewUpdatePassword(),
		admin:          NewAdmin(),
		worker:         NewWorker(),
	}
}

func (view *Home) Init() {
	if cache.C().Username() == types.AdminUsername {
		view.activeView = view.admin
		control.DispatchEvent(control.NewEvent(control.EListUser))
	} else {
		view.activeView = view.worker
	}
	control.DispatchEvent(control.NewEvent(control.EListWorker))
}

func (view *Home) Reset() {
	view.admin.Reset()
	view.worker.Reset()
	view.Base.Reset()
}

func (view *Home) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Section(
			addClass("hero", "is-small"),

			elem.Div(
				addClass("hero-head"),
				view.renderNavbar(),
			),

			elem.Div(
				addClass("hero-body"),

				view.renderNote(1),

				view.activeView,
			),

			elem.Div(
				addClass("hero-foot"),
			),
		),
		view.updatePassword,
	)
}

func (view *Home) renderNavbar() *vecty.HTML {
	return elem.Navigation(
		addClass("navbar"),

		elem.Div(
			addClass("navbar-end"),

			elem.Div(
				addClass("navbar-item"),

				elem.Anchor(
					addClass("button", "is-link", "is-outlined"),
					elem.Span(
						addClass("icon"),
						addIcon("user-cog"),
					),
					elem.Span(
						addText(cache.C().Username()),
					),
					onClick(func() {
						view.updatePassword.Active()
						rerender()
					}),
				),
			),

			elem.Div(
				addClass("navbar-item"),

				elem.Anchor(
					addClass("button"),
					elem.Span(
						addClass("icon"),
						addIcon("sign-out-alt"),
					),
					elem.Span(
						addText("Logout"),
					),
					onClick(func() {
						control.DispatchEvent(control.NewEvent(control.ELogout))
					}),
				),
			),
		),
	)
}

type UpdatePassword struct {
	vecty.Core
	Modal
	passwordHelp *Component
	confirmHelp  *Component
	updateButton *Component

	password string
	confirm  string
}

func NewUpdatePassword() *UpdatePassword {
	view := &UpdatePassword{}
	view.passwordHelp = NewComponent(func() vecty.ComponentOrHTML {
		if len(view.password) == 0 {
			return elem.Paragraph(
				addClass("help", "is-danger"),
				addText("Password can't be empty"),
			)
		}

		return elem.Div()
	})
	view.confirmHelp = NewComponent(func() vecty.ComponentOrHTML {
		if view.confirm != view.password {
			return elem.Paragraph(
				addClass("help", "is-danger"),
				addText("Password mismatch"),
			)
		}

		return elem.Div()
	})
	view.updateButton = NewComponent(func() vecty.ComponentOrHTML {
		if len(view.password) > 0 && view.confirm == view.password {
			return elem.Anchor(
				addClass("button", "is-success"),
				elem.Span(
					addText("Update"),
				),
				onClick(func() {
					control.DispatchEvent(
						control.NewEvent(control.EUpdatePassword).
							Set("password", view.password))
					view.Reset()
				}),
			)
		}

		return elem.Anchor(
			addClass("button", "is-success"),
			addAttribute("disabled", ""),
			elem.Span(
				addText("Update"),
			),
		)
	})
	return view
}

func (view *UpdatePassword) Reset() {
	view.password = ""
	view.confirm = ""
	view.Modal.Reset()
}

func (view *UpdatePassword) rerender() {
	vecty.Rerender(view.passwordHelp)
	vecty.Rerender(view.confirmHelp)
	vecty.Rerender(view.updateButton)
}

func (view *UpdatePassword) Render() vecty.ComponentOrHTML {
	return view.Modal.Render("Update password:", elem.Div(
		elem.Div(
			addClass("field"),

			elem.Div(
				addClass("control", "has-icons-left", "is-expanded"),

				addInput([]vecty.MarkupOrChild{
					addClass("input"),
					addProprety("type", "password"),
					addProprety("placeholder", "Password"),
					addProprety("value", ""),
				}, func(s string) {
					view.password = s
					view.rerender()
				}),

				elem.Span(
					addClass("icon", "is-small", "is-left"),
					addIcon("lock"),
				),

				view.passwordHelp,
			),
		),

		elem.Div(
			addClass("field"),

			elem.Div(
				addClass("control", "has-icons-left", "is-expanded"),

				addInput([]vecty.MarkupOrChild{
					addClass("input"),
					addProprety("type", "password"),
					addProprety("placeholder", "Confirm"),
					addProprety("value", ""),
				}, func(s string) {
					view.confirm = s
					view.rerender()
				}),

				elem.Span(
					addClass("icon", "is-small", "is-left"),
					addIcon("lock"),
				),

				view.confirmHelp,
			),
		),

		elem.Div(
			addClass("field"),

			elem.Div(
				addClass("control"),
				view.updateButton,
			),
		),
	), view.Reset)
}
