package internal

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type CloudWatchLogEvents struct {
	*ui.Table
	view.CloudWatch
	repo          *repo.CloudWatch
	app           *Application
	model         []model.CloudWatchLogEvent
	logGroupName  string
	logStreamName string
	expanded      map[int]bool
	expandAll     bool
	// rowToEvent maps a table row index to the event index it belongs to.
	// The header row and continuation rows of an expanded event both point
	// to their owning event.
	rowToEvent []int
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
		expanded:      make(map[int]bool),
	}
	c.SetSelectedFunc(c.toggleHandler)
	return c
}

func (c CloudWatchLogEvents) GetLabels() []string {
	return []string{c.logStreamName, "Events"}
}

// prettifyMessage indents a log message when it is valid JSON, leaving plain
// text untouched.
func prettifyMessage(msg string) string {
	if json.Valid([]byte(msg)) {
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(msg), "", "  "); err == nil {
			return buf.String()
		}
	}
	return msg
}

func (c *CloudWatchLogEvents) toggleHandler(row, col int) {
	if row < 0 || row >= len(c.rowToEvent) {
		return
	}
	idx := c.rowToEvent[row]
	if idx < 0 {
		return
	}
	c.expanded[idx] = !c.expanded[idx]
	c.Render()
}

func (c *CloudWatchLogEvents) toggleExpandAll() {
	c.expandAll = !c.expandAll
	c.Render()
}

func (c CloudWatchLogEvents) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Toggle Prettify All",
			Action:      c.toggleExpandAll,
		},
	}
}

func (c *CloudWatchLogEvents) Render() {
	model, err := c.repo.GetLogEvents(c.logGroupName, c.logStreamName)
	if err != nil {
		panic(err)
	}
	c.model = model

	var data [][]string
	// row 0 is the header and belongs to no event
	c.rowToEvent = []int{-1}
	for i, v := range model {
		timestamp := formatLogTimestamp(v.Timestamp)
		message := utils.DerefString(v.Message, "")
		if c.expandAll || c.expanded[i] {
			pretty := strings.TrimRight(prettifyMessage(message), "\n")
			for j, line := range strings.Split(pretty, "\n") {
				ts := ""
				if j == 0 {
					ts = timestamp
				}
				data = append(data, []string{ts, line})
				c.rowToEvent = append(c.rowToEvent, i)
			}
		} else {
			// collapse any embedded newlines so the event stays on one row
			singleLine := strings.ReplaceAll(message, "\n", " ")
			data = append(data, []string{timestamp, singleLine})
			c.rowToEvent = append(c.rowToEvent, i)
		}
	}
	c.Reset()
	c.SetData(data)
}
