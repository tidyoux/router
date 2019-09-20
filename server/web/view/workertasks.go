package view

import (
	"fmt"
	"strings"
	"time"

	"github.com/tidyoux/router/common/types"
	"github.com/tidyoux/router/server/web/cache"
	"github.com/tidyoux/router/server/web/control"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/tidyoux/lisc"
)

type WorkerTasks struct {
	taskList   *TaskList
	createTask *CreateTask

	worker *types.Worker
}

func NewWorkerTasks() *WorkerTasks {
	return &WorkerTasks{
		taskList:   NewTaskList(),
		createTask: NewCreateTask(),
	}
}

func (view *WorkerTasks) Reset() {
	view.worker = nil
	view.taskList.Reset()
	view.createTask.Reset()
}

func (view *WorkerTasks) SetWorker(worker *types.Worker) {
	view.worker = worker
	view.taskList.SetWorkerID(worker.ID)
}

func (view *WorkerTasks) Render() vecty.ComponentOrHTML {
	if view.worker != nil {
		if w, ok := cache.C().WorkerByID(view.worker.ID); ok {
			view.worker = w
		} else {
			view.worker = nil
			view.taskList.Reset()
		}
	}

	if view.worker == nil {
		return elem.Div()
	}

	return elem.Div(
		addClass("box"),
		view.renderHeader(),
		view.taskList,
		view.createTask,
	)
}

func (view *WorkerTasks) renderHeader() *vecty.HTML {
	var newBtn *vecty.HTML
	if view.worker.Status == types.WorkerEnabled {
		newBtn = elem.Anchor(
			addClass("level-item", "button", "is-success"),

			elem.Span(
				addText("New"),
			),

			onClick(func() {
				view.createTask.SetWorker(view.worker.ID)
				view.createTask.Active()
				rerender()
			}),
		)
	} else {
		newBtn = elem.Anchor(
			addClass("level-item", "button", "is-success"),
			addAttribute("disabled", ""),

			elem.Span(
				addText("New"),
			),
		)
	}

	return elem.Div(
		addClass("level"),

		elem.Div(
			addClass("level-left"),
			elem.Span(
				addClass("level-item"),
				elem.Strong(
					addText("Tasks:"),
				),
			),
		),

		elem.Div(
			addClass("level-right"),
			newBtn,
		),
	)
}

type TaskList struct {
	vecty.Core

	rendered bool
	workerID uint64
}

func NewTaskList() *TaskList {
	view := &TaskList{}
	go view.updateList()
	return view
}

func (view *TaskList) Reset() {
	view.workerID = 0
	cache.C().SetTasks(0, 0, nil)
}

func (view *TaskList) SetWorkerID(workerID uint64) {
	view.workerID = workerID
	cache.C().SetTasks(workerID, 0, nil)
}

func (view *TaskList) Render() vecty.ComponentOrHTML {
	defer func() { view.rendered = true }()

	var (
		weights    = []int{1, 4, 1, 1, 2, 2}
		titleColor = "is-info"
		nodes      = make([]vecty.MarkupOrChild, 0, 1+len(cache.C().Tasks()))
	)

	nodes = append(nodes, addColumns(0, []vecty.MarkupOrChild{
		elem.Span(
			addClass("tag", titleColor),
			addText("ID"),
		),

		elem.Span(
			addClass("tag", titleColor),
			addText("Params"),
		),

		elem.Span(
			addClass("tag", titleColor),
			addText("By"),
		),

		elem.Span(
			addClass("tag", titleColor),
			addText("Duration"),
		),

		elem.Span(
			addClass("tag", titleColor),
			addText("CreatedAt"),
		),

		elem.Span(
			addClass("tag", titleColor),
			addText("UpdatedAt"),
		),
	}, weights))

	for i := 0; i < len(cache.C().Tasks()); i++ {
		task := cache.C().Tasks()[i]

		var (
			statusColor = "is-light"
			statusIcon  = "spinner"
			duration    = task.UpdatedAt - task.CreatedAt
			now         = time.Now().Unix()
		)
		switch task.Status {
		case types.TaskRecord:
			if task.UpdatedAt < now {
				duration = now - task.CreatedAt
			}
		case types.TaskAccepted:
			statusColor = "is-loading"
			if task.UpdatedAt < now {
				duration = now - task.CreatedAt
			}
		case types.TaskSuccess:
			statusColor = "is-success"
			statusIcon = "check-circle"
		case types.TaskFailed:
			statusColor = "is-danger"
			statusIcon = "exclamation-circle"
		}

		nodes = append(nodes, addColumns(0, []vecty.MarkupOrChild{
			elem.Div(
				elem.Button(
					addClass("button", "is-small", statusColor),
					elem.Span(
						addClass("icon"),
						addIcon(statusIcon),
					),
					elem.Span(
						addText(fmt.Sprintf("#%d", task.ID)),
					),
				),
			),

			elem.Preformatted(
				addText(task.Params),
			),

			elem.Span(
				addClass("tag"),
				addText(task.Creator),
			),

			elem.Span(
				addClass("tag"),
				addText((time.Duration(duration) * time.Second).String()),
			),

			elem.Span(
				addClass("tag"),
				addText(time.Unix(task.CreatedAt, 0).Format("2006-01-02 15:04:05")),
			),

			elem.Span(
				addClass("tag"),
				addText(time.Unix(task.UpdatedAt, 0).Format("2006-01-02 15:04:05")),
			),
		}, weights))

		nodes = append(nodes, addColumns(1, []vecty.MarkupOrChild{
			elem.Preformatted(
				addText(task.Detail),
			),
		}, []int{10}))
	}

	return elem.Div(nodes...)
}

