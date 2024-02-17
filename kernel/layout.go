package kernel

import (
	"bufio"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type LayoutContract interface {
	SetContent(content fyne.CanvasObject) *fyne.Container
	NewProgressbar(queueId int, totalDuration float64) ProgressContract
	ChangeQueueStatus(queueId int, queue *Queue)
}

type Layout struct {
	layout            *fyne.Container
	queueLayoutObject QueueLayoutObjectContract
	localizerService  LocalizerContract
}

func NewLayout(queueLayoutObject QueueLayoutObjectContract, localizerService LocalizerContract) *Layout {
	layout := container.NewAdaptiveGrid(2, widget.NewLabel(""), container.NewVScroll(queueLayoutObject.GetCanvasObject()))

	return &Layout{
		layout:            layout,
		queueLayoutObject: queueLayoutObject,
		localizerService:  localizerService,
	}
}

func (l Layout) SetContent(content fyne.CanvasObject) *fyne.Container {
	l.layout.Objects[0] = content
	return l.layout
}

func (l Layout) NewProgressbar(queueId int, totalDuration float64) ProgressContract {
	progressbar := l.queueLayoutObject.GetProgressbar(queueId)
	return NewProgress(totalDuration, progressbar, l.localizerService)
}

func (l Layout) ChangeQueueStatus(queueId int, queue *Queue) {
	l.queueLayoutObject.ChangeQueueStatus(queueId, queue)
}

type QueueLayoutObjectContract interface {
	GetCanvasObject() fyne.CanvasObject
	GetProgressbar(queueId int) *widget.ProgressBar
	ChangeQueueStatus(queueId int, queue *Queue)
}

type QueueLayoutObject struct {
	QueueListContract QueueListContract

	queue                   QueueListContract
	container               *fyne.Container
	items                   map[int]QueueLayoutItem
	localizerService        LocalizerContract
	layoutLocalizerListener LayoutLocalizerListenerContract
}

type QueueLayoutItem struct {
	CanvasObject  fyne.CanvasObject
	ProgressBar   *widget.ProgressBar
	StatusMessage *canvas.Text
	MessageError  *canvas.Text
}

func NewQueueLayoutObject(queue QueueListContract, localizerService LocalizerContract, layoutLocalizerListener LayoutLocalizerListenerContract) *QueueLayoutObject {
	title := widget.NewLabel(localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: "queue"}) + ":")
	title.TextStyle.Bold = true

	layoutLocalizerListener.AddItem("queue", title)

	queueLayoutObject := &QueueLayoutObject{
		queue:                   queue,
		container:               container.NewVBox(title),
		items:                   map[int]QueueLayoutItem{},
		localizerService:        localizerService,
		layoutLocalizerListener: layoutLocalizerListener,
	}

	queue.AddListener(queueLayoutObject)

	return queueLayoutObject
}

func (o QueueLayoutObject) GetCanvasObject() fyne.CanvasObject {
	return o.container
}

func (o QueueLayoutObject) GetProgressbar(queueId int) *widget.ProgressBar {
	if item, ok := o.items[queueId]; ok {
		return item.ProgressBar
	}
	return widget.NewProgressBar()
}

func (o QueueLayoutObject) Add(id int, queue *Queue) {
	progressBar := widget.NewProgressBar()
	statusMessage := canvas.NewText(o.getStatusTitle(queue.Status), theme.PrimaryColor())
	messageError := canvas.NewText("", theme.ErrorColor())

	content := container.NewVBox(
		container.NewHScroll(widget.NewLabel(queue.Setting.VideoFileInput.Name)),
		progressBar,
		container.NewHScroll(statusMessage),
		container.NewHScroll(messageError),
		canvas.NewLine(theme.FocusColor()),
		container.NewPadded(),
	)
	o.items[id] = QueueLayoutItem{
		CanvasObject:  content,
		ProgressBar:   progressBar,
		StatusMessage: statusMessage,
		MessageError:  messageError,
	}
	o.container.Add(content)
}

