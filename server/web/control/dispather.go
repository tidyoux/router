package control

var (
	listeners = make(map[string]map[int]Listener)
	idCounter int
)

type Listener func(*Event)

func AddListener(eventName string, l func(*Event)) int {
	ls, ok := listeners[eventName]
	if !ok {
		ls = make(map[int]Listener)
		listeners[eventName] = ls
	}

	idCounter++
	ls[idCounter] = l
	return idCounter
}

func RemoveListener(e *Event, id int) {
	ls, ok := listeners[e.Name()]
	if !ok {
		return
	}

	delete(ls, id)
}

func DispatchEvent(e *Event) {
	ls, ok := listeners[e.Name()]
	if !ok {
		return
	}

	for _, l := range ls {
		go l(e)
	}
}
