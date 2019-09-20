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

type WorkerDetail struct {
	worker *types.Worker

	workerName string
	editName   bool

	workerDesc string
	editDesc   bool
}

func NewWorkerDetail() *WorkerDetail {
	return &WorkerDetail{}
}

func (view *WorkerDetail) Reset() {
	view.resetEdit()
	view.worker = nil
}

func (view *WorkerDetail) resetEdit() {
	view.workerName = ""
	view.editName = false
	view.workerDesc = ""
	view.editDesc = false
}

func (view *WorkerDetail) SetWorker(worker *types.Worker) {
	view.worker = worker
}

func (view *WorkerDetail) Worker() *types.Worker {
	return view.worker
}

func (view *WorkerDetail) Render() vecty.ComponentOrHTML {
	if view.worker != nil {
		if w, ok := cache.C().WorkerByID(view.worker.ID); ok {
			view.worker = w
		} else {
			view.worker = nil
		}
	}

	if view.worker == nil {
		return elem.Div()
	}

	return elem.Div(
		addClass("card"),
		view.renderHeader(),
		view.renderBody(),
	)
}

func (view *WorkerDetail) renderHeader() *vecty.HTML {
	var nameNode *vecty.HTML
	if view.editName {
		view.workerName = view.worker.Name
		nameNode = elem.Div(
			addClass("field", "is-grouped"),
			elem.Div(
				addClass("control"),
				addInput([]vecty.MarkupOrChild{
					addClass("input", "is-primary"),
					addProprety("type", "text"),
					addProprety("value", view.worker.Name),
				}, func(s string) {
					view.workerName = s
				}),
			),

			elem.Div(
				addClass("control"),
				elem.Anchor(
					addClass("button", "is-info"),
					addText("Save"),
					onClick(func() {
						name := view.workerName
						view.resetEdit()
						if name != view.worker.Name {
							control.DispatchEvent(
								control.NewEvent(control.EUpdateWorker).
									Set("workerID", view.worker.ID).
									Set("name", name))
						} else {
							rerender()
						}
					}),
				),
			),

			view.renderCancelButton(),
		)
	} else {
		var editBtn *vecty.HTML
		if view.worker.Status == types.WorkerEnabled {
			editBtn = elem.Anchor(
				elem.Span(
					addClass("icon"),
					addIcon("edit"),
				),
				onClick(func() {
					view.editName = true
					rerender()
				}),
			)
		} else {
			editBtn = elem.Anchor(
				elem.Span(
					addClass("icon", "has-text-grey"),
					addIcon("edit"),
				),
			)
		}

		nameNode = elem.Div(
			elem.Span(
				addClass("title", "is-4"),
				addText(view.worker.Name),
			),
			editBtn,
		)
	}

	return elem.Header(
		addClass("card-header"),
		elem.Div(
			addClass("card-header-title"),
			elem.Div(
				addClass("content"),

				nameNode,

				elem.Div(
					elem.Span(
						addClass("has-text-grey"),
						addText(fmt.Sprintf("id: %d", view.worker.ID)),
					),
				),
			),
		),
	)
}

func (view *WorkerDetail) renderBody() *vecty.HTML {
	var descNode *vecty.HTML
	if view.editDesc {
		view.workerDesc = view.worker.Desc
		const MinTextAreaLines = 5
		lines := strings.Count(view.worker.Desc, "\n") + 1
		if lines < MinTextAreaLines {
			lines = MinTextAreaLines
		}

		descNode = elem.Div(
			elem.Div(
				addClass("field"),
				elem.Div(
					addClass("control"),
					elem.TextArea(
						addClass("textarea"),
						addProprety("value", view.worker.Desc),
						addProprety("rows", lines),
						onInput(func(s string) {
							view.workerDesc = s
						}),
					),
				),
			),

			elem.Div(
				addClass("field", "is-grouped"),
				elem.Div(
					addClass("control"),
					elem.Anchor(
						addClass("button", "is-info"),
						addText("Save"),
						onClick(func() {
							desc := view.workerDesc
							view.resetEdit()
							if desc != view.worker.Desc {
								control.DispatchEvent(
									control.NewEvent(control.EUpdateWorker).
										Set("workerID", view.worker.ID).
										Set("desc", desc))
							} else {
								rerender()
							}
						}),
					),
				),

				view.renderCancelButton(),
			),
		)
	} else {
		var editBtn *vecty.HTML
		if view.worker.Status == types.WorkerEnabled {
			editBtn = elem.Anchor(
				addClass("level-item"),
				elem.Span(
					addClass("icon"),
					addIcon("edit"),
				),
				onClick(func() {
					view.editDesc = true
					rerender()
				}),
			)
		} else {
			editBtn = elem.Anchor(
				addClass("level-item"),
				elem.Span(
					addClass("icon", "has-text-grey"),
					addIcon("edit"),
				),
			)
		}

		descNode = elem.Div(
			addClass("has-background-light"),

			elem.Div(
				addClass("level"),

				elem.Div(
					addClass("level-left"),
					elem.Span(
						addClass("level-item", "tag"),
						addText("Description:"),
					),
				),

				elem.Div(
					addClass("level-right"),
					editBtn,
				),
			),

			elem.Paragraph(
				elem.Preformatted(
					addText(view.worker.Desc),
				),
			),
		)
	}

	return elem.Div(
		addClass("card-content"),
		elem.Div(
			addClass("content"),

			elem.Paragraph(
				elem.Span(
					addClass("tag", "has-text-danger"),
					addText(fmt.Sprintf("key: %s", view.worker.Key)),
				),
			),

			elem.Div(
				addClass("box"),

				descNode,
			),
		),
	)
}

func (view *WorkerDetail) renderCancelButton() *vecty.HTML {
	return elem.Div(
		addClass("control"),
		elem.Anchor(
			addClass("button", "is-text"),
			addText("Cancel"),
			onClick(func() {
				view.resetEdit()
				rerender()
			}),
		),
	)
}
