package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb"
)

// Checkpointer allows us to restart a service from the last time we executed.
type TaskServiceCheckpointer struct {
	ts influxdb.TaskService
}

func NewTaskServiceCheckpointer(ts influxdb.TaskService) *TaskServiceCheckpointer {
	return &TaskServiceCheckpointer{
		ts: ts,
	}
}

// Checkpoint updates a task's LatestCompleted value with the given time
func (c *TaskServiceCheckpointer) Checkpoint(ctx context.Context, id influxdb.ID, t time.Time) error {
	_, err := c.ts.UpdateTask(ctx, id, influxdb.TaskUpdate{
		LatestCompleted: &t,
	})

	if err != nil {
		return fmt.Errorf("could not update checkpoint for task: %v", err)
	}
	return nil
}

// Last retrieves a task by its ID and returns its LatestCompleted value
func (c *TaskServiceCheckpointer) Last(ctx context.Context, id influxdb.ID) (time.Time, error) {
	task, err := c.ts.FindTaskByID(ctx, id)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not fetch task: %v", err)
	}

	return task.LatestCompleted, nil
}
