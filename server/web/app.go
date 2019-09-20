package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/tidyoux/router/client"
	"github.com/tidyoux/router/server/web/cache"
	"github.com/tidyoux/router/server/web/control"
	"github.com/tidyoux/router/server/web/view"
	"github.com/gopherjs/vecty"
)

func main() {
	vecty.SetTitle("Router")
	vecty.AddStylesheet("https://cdnjs.cloudflare.com/ajax/libs/bulma/0.7.5/css/bulma.min.css")
	vecty.RenderBody(New())
}

type App struct {
	vecty.Core
	client *client.Client

	loginView *view.Login
	homeView  *view.Home

	activeView view.View

	hasLogin bool
	userName string
}

func New() *App {
	app := &App{
		client:    client.New(""),
		loginView: view.NewLogin(),
		homeView:  view.NewHome(),
	}
	app.Init()
	return app
}

func (a *App) Init() {
	control.AddListener(control.ERerender, a.onRerender)

	control.AddListener(control.ELogin, a.onLogin)
	control.AddListener(control.ELogout, a.onLogout)

	control.AddListener(control.EUpdatePassword, a.onUpdatePassword)

	control.AddListener(control.EListWorker, a.onListWorker)
	control.AddListener(control.EUpdateWorker, a.onUpdateWorker)

	control.AddListener(control.EListTask, a.onListTask)
	control.AddListener(control.ESendTask, a.onSendTask)

	control.AddListener(control.EListUser, a.onListUser)
	control.AddListener(control.EUpdateUserDetail, a.onUpdateUserDetail)
	control.AddListener(control.EResetUserPassword, a.onResetUserPassword)
	control.AddListener(control.EEnableUser, a.onEnableUser)
	control.AddListener(control.EDisableUser, a.onDisableUser)
	control.AddListener(control.EAddWorker, a.onAddWorker)
	control.AddListener(control.EEnableWorker, a.onEnableWorker)
	control.AddListener(control.EDisableWorker, a.onDisableWorker)
	control.AddListener(control.EAddWorkerUser, a.onAddWorkerUser)
	control.AddListener(control.ERemoveWorkerUser, a.onRemoveWorkerUser)
	control.AddListener(control.ERemoveWorker, a.onRemoveWorker)

	if !a.restoreTokenUsername() {
		a.changeView(a.loginView)
	} else {
		cache.C().SetUsername(a.userName)
		a.changeView(a.homeView)
	}
}

func (a *App) Render() vecty.ComponentOrHTML {
	return a.activeView.Render()
}

func (a *App) changeView(view view.View) {
	view.Reset()
	a.activeView = view
	a.activeView.Init()
}

func (a *App) storeTokenUsername() {
	localStore("token", a.client.Token())
	localStore("username", a.userName)
}

func (a *App) restoreTokenUsername() bool {
	if data := localLoad("token"); data != js.Undefined() {
		token := data.String()
		a.client.SetToken(token)

		if data := localLoad("username"); data != js.Undefined() {
			a.userName = data.String()

			err := a.client.Ping()
			if err == nil {
				return true
			}
		}
	}

	a.clearTokenUsername()
	return false
}

func (a *App) clearTokenUsername() {
	localStore("token", "")
	localStore("username", "")
	a.client.SetToken("")
	a.userName = ""
	cache.C().Clear()
}

func (a *App) onRerender(e *control.Event) {
	vecty.Rerender(a)
}

func (a *App) onLogin(e *control.Event) {
	username, _ := e.Get("username")
	password, _ := e.Get("password")
	_, err := a.client.Login(username.(string), password.(string))
	if err != nil {
		a.loginView.SetNode(err.Error())
	} else {
		a.userName = username.(string)
		a.storeTokenUsername()
		cache.C().SetUsername(a.userName)
		a.changeView(a.homeView)
	}

	vecty.Rerender(a)
}

func (a *App) onLogout(e *control.Event) {
	err := a.client.Logout()
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.clearTokenUsername()
		a.changeView(a.loginView)
	}

	vecty.Rerender(a)
}

