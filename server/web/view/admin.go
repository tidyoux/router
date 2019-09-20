package view

import (
	"fmt"
	"strings"

	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server/web/cache"
	"github.com/tidyoux/router/server/web/control"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Admin struct {
	vecty.Core
	updateDetail *UpdateUserDetail
	resetUser    *ResetUser

	selectWorkerUser *SelectWorkerUser
	removeWorker     *RemoveWorker
}

func NewAdmin() *Admin {
	return &Admin{
		updateDetail:     NewUpdateUserDetail(),
		resetUser:        NewResetUser(),
		selectWorkerUser: NewSelectWorkerUser(),
		removeWorker:     NewRemoveWorker(),
	}
}

func (view *Admin) Reset() {
	view.updateDetail.Reset()
	view.resetUser.Reset()
	view.selectWorkerUser.Reset()
	view.removeWorker.Reset()
}

func (view *Admin) Render() vecty.ComponentOrHTML {
	return elem.Div(
		addColumns(1, []vecty.MarkupOrChild{
			view.renderUserList(),
			view.renderWorkerList(),
		}, []int{4, 6}),
		view.updateDetail,
		view.resetUser,
		view.selectWorkerUser,
		view.removeWorker,
	)
}

func (view *Admin) renderUserList() *vecty.HTML {
	nodes := make([]vecty.MarkupOrChild, 0, len(cache.C().Users())+2)
	nodes = append(nodes, addClass("panal"))
	nodes = append(nodes, elem.Paragraph(
		addClass("panel-heading"),
		addText(fmt.Sprintf("users: %d", len(cache.C().Users()))),
	))

	for i := 0; i < len(cache.C().Users()); i++ {
		user := cache.C().Users()[i]

		var (
			statusColor = "has-text-grey"
			opName      = "Activate"
			opEventName = control.EEnableUser
			opColor     = "is-success"
		)
		if user.Status == types.UserEnabled {
			statusColor = "has-text-success"
			opName = "Deactivate"
			opEventName = control.EDisableUser
			opColor = "is-gray"
		}

		var (
			detailColor = "has-text-dark"
			detail      = user.Detail
		)
		if len(user.Detail) == 0 {
			detailColor = "has-text-grey"
			detail = "click to edit detail"
		}

		nodes = append(nodes, elem.Div(
			addClass("panel-block"),

			elem.Span(
				addClass("panel-icon", statusColor),
				elem.Article(
					addIcon("user"),
				),
			),

			elem.Span(
				addClass("tag", "has-text-link"),
				addText(user.Name),
			),

			elem.Div(
				addClass("control"),

				elem.Span(
					addClass(detailColor),
					addText(detail),
					onClick(func() {
						view.updateDetail.SetUser(user)
						view.updateDetail.Active()
						rerender()
					}),
				),
			),

			elem.Div(
				addClass("field", "is-grouped"),

				elem.Div(
					addClass("control"),

					elem.Button(
						addClass("button", opColor),
						addText(opName),
						onClick(func() {
							control.DispatchEvent(
								control.NewEvent(opEventName).
									Set("userID", user.ID))
						}),
					),
				),

				elem.Div(
					addClass("control"),
					elem.Button(
						addClass("button", "has-text-danger"),
						addText("Reset"),
						onClick(func() {
							view.resetUser.SetUser(user.ID)
							view.resetUser.Active()
							rerender()
						}),
					),
				),
			),
		))
	}

	return elem.Div(nodes...)
}

func (view *Admin) userNames(userIDs []uint64) []string {
	names := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		if u, ok := cache.C().UserByID(id); ok {
			names = append(names, u.Name)
		}
	}
	return names
}

