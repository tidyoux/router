package view

import (
	"fmt"

	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server/web/cache"
	"github.com/tidyoux/router/server/web/control"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Worker struct {
	vecty.Core

	workerDetail *WorkerDetail
	workerTasks  *WorkerTasks
}

func NewWorker() *Worker {
	return &Worker{
		workerDetail: NewWorkerDetail(),
		workerTasks:  NewWorkerTasks(),
	}
}

func (view *Worker) Reset() {
	view.workerTasks.Reset()
	view.workerDetail.Reset()
}

func (view *Worker) Render() vecty.ComponentOrHTML {
	return addColumns(0, []vecty.MarkupOrChild{
		view.renderWorkerList(),
		elem.Div(
			view.workerDetail.Render(),
			view.workerTasks.Render(),
		),
	}, []int{2, 10})
}

func (view *Worker) renderWorkerList() *vecty.HTML {
	nodes := make([]vecty.MarkupOrChild, 0, len(cache.C().Workers())+2)
	nodes = append(nodes, addClass("panal"))
	nodes = append(nodes, elem.Paragraph(
		addClass("panel-heading"),
		addText(fmt.Sprintf("workers: %d", len(cache.C().Workers()))),
	))

	for i := 0; i < len(cache.C().Workers()); i++ {
		worker := cache.C().Workers()[i]
		statusColor := "has-text-grey"
		if worker.Status == types.WorkerEnabled {
			statusColor = "has-text-success"
		}

		nameBtnCls := []string{"button", "is-outlined", "is-link", "is-fullwidth"}
		if view.workerDetail.Worker() != nil && worker.ID == view.workerDetail.Worker().ID {
			nameBtnCls = append(nameBtnCls, "is-focused")
		}

		nodes = append(nodes, elem.Div(
			addClass("panel-block"),

			elem.Span(
				addClass("panel-icon", statusColor),
				elem.Article(
					addIcon("user-astronaut"),
				),
			),

			elem.Button(
				addClass(nameBtnCls...),
				addText(worker.Name),
				onClick(func() {
					if view.workerDetail.Worker() == nil || worker.ID != view.workerDetail.Worker().ID {
						view.workerDetail.Reset()
						view.workerDetail.SetWorker(worker)

						view.workerTasks.Reset()
						view.workerTasks.SetWorker(worker)

						rerender()
						control.DispatchEvent(control.NewEvent(control.EListTask).Set("workerID", worker.ID))
					}
				}),
			),
		))
	}

	return elem.Div(nodes...)
}
