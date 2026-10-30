package internal

import (
	"fmt"
	"time"

	"github.com/verdex4/go-course-ozon/task-tracker/internal/domain"
)

type Storage struct {
	data   map[int]*domain.Task
	nextID int
}

func NewStorage() *Storage {
	d := make(map[int]*domain.Task)
	return &Storage{data: d, nextID: 1}
}

func (s *Storage) Length() int {
	return len(s.data)
}

func (s *Storage) Add(title string, deadline time.Time) {
	s.data[s.nextID] = &domain.Task{ID: s.nextID, Title: title, Status: "planned", Deadline: deadline}
	s.nextID++
}

func (s *Storage) List() []*domain.Task {
	res := make([]*domain.Task, 0, len(s.data))
	for _, task := range s.data {
		res = append(res, task)
	}
	return res
}

func (s *Storage) Find(id int) (*domain.Task, error) {
	t, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("Задача с id %d не найдена", id)
	}
	return t, nil
}

func (s *Storage) ChangeStatus(id int, status string) {
	s.data[id].Status = status
}

func (s *Storage) ChangeTitle(id int, title string) {
	s.data[id].Title = title
}

func (s *Storage) ChangeDeadline(id int, deadline time.Time) {
	s.data[id].Deadline = deadline
}

func (s *Storage) Delete(id int) {
	delete(s.data, id)
}
