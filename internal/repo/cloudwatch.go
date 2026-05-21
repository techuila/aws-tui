package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cwLogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwLogsTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/bporter816/aws-tui/internal/model"
)

type CloudWatch struct {
	cwLogsClient *cwLogs.Client
}

func NewCloudWatch(cwLogsClient *cwLogs.Client) *CloudWatch {
	return &CloudWatch{
		cwLogsClient: cwLogsClient,
	}
}

func (c CloudWatch) ListLogGroups() ([]model.CloudWatchLogGroup, error) {
	pg := cwLogs.NewDescribeLogGroupsPaginator(
		c.cwLogsClient,
		&cwLogs.DescribeLogGroupsInput{},
	)
	var logGroups []model.CloudWatchLogGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.CloudWatchLogGroup{}, err
		}
		for _, v := range out.LogGroups {
			logGroups = append(logGroups, model.CloudWatchLogGroup(v))
		}
	}
	return logGroups, nil
}

func (c CloudWatch) ListLogStreams(logGroupName string) ([]model.CloudWatchLogStream, error) {
	pg := cwLogs.NewDescribeLogStreamsPaginator(
		c.cwLogsClient,
		&cwLogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(logGroupName),
			OrderBy:      cwLogsTypes.OrderByLastEventTime,
			Descending:   aws.Bool(true),
		},
	)
	var logStreams []model.CloudWatchLogStream
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.CloudWatchLogStream{}, err
		}
		for _, v := range out.LogStreams {
			logStreams = append(logStreams, model.CloudWatchLogStream(v))
		}
	}
	return logStreams, nil
}

func (c CloudWatch) GetLogEvents(logGroupName string, logStreamName string) ([]model.CloudWatchLogEvent, error) {
	out, err := c.cwLogsClient.GetLogEvents(
		context.TODO(),
		&cwLogs.GetLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
			Limit:         aws.Int32(1000),
			StartFromHead: aws.Bool(true),
		},
	)
	if err != nil {
		return []model.CloudWatchLogEvent{}, err
	}
	var events []model.CloudWatchLogEvent
	for _, v := range out.Events {
		events = append(events, model.CloudWatchLogEvent(v))
	}
	return events, nil
}

func (c CloudWatch) ListTags(resourceArn string) (model.Tags, error) {
	out, err := c.cwLogsClient.ListTagsForResource(
		context.TODO(),
		&cwLogs.ListTagsForResourceInput{
			ResourceArn: aws.String(resourceArn),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for k, v := range out.Tags {
		tags = append(tags, model.Tag{Key: k, Value: v})
	}
	return tags, nil
}
