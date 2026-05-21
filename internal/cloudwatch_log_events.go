package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type CloudWatchLogEvents struct {
	*ui.Table
	view.CloudWatch
	repo          *repo.CloudWatch
	app           *Application
	model         []model.CloudWatchLogEvent
	logGroupName  string
	logStreamName string
}

func NewCloudWatchLogEvents(logGroupName string, logStreamName string, repo *repo.CloudWatch, app *Application) *CloudWatchLogEvents {
	c := &CloudWatchLogEvents{
		Table: ui.NewTable([]string{
			"TIMESTAMP",
			"MESSAGE",
		}, 1, 0),
		repo:          repo,
		app:           app,
		logGroupName:  logGroupName,
		logStreamName: logStreamName,
	}
	return c
}

func (c CloudWatchLogEvents) GetLabels() []string {
	return []string{c.logStreamName, "Events"}
}

func (c CloudWatchLogEvents) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c *CloudWatchLogEvents) Render() {
	model, err := c.repo.GetLogEvents(c.logGroupName, c.logStreamName)
	if err != nil {
		panic(err)
	}
	c.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			formatLogTimestamp(v.Timestamp),
			utils.DerefString(v.Message, ""),
		})
	}
	c.SetData(data)
}
