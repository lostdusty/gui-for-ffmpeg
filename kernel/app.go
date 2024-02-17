package kernel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"time"
)

type AppContract interface {
	GetAppFyne() fyne.App
	GetWindow() WindowContract
	GetQueue() QueueListContract
	GetLocalizerService() LocalizerContract
	GetConvertorService() ConvertorContract
	AfterClosing()
	RunConvertor()
}

type App struct {
	AppFyne fyne.App
	Window  WindowContract
	Queue   QueueListContract

	localizerService LocalizerContract
	convertorService ConvertorContract
}

func NewApp(
	metadata *fyne.AppMetadata,
	localizerService LocalizerContract,
	queue QueueListContract,
	queueLayoutObject QueueLayoutObjectContract,
	convertorService ConvertorContract,
) *App {
	app.SetMetadata(*metadata)
	a := app.New()

	return &App{
		AppFyne: a,
		Window:  newWindow(a.NewWindow("GUI for FFmpeg"), NewLayout(queueLayoutObject, localizerService)),
		Queue:   queue,

		localizerService: localizerService,
		convertorService: convertorService,
	}
}

func (a App) GetAppFyne() fyne.App {
	return a.AppFyne
}

func (a App) GetQueue() QueueListContract {
	return a.Queue
}

func (a App) GetWindow() WindowContract {
	return a.Window
}

func (a App) GetLocalizerService() LocalizerContract {
	return a.localizerService
}

func (a App) GetConvertorService() ConvertorContract {
	return a.convertorService
}

func (a App) AfterClosing() {
	for _, cmd := range a.convertorService.GetRunningProcesses() {
		_ = cmd.Process.Kill()
	}
}

func (a App) RunConvertor() {
	go func() {
		for {
			time.Sleep(time.Millisecond * 3000)
			queueId, queue := a.Queue.Next()
			if queue == nil {
				continue
			}
			queue.Status = StatusType(InProgress)
			a.Window.GetLayout().ChangeQueueStatus(queueId, queue)

			totalDuration, err := a.convertorService.GetTotalDuration(&queue.Setting.VideoFileInput)
			if err != nil {
				queue.Status = StatusType(Error)
				queue.Error = err
				a.Window.GetLayout().ChangeQueueStatus(queueId, queue)
				continue
			}
			progress := a.Window.GetLayout().NewProgressbar(queueId, totalDuration)

			err = a.convertorService.RunConvert(*queue.Setting, progress)
			if err != nil {
				queue.Status = StatusType(Error)
				queue.Error = err
				a.Window.GetLayout().ChangeQueueStatus(queueId, queue)
				continue
			}
			queue.Status = StatusType(Completed)
			a.Window.GetLayout().ChangeQueueStatus(queueId, queue)
		}
	}()
}