func (view *TaskList) updateList() {
	for {
		time.Sleep(time.Second)
		if view.workerID > 0 {
			control.DispatchEvent(
				control.NewEvent(control.EListTask).
					Set("workerID", view.workerID))
			if view.rendered {
				vecty.Rerender(view)
			} else {
				rerender()
			}
		}
	}
}

type Params struct {
	Value string
	Subs  []*Params
}

type CreateTask struct {
	vecty.Core
	Modal
	workerID    uint64
	paramFormat []*Params

	params    []string
	extParams string
}

func NewCreateTask() *CreateTask {
	return &CreateTask{}
}

func (view *CreateTask) SetWorker(workerID uint64) {
	view.workerID = workerID
	err := view.parseParamFormat()
	if err != nil {
		fmt.Printf("parse task params format failed, %v\n", err)
	}
}

func parseParamFormat(pair *lisc.Pair) ([]*Params, error) {
	if pair.ValueCount() == 0 {
		return nil, nil
	}

	params := make([]*Params, 0, pair.ValueCount())
	for i := 0; i < pair.ValueCount(); i++ {
		v, err := pair.Pair(i)
		if err != nil {
			return nil, fmt.Errorf("parse %s at index %d failed, %v",
				pair.Key(), i, err)
		}

		if len(v.Key()) == 0 {
			return nil, fmt.Errorf("parse %s at index %d failed, key can't be emtpy",
				pair.Key(), i)
		}

		subs, err := parseParamFormat(v)
		if err != nil {
			return nil, err
		}

		p := &Params{
			Value: v.Key(),
			Subs:  subs,
		}

		params = append(params, p)
	}

	return params, nil
}

func (view *CreateTask) parseParamFormat() error {
	worker, ok := cache.C().WorkerByID(view.workerID)
	if !ok {
		return fmt.Errorf("can't find worker by id")
	}

	format := lisc.New()
	err := format.Parse(worker.Desc)
	if err != nil {
		return err
	}

	params, err := parseParamFormat(format.Pair)
	if err != nil {
		return err
	}

	view.paramFormat = params
	return nil
}

func (view *CreateTask) Reset() {
	view.workerID = 0
	view.paramFormat = nil
	view.params = nil
	view.extParams = ""
	view.Modal.Reset()
}

func (view *CreateTask) Render() vecty.ComponentOrHTML {
	if view.workerID == 0 {
		return elem.Div()
	}

	if len(view.paramFormat) == 0 {
		return view.Modal.Render("Create task:", elem.Div(
			addClass("field", "has-addons"),
			view.renderInput(),
			view.renderSendButton(),
		), view.Reset)
	}

	markups := []vecty.MarkupOrChild{
		addClass("field", "has-addons"),
	}
	markups = append(markups, view.renderInputs()...)
	markups = append(markups, view.renderSendButton())

	return view.Modal.Render("Create task:", elem.Div(markups...), view.Reset)
}

func (view *CreateTask) renderInputs() []vecty.MarkupOrChild {
	var (
		inputs []vecty.MarkupOrChild
		subs   = view.paramFormat
		idx    = 0
	)
	for len(subs) > 0 {
		if subs[0].Value == "-" {
			inputs = append(inputs, view.renderInput())
			break
		}

		var node vecty.MarkupOrChild
		node, subs = view.renderSelect(idx, subs)
		inputs = append(inputs, node)
		idx++
	}
	return inputs
}

func (view *CreateTask) renderSelect(idx int, params []*Params) (vecty.MarkupOrChild, []*Params) {
	var subs []*Params
	options := make([]string, 0, len(params))
	for i, p := range params {
		if i == 0 && idx >= len(view.params) {
			view.params = append(view.params, p.Value)
		}

		options = append(options, p.Value)
		if p.Value == view.params[idx] {
			subs = p.Subs
		}
	}

	return addSelect(view.params[idx], options, func(value string) {
		view.params[idx] = value
		view.params = view.params[:idx+1]
		view.extParams = ""
		vecty.Rerender(view)
	}), subs
}

func (view *CreateTask) renderInput() vecty.MarkupOrChild {
	return elem.Div(
		addClass("control", "is-expanded"),

		addInput([]vecty.MarkupOrChild{
			addClass("input"),
			addProprety("type", "text"),
			addProprety("placeholder", "params"),
			addProprety("value", ""),
		}, func(s string) {
			view.extParams = s
		}),
	)
}

func (view *CreateTask) renderSendButton() vecty.MarkupOrChild {
	return elem.Div(
		addClass("control"),

		elem.Anchor(
			addClass("button", "is-success"),
			elem.Span(
				addClass("icon"),
				addIcon("paper-plane"),
			),
			onClick(func() {
				var params string
				if len(view.params) == 0 {
					params = view.extParams
				} else {
					params = strings.Join(view.params, " ")
					if len(view.extParams) > 0 {
						params += " " + view.extParams
					}
				}
				control.DispatchEvent(
					control.NewEvent(control.ESendTask).
						Set("workerID", view.workerID).
						Set("params", params))
				view.Reset()
			}),
		),
	)
}
