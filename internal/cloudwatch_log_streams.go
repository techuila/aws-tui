package internal

import (
	"time"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type CloudWatchLogStreams struct {
	*ui.Table
	view.CloudWatch
	repo         *repo.CloudWatch
	app          *Application
	model        []model.CloudWatchLogStream
	logGroupName string
}

func NewCloudWatchLogStreams(logGroupName string, repo *repo.CloudWatch, app *Application) *CloudWatchLogStreams {
	c := &CloudWatchLogStreams{
		Table: ui.NewTable([]string{
			"NAME",
			"FIRST EVENT",
			"LAST EVENT",
			"LAST INGESTION",
		}, 1, 0),
		repo:         repo,
		app:          app,
		logGroupName: logGroupName,
	}
	c.SetSelectedFunc(c.eventsHandler)
	return c
}

func (c CloudWatchLogStreams) GetLabels() []string {
	return []string{c.logGroupName, "Log Streams"}
}

func (c CloudWatchLogStreams) eventsHandler(row, col int) {
	logStreamName, err := c.GetColSelection("NAME")
	if err != nil {
		return
	}
	eventsView := NewCloudWatchLogEvents(c.logGroupName, logStreamName, c.repo, c.app)
	c.app.AddAndSwitch(eventsView)
}

func (c CloudWatchLogStreams) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func formatLogTimestamp(ms *int64) string {
	if ms == nil {
		return "-"
	}
	return time.UnixMilli(*ms).Local().Format("2006-01-02 15:04:05")
}

func (c *CloudWatchLogStreams) Render() {
	model, err := c.repo.ListLogStreams(c.logGroupName)
	if err != nil {
		panic(err)
	}
	c.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.LogStreamName, ""),
			formatLogTimestamp(v.FirstEventTimestamp),
			formatLogTimestamp(v.LastEventTimestamp),
			formatLogTimestamp(v.LastIngestionTime),
		})
	}
	c.SetData(data)
}
