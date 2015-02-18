package basecrm

import (
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	Id           int          `json:"id,omitempty"`
	CreatorId    int          `json:"creator_id,omitempty"`
	OwnerId      int          `json:"owner_id,omitempty"`
	ResourceType ResourceType `json:"resource_type,omitempty"`
	ResourceId   int          `json:"resource_id,omitempty"`
	Completed    bool         `json:"completed,omitempty"`
	CompletedAt  time.Time    `json:"completed_at,omitempty"`
	DueDate      time.Time    `json:"due_date,omitempty"`
	Overdue      bool         `json:"overdue,omitempty"`
	Remind       bool         `json:"remind,omitempty"`
	RemindAt     time.Time    `json:"remind_at,omitempty"`
	Content      string       `json:"content,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
}

type TaskListOptions struct {
	Q string `url:"q,omitempty"`

	CreatorId int `url:"creator_id,omitempty"`
	OwnerId   int `url:"owner_id,omitempty"`

	ResourceType ResourceType `url:"resource_type,omitempty"`
	ResourceId   int          `url:"resource_id,omitempty"`

	Completed bool `url:"completed,omitempty"`
	Overdue   bool `url:"overdue,omitempty"`
	Remind    bool `url:"remind,omitempty"`

	ListOptions
}

type TasksService interface {
	List(opt *TaskListOptions) ([]*Task, *Response, error)
	Get(id int) (*Task, *Response, error)
	Create(task *Task) (*Task, *Response, error)
	Edit(id int, task *Task) (*Task, *Response, error)
	Delete(id int) (bool, *Response, error)
}

func NewTasksService(client *Client) TasksService {
	return &TasksServiceOp{client}
}

type taskRoot struct {
	Task *Task `json:"data"`
	Meta *struct {
		Type string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type tasksRoot struct {
	Items []*taskRoot `json:"items"`
	Meta  *Meta       `json:"meta"`
}

func (r *tasksRoot) Tasks() []*Task {
	tasks := make([]*Task, len(r.Items))
	for i, root := range r.Items {
		tasks[i] = root.Task
	}
	return tasks
}

type TasksServiceOp struct {
	client *Client
}

func (s *TasksServiceOp) List(opt *TaskListOptions) ([]*Task, *Response, error) {
	u, err := addOptions("/v2/tasks", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tasksRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Tasks(), res, err
}

func (s *TasksServiceOp) Get(id int) (*Task, *Response, error) {
	u := fmt.Sprintf("/v2/tasks/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(taskRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Task, res, err
}

func (s *TasksServiceOp) Create(task *Task) (*Task, *Response, error) {
	u := "/v2/tasks"
	envelope := &taskRoot{Task: task}
	req, err := s.client.NewRequest("POST", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(taskRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Task, res, err
}

func (s *TasksServiceOp) Edit(id int, task *Task) (*Task, *Response, error) {
	u := fmt.Sprintf("/v2/tasks/%d", id)
	envelope := &taskRoot{Task: task}
	req, err := s.client.NewRequest("PUT", u, envelope)
	if err != nil {
		return nil, nil, err
	}

	root := new(taskRoot)
	res, err := s.client.Do(req, root)
	if err != nil {
		return nil, res, err
	}

	return root.Task, res, err
}

func (s *TasksServiceOp) Delete(id int) (bool, *Response, error) {
	u := fmt.Sprintf("/v2/tasks/%d", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return false, nil, err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return false, res, err
	}

	return res.StatusCode == http.StatusNoContent, res, err
}
