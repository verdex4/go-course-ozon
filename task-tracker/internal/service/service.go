package service

import (
	"fmt"
	"time"

	domain "github.com/verdex4/go-course-ozon/task-tracker/internal/domain"
	storage "github.com/verdex4/go-course-ozon/task-tracker/internal/storage"
)

type Service struct {
	store *storage.Storage
}

func NewService() *Service {
	store := storage.NewStorage()
	return &Service{store}
}

func (s *Service) AddTask(title string, deadline time.Time) error {
	if title == "" {
		return fmt.Errorf("Ошибка: Название задачи не может быть пустым.")
	}
	if deadline.Before(time.Now()) {
		return fmt.Errorf("Ошибка: дата '%s' уже прошла.\n", deadline)
	}
	s.store.Add(title, deadline)
	return nil
}

func (s *Service) ListTasks() []*domain.Task {
	return s.store.List()
}

func (s *Service) FindTask(id int) (*domain.Task, error) {
	return s.store.Find(id)
}

func (s *Service) ChangeStatus(id int, status string) error {
	if status != "planned" && status != "in_progress" && status != "cancelled" && status != "done" {
		return fmt.Errorf("Ошибка: неверный статус задачи '%s'. Доступны: planned, in_progress, cancelled, done\n", status)
	}
	cur, err := s.store.Find(id)
	if err != nil {
		return err
	}
	if cur.Status == status {
		return fmt.Errorf("Ошибка: задача с id %d уже имеет статус '%s'\n", id, status)
	}
	s.store.ChangeStatus(id, status)
	return nil
}

func (s *Service) ChangeTitle(id int, title string) error {
	if title == "" {
		return fmt.Errorf("Ошибка: Название задачи не может быть пустым.")
	}
	cur, err := s.store.Find(id)
	if err != nil {
		return err
	}
	if cur.Title == title {
		return fmt.Errorf("Ошибка: задача с id %d уже имеет название '%s'\n", id, title)
	}
	s.store.ChangeTitle(id, title)
	return nil
}

func (s *Service) ChangeDeadline(id int, deadline time.Time) error {
	if deadline.Before(time.Now()) {
		return fmt.Errorf("Ошибка: дата '%s' уже прошла.\n", deadline)
	}
	cur, err := s.store.Find(id)
	if err != nil {
		return err
	}
	if time.Time.Equal(cur.Deadline, deadline) {
		return fmt.Errorf("Ошибка: задача с id %d уже имеет дедлайн '%s'\n", id, deadline)
	}
	s.store.ChangeDeadline(id, deadline)
	return nil
}

func (s *Service) DeleteTask(id int) bool {
	prev := s.store.Length()
	s.store.Delete(id)
	if prev == s.store.Length() {
		return false
	}
	return true
}