func (view *Admin) renderWorkerList() *vecty.HTML {
	nodes := make([]vecty.MarkupOrChild, 0, len(cache.C().Workers())+2)
	nodes = append(nodes, addClass("panal"))
	nodes = append(nodes, elem.Paragraph(
		addClass("panel-heading"),
		addText(fmt.Sprintf("workers: %d", len(cache.C().Workers()))),
	))

	for i := 0; i < len(cache.C().Workers()); i++ {
		worker := cache.C().Workers()[i]

		var (
			statusColor = "has-text-grey"
			opName      = "Activate"
			opEventName = control.EEnableWorker
			opColor     = "is-success"
		)
		if worker.Status == types.WorkerEnabled {
			statusColor = "has-text-success"
			opName = "Deactivate"
			opEventName = control.EDisableWorker
			opColor = "is-gray"
		}

		var (
			usersColor = "has-text-grey"
			userNames  = "click to add owners"
		)
		if len(worker.Users) > 0 {
			usersColor = "has-text-info"
			userNames = "[" + strings.Join(view.userNames(worker.Users), "|") + "]"
		}

		nodes = append(nodes, elem.Div(
			addClass("panel-block"),

			elem.Span(
				addClass("panel-icon", statusColor),
				elem.Article(
					addIcon("user-astronaut"),
				),
			),

			elem.Span(
				addClass("tag", "has-text-link"),
				addText(worker.Name),
			),

			elem.Div(
				addClass("control"),

				elem.Span(
					addClass(usersColor),
					addText(userNames),
					onClick(func() {
						view.selectWorkerUser.SetWorker(worker)
						view.selectWorkerUser.Active()
						rerender()
					}),
				),
			),

			elem.Div(
				addClass("field", "is-grouped"),

				elem.Div(
					addClass("control"),

					elem.Button(
						addClass("button", opColor),
						addText(opName),
						onClick(func() {
							control.DispatchEvent(
								control.NewEvent(opEventName).
									Set("workerID", worker.ID))
						}),
					),
				),

				elem.Div(
					addClass("control"),
					elem.Button(
						addClass("button", "has-text-danger"),
						elem.Span(
							addClass("icon"),
							addIcon("trash-alt"),
						),
						onClick(func() {
							view.removeWorker.SetWorker(worker.ID)
							view.removeWorker.Active()
							rerender()
						}),
					),
				),
			),
		))
	}

	nodes = append(nodes, elem.Div(
		addClass("panel-block"),

		elem.Button(
			addClass("button", "is-outlined", "is-fullwidth"),

			elem.Span(
				addClass("panel-icon", "has-text-success"),
				addIcon("plus-square"),
			),

			onClick(func() {
				control.DispatchEvent(control.NewEvent(control.EAddWorker))
			}),
		),
	))

	return elem.Div(nodes...)
}

type UpdateUserDetail struct {
	vecty.Core
	Modal
	user *types.User

	detail string
}

func NewUpdateUserDetail() *UpdateUserDetail {
	return &UpdateUserDetail{}
}

func (view *UpdateUserDetail) SetUser(user *types.User) {
	view.user = user
}

func (view *UpdateUserDetail) Reset() {
	view.user = nil
	view.detail = ""
	view.Modal.Reset()
}

func (view *UpdateUserDetail) Render() vecty.ComponentOrHTML {
	var userDetail string
	if view.user != nil {
		userDetail = view.user.Detail
	}

	return view.Modal.Render("Update user detail:", elem.Div(
		elem.Div(
			addClass("field", "has-addons"),

			elem.Div(
				addClass("control", "has-icons-left", "is-expanded"),

				addInput([]vecty.MarkupOrChild{
					addClass("input"),
					addProprety("type", "text"),
					addProprety("placeholder", "Edit user detail"),
					addProprety("value", userDetail),
				}, func(s string) {
					view.detail = s
				}),

				elem.Span(
					addClass("icon", "is-small", "is-left"),
					addIcon("info"),
				),
			),

			elem.Div(
				addClass("control"),

				elem.Anchor(
					addClass("button", "is-success"),
					elem.Span(
						addClass("icon"),
						addIcon("paper-plane"),
					),
					onClick(func() {
						if view.detail != userDetail {
							control.DispatchEvent(
								control.NewEvent(control.EUpdateUserDetail).
									Set("userID", view.user.ID).
									Set("detail", view.detail))
						}

						view.Reset()
						vecty.Rerender(view)
					}),
				),
			),
		),
	), view.Reset)
}

type ResetUser struct {
	vecty.Core
	Modal
	userID uint64
}

func NewResetUser() *ResetUser {
	return &ResetUser{}
}

func (view *ResetUser) SetUser(userID uint64) {
	view.userID = userID
}

func (view *ResetUser) Reset() {
	view.userID = 0
	view.Modal.Reset()
}

