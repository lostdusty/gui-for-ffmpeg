package kernel

import (
	"errors"
)

type Queue struct {
	Setting *ConvertSetting
	Status  StatusContract
	Error   error
}

type File struct {
	Path string
	Name string
	Ext  string
}

type ConvertSetting struct {
	VideoFileInput       File
	VideoFileOut         File
	OverwriteOutputFiles bool
}

type StatusContract interface {
	name() string
	ordinal() int
}

const (
	Waiting = iota
	InProgress
	Completed
	Error
)

type StatusType uint

var statusTypeStrings = []string{
	"waiting",
	"inProgress",
	"completed",
	"error",
}

func (status StatusType) name() string {
	return statusTypeStrings[status]
}

func (status StatusType) ordinal() int {
	return int(status)
}

type QueueListenerContract interface {
	Add(key int, queue *Queue)
	Remove(key int)
}

type QueueListContract interface {
	AddListener(queueListener QueueListenerContract)
	GetItems() map[int]*Queue
	Add(setting *ConvertSetting)
	Remove(key int)
	GetItem(key int) (*Queue, error)
	Next() (key int, queue *Queue)
}

type QueueList struct {
	currentKey    *int
	items         map[int]*Queue
	queueListener map[int]QueueListenerContract
}

func NewQueueList() *QueueList {
	currentKey := 0
	return &QueueList{
		currentKey:    &currentKey,
		items:         map[int]*Queue{},
		queueListener: map[int]QueueListenerContract{},
	}
}

func (l QueueList) GetItems() map[int]*Queue {
	return l.items
}

func (l QueueList) Add(setting *ConvertSetting) {
	queue := Queue{
		Setting: setting,
		Status:  StatusType(Waiting),
	}

	*l.currentKey += 1
	l.items[*l.currentKey] = &queue
	l.eventAdd(*l.currentKey, &queue)
}

func (l QueueList) Remove(key int) {
	if _, ok := l.items[key]; ok {
		delete(l.items, key)
		l.eventRemove(key)
	}
}

func (l QueueList) GetItem(key int) (*Queue, error) {
	if item, ok := l.items[key]; ok {
		return item, nil
	}

	return nil, errors.New("key not found")
}

func (l QueueList) AddListener(queueListener QueueListenerContract) {
	l.queueListener[len(l.queueListener)] = queueListener
}

func (l QueueList) eventAdd(key int, queue *Queue) {
	for _, listener := range l.queueListener {
		listener.Add(key, queue)
	}
}

func (l QueueList) eventRemove(key int) {
	for _, listener := range l.queueListener {
		listener.Remove(key)
	}
}

func (l QueueList) Next() (key int, queue *Queue) {
	statusWaiting := StatusType(Waiting)
	for key, item := range l.items {
		if item.Status == statusWaiting {
			return key, item
		}
	}
	return -1, nil
}