func (a *App) onUpdatePassword(e *control.Event) {
	password, _ := e.Get("password")
	err := a.client.UpdatePassword(password.(string))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.clearTokenUsername()
		a.changeView(a.loginView)
	}

	vecty.Rerender(a)
}

func (a *App) updateWorkers() {
	workers, err := a.client.ListWorker()
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		cache.C().SetWorkers(workers)
	}
}

func (a *App) onListWorker(e *control.Event) {
	a.updateWorkers()
	vecty.Rerender(a)
}

func (a *App) onUpdateWorker(e *control.Event) {
	workerID, _ := e.Get("workerID")

	var err error
	if name, ok := e.Get("name"); ok {
		err = a.client.UpdateWorkerName(workerID.(uint64), name.(string))
	} else if desc, ok := e.Get("desc"); ok {
		err = a.client.UpdateWorkerDesc(workerID.(uint64), desc.(string))
	} else {
		return
	}

	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}

func (a *App) updateTasks(workerID uint64) {
	resp, err := a.client.ListTask(workerID, 0, 10)
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		cache.C().SetTasks(workerID, resp.Total, resp.Tasks)
	}
}

func (a *App) onListTask(e *control.Event) {
	workerID, _ := e.Get("workerID")
	a.updateTasks(workerID.(uint64))
}

func (a *App) onSendTask(e *control.Event) {
	workerID, _ := e.Get("workerID")
	params, _ := e.Get("params")
	_, err := a.client.SendTask(workerID.(uint64), params.(string))
	if err != nil {
		a.homeView.SetNode(err.Error())
	}

	a.updateTasks(workerID.(uint64))

	vecty.Rerender(a)
}

func (a *App) updateUsers() {
	users, err := a.client.ListUser()
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		cache.C().SetUsers(users)
	}
}

func (a *App) onListUser(e *control.Event) {
	a.updateUsers()
	vecty.Rerender(a)
}

func (a *App) onUpdateUserDetail(e *control.Event) {
	userID, _ := e.Get("userID")
	detail, _ := e.Get("detail")
	err := a.client.UpdateUserDetail(userID.(uint64), detail.(string))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateUsers()
	}

	vecty.Rerender(a)
}

func (a *App) onResetUserPassword(e *control.Event) {
	userID, _ := e.Get("userID")
	err := a.client.ResetUserPassword(userID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
		vecty.Rerender(a)
	}
}

func (a *App) onEnableUser(e *control.Event) {
	userID, _ := e.Get("userID")
	err := a.client.EnableUser(userID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateUsers()
	}

	vecty.Rerender(a)
}

func (a *App) onDisableUser(e *control.Event) {
	userID, _ := e.Get("userID")
	err := a.client.DisableUser(userID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateUsers()
	}

	vecty.Rerender(a)
}

func (a *App) onAddWorker(e *control.Event) {
	_, err := a.client.AddWorker(fmt.Sprintf("worker%d", time.Now().Unix()), "; Task command.\n\n(run (-))")
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}

func (a *App) onEnableWorker(e *control.Event) {
	workerID, _ := e.Get("workerID")
	err := a.client.EnableWorker(workerID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}

func (a *App) onDisableWorker(e *control.Event) {
	workerID, _ := e.Get("workerID")
	err := a.client.DisableWorker(workerID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}

func (a *App) onAddWorkerUser(e *control.Event) {
	workerID, _ := e.Get("workerID")
	userID, _ := e.Get("userID")
	err := a.client.AddWorkerUser(workerID.(uint64), userID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}

func (a *App) onRemoveWorkerUser(e *control.Event) {
	workerID, _ := e.Get("workerID")
	userID, _ := e.Get("userID")
	err := a.client.RemoveWorkerUser(workerID.(uint64), userID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}

func (a *App) onRemoveWorker(e *control.Event) {
	workerID, _ := e.Get("workerID")
	err := a.client.RemoveWorker(workerID.(uint64))
	if err != nil {
		a.homeView.SetNode(err.Error())
	} else {
		a.updateWorkers()
	}

	vecty.Rerender(a)
}