func (view *ResetUser) Render() vecty.ComponentOrHTML {
	return view.Modal.Render("Reset user password:", elem.Div(
		elem.Div(
			addClass("field"),

			elem.Label(
				addClass("label"),
				addText("Reset user password to"),
				elem.Code(
					addText("123456"),
				),
				addText("."),
			),
		),

		elem.Anchor(
			addClass("button", "is-fullwidth"),
			elem.Span(
				addClass("has-text-danger"),
				addText("I understand the consequences, reset this user"),
			),
			onClick(func() {
				control.DispatchEvent(
					control.NewEvent(control.EResetUserPassword).
						Set("userID", view.userID))
				view.Reset()
				vecty.Rerender(view)
			}),
		),
	), view.Reset)
}

type SelectWorkerUser struct {
	vecty.Core
	Modal

	workerID     uint64
	workerOwners map[uint64]struct{}
}

func NewSelectWorkerUser() *SelectWorkerUser {
	return &SelectWorkerUser{}
}

func (view *SelectWorkerUser) SetWorker(worker *types.Worker) {
	view.workerID = worker.ID
	view.workerOwners = make(map[uint64]struct{}, len(worker.Users))
	for _, id := range worker.Users {
		view.workerOwners[id] = struct{}{}
	}
}

func (view *SelectWorkerUser) Reset() {
	view.workerID = 0
	view.workerOwners = nil
	view.Modal.Reset()
}

func (view *SelectWorkerUser) Render() vecty.ComponentOrHTML {
	nodes := make([]vecty.MarkupOrChild, 0, 1+len(cache.C().Users()))
	nodes = append(nodes, addClass("level"))
	for i := 0; i < len(cache.C().Users()); i++ {
		user := cache.C().Users()[i]
		_, isOwner := view.workerOwners[user.ID]

		nodes = append(nodes, elem.Label(
			addClass("checkbox"),
			elem.Input(
				addProprety("type", "checkbox"),
				addProprety("checked", isOwner),
				onCheckChange(func(checked bool) {
					if checked {
						view.workerOwners[user.ID] = struct{}{}
						control.DispatchEvent(
							control.NewEvent(control.EAddWorkerUser).
								Set("workerID", view.workerID).
								Set("userID", user.ID))
					} else {
						delete(view.workerOwners, user.ID)
						control.DispatchEvent(
							control.NewEvent(control.ERemoveWorkerUser).
								Set("workerID", view.workerID).
								Set("userID", user.ID))
					}
				}),
			),
			addText(" "+user.Name),
		))
	}
	return view.Modal.Render("Select owners:", elem.Div(nodes...), view.Reset)
}

const (
	RemoveWorkerConfirm = "delete"
)

type RemoveWorker struct {
	vecty.Core
	Modal
	workerID uint64

	confirm string
}

func NewRemoveWorker() *RemoveWorker {
	return &RemoveWorker{}
}

func (view *RemoveWorker) SetWorker(workerID uint64) {
	view.workerID = workerID
}

func (view *RemoveWorker) Reset() {
	view.workerID = 0
	view.confirm = ""
	view.Modal.Reset()
}

func (view *RemoveWorker) Render() vecty.ComponentOrHTML {
	var confirmBtn *vecty.HTML
	if strings.EqualFold(view.confirm, RemoveWorkerConfirm) {
		confirmBtn = elem.Anchor(
			addClass("button", "is-fullwidth"),
			elem.Span(
				addClass("has-text-danger"),
				addText("I understand the consequences, delete this worker"),
			),
			onClick(func() {
				control.DispatchEvent(
					control.NewEvent(control.ERemoveWorker).
						Set("workerID", view.workerID))
				view.Reset()
			}),
		)
	} else {
		confirmBtn = elem.Anchor(
			addClass("button", "is-fullwidth"),
			addAttribute("disabled", ""),
			elem.Span(
				addClass("has-text-danger"),
				addText("I understand the consequences, delete this worker"),
			),
		)
	}

	return view.Modal.Render("Are you absolutely sure?", elem.Div(
		elem.Div(
			addClass("field"),

			elem.Label(
				addClass("label"),
				addText("Please type in"),
				elem.Code(
					addText(RemoveWorkerConfirm),
				),
				addText("to confirm."),
			),

			elem.Div(
				addClass("control", "is-expanded"),
				addInput([]vecty.MarkupOrChild{
					addClass("input"),
					addProprety("type", "text"),
					addProprety("value", view.confirm),
				}, func(s string) {
					view.confirm = s
					vecty.Rerender(view)
				}),
			),
		),

		confirmBtn,
	), view.Reset)
}