func (o QueueLayoutObject) Remove(id int) {
	if item, ok := o.items[id]; ok {
		o.container.Remove(item.CanvasObject)
		o.items[id] = QueueLayoutItem{}
	}
}

func (o QueueLayoutObject) ChangeQueueStatus(queueId int, queue *Queue) {
	if item, ok := o.items[queueId]; ok {
		statusColor := o.getStatusColor(queue.Status)
		item.StatusMessage.Text = o.getStatusTitle(queue.Status)
		item.StatusMessage.Color = statusColor
		item.StatusMessage.Refresh()
		if queue.Error != nil {
			item.MessageError.Text = queue.Error.Error()
			item.MessageError.Color = statusColor
			item.MessageError.Refresh()
		}
	}
}

func (o QueueLayoutObject) getStatusColor(status StatusContract) color.Color {
	if status == StatusType(Error) {
		return theme.ErrorColor()
	}

	if status == StatusType(Completed) {
		return color.RGBA{R: 49, G: 127, B: 114, A: 255}
	}

	return theme.PrimaryColor()
}

func (o QueueLayoutObject) getStatusTitle(status StatusContract) string {
	return o.localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: status.name()})
}

type Progress struct {
	totalDuration    float64
	progressbar      *widget.ProgressBar
	protocol         string
	localizerService LocalizerContract
}

func NewProgress(totalDuration float64, progressbar *widget.ProgressBar, localizerService LocalizerContract) Progress {
	return Progress{
		totalDuration:    totalDuration,
		progressbar:      progressbar,
		protocol:         "pipe:",
		localizerService: localizerService,
	}
}

func (p Progress) GetProtocole() string {
	return p.protocol
}

func (p Progress) Run(stdOut io.ReadCloser, stdErr io.ReadCloser) error {
	isProcessCompleted := false
	var errorText string

	p.progressbar.Value = 0
	p.progressbar.Max = p.totalDuration
	p.progressbar.Refresh()

	progress := 0.0

	go func() {
		scannerErr := bufio.NewReader(stdErr)
		for {
			line, _, err := scannerErr.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}
			data := strings.TrimSpace(string(line))
			errorText = data
		}
	}()

	scannerOut := bufio.NewReader(stdOut)
	for {
		line, _, err := scannerOut.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		data := strings.TrimSpace(string(line))
		if strings.Contains(data, "progress=end") {
			p.progressbar.Value = p.totalDuration
			p.progressbar.Refresh()
			isProcessCompleted = true
			break
		}

		re := regexp.MustCompile(`frame=(\d+)`)
		a := re.FindAllStringSubmatch(data, -1)

		if len(a) > 0 && len(a[len(a)-1]) > 0 {
			c, err := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
			if err != nil {
				continue
			}
			progress = float64(c)
		}
		if p.progressbar.Value != progress {
			p.progressbar.Value = progress
			p.progressbar.Refresh()
		}
	}

	if isProcessCompleted == false {
		if len(errorText) == 0 {
			errorText = p.localizerService.GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorConverter",
			})
		}
		return errors.New(errorText)
	}

	return nil
}

type LayoutLocalizerItem struct {
	messageID string
	object    *widget.Label
}

type LayoutLocalizerListener struct {
	itemCurrentId int
	items         map[int]*LayoutLocalizerItem
}

type LayoutLocalizerListenerContract interface {
	AddItem(messageID string, object *widget.Label)
}

func NewLayoutLocalizerListener() *LayoutLocalizerListener {
	return &LayoutLocalizerListener{
		itemCurrentId: 0,
		items:         map[int]*LayoutLocalizerItem{},
	}
}

func (l LayoutLocalizerListener) AddItem(messageID string, object *widget.Label) {
	l.itemCurrentId += 1
	l.items[l.itemCurrentId] = &LayoutLocalizerItem{messageID: messageID, object: object}
}

func (l LayoutLocalizerListener) Change(localizerService LocalizerContract) {
	for _, item := range l.items {
		item.object.Text = localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: item.messageID})
		item.object.Refresh()
	}
}
