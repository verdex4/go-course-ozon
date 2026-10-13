package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID       int
	Title    string
	Status   string // "planned", "in_progress", "cancelled", "done"
	Deadline time.Time
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	tasks := make([]*Task, 0)
	nextId := 1
	fmt.Println("Трекер задач запущен, введите команды")

	for {
		if scanner.Err() != nil {
			break
		}

		var tokens []string
		if scanner.Scan() {
			line := scanner.Text()
			tokens = strings.Fields(line)
			if len(tokens) == 0 {
				continue
			}
		}

		command := strings.ToLower(tokens[0])
		args := tokens[1:]

		switch command {
		case "help":
			fmt.Println("Доступные команды:")
			fmt.Println("add <title> <deadline> - добавляет новую задачу")
			fmt.Println("list - выводит список задач")
			fmt.Println("task <id> - выводит информацию о задаче")
			fmt.Println("new-status <id> <status> - меняет статус задачи")
			fmt.Println("change <id> title <title> deadline <deadline> - меняет название и дату задачи (указать достаточно один из аргументов, title или deadline)")

		case "exit":
			fmt.Println("Завершение работы...")
			return

		case "add":
			if len(args) != 3 {
				fmt.Println("Ошибка: команда add требует 2 аргумента: title и deadline в формате YYYY-MM-DD HH:MM:SS")
				continue
			}

			title := args[0]
			if title == "" {
				fmt.Println("Ошибка: название задачи не может быть пустым.")
				continue
			}
			deadline, err := time.Parse("2006-01-02 15:04:05", args[1]+" "+args[2])
			if err != nil {
				fmt.Printf("Ошибка: неверный формат даты '%s'.\n", args[1])
				continue
			}
			if deadline.Before(time.Now()) {
				fmt.Printf("Ошибка: дата '%s' уже прошла.\n", args[1])
				continue
			}

			task := Task{
				ID:       nextId,
				Title:    title,
				Status:   "planned",
				Deadline: deadline,
			}
			nextId++
			tasks = append(tasks, &task)

		case "list":
			if len(args) != 0 {
				fmt.Println("Ошибка: команда list не требует аргументов")
				continue
			}
			for _, task := range tasks {
				s := fmt.Sprintf("ID: %d; Title: %s; Status: %s; Deadline: %s", task.ID, task.Title, task.Status, task.Deadline.Format("2006-01-02 15:04:05"))
				fmt.Println(s)
			}

		case "task":
			if len(args) != 1 {
				fmt.Println("Ошибка: команда task требует 1 аргумент: id")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Ошибка: ожидалось целое число. Получено: '%s'\n", args[0])
				continue
			}
			for _, task := range tasks {
				if task.ID == id {
					s := fmt.Sprintf("ID: %d; Title: %s; Status: %s; Deadline: %s", task.ID, task.Title, task.Status, task.Deadline.Format("2006-01-02 15:04:05"))
					fmt.Println(s)
					break
				}
			}

		case "new-status":
			if len(args) != 2 {
				fmt.Println("Ошибка: команда new-status требует 2 аргумента: id и status")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Ошибка: ожидалось целое число. Получено: '%s'\n", args[0])
				continue
			}
			status := args[1]
			if status != "planned" && status != "in_progress" && status != "cancelled" && status != "done" {
				fmt.Printf("Ошибка: неверный статус задачи '%s'. Доступны: planned, in_progress, cancelled, done\n", status)
				continue
			}
			for _, task := range tasks {
				if task.ID == id {
					if status == task.Status {
						fmt.Printf("Ошибка: задача уже имеет статус '%s'\n", status)
						break
					}
					task.Status = status
					break
				}
			}
			fmt.Println("Статус задачи успешно изменен")

		case "change":
			if len(args) <= 2 {
				fmt.Println("Ошибка: команда change требует хотя бы 3 аргумента: id, title <title> или deadline <deadline>")
				continue
			}
			if len(args) > 6 {
				fmt.Println("Ошибка: введено слишком много аргументов. Ожидалось: change <id> title <title> deadline <deadline> (наличие обоих title и deadline необязательно)")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Ошибка: ожидалось целое число для первого аргумента ID. Получено: '%s'\n", args[0])
				continue
			}
			arg := args[1]
			if arg != "title" && arg != "deadline" {
				fmt.Printf("Ошибка: неверная команда '%s'. Доступные команды: title, deadline\n", arg)
				continue
			}
			index := 0
			if arg == "title" {
				title := args[2]
				if title == "" {
					fmt.Println("Ошибка: название задачи не может быть пустым.")
					continue
				}
				for _, task := range tasks {
					if task.ID == id {
						if title == task.Title {
							fmt.Printf("Ошибка: задача уже имеет название '%s'\n", title)
							break
						}
						task.Title = title
						break
					}
				}
				index = 3
			} else {
				d := args[2] + " " + args[3]
				deadline, err := time.Parse("2006-01-02 15:04:05", d)
				if err != nil {
					fmt.Printf("Ошибка: неверный формат даты '%s'. Ожидалось: YYYY-MM-DD HH:MM:SS\n", d)
					continue
				}
				if deadline.Before(time.Now()) {
					fmt.Printf("Ошибка: дата '%s' уже прошла.\n", d)
					continue
				}
				for _, task := range tasks {
					if task.ID == id {
						if time.Time.Equal(deadline, task.Deadline) {
							fmt.Printf("Ошибка: задача уже имеет дедлайн '%s'\n", deadline)
							break
						}
						task.Deadline = deadline
						break
					}
				}
				index = 4
			}
			if index == len(args) {
				fmt.Printf("%s успешно изменен\n", args[1])
				continue
			}

			if len(args) == index+1 {
				fmt.Println("Ошибка: неверный формат команды. Ожидалось: change <id> title <title> deadline <deadline> (наличие обоих title и deadline необязательно)")
				fmt.Println("Обработана только первая команда")
			}
			arg = args[index]
			if arg != "title" && arg != "deadline" {
				fmt.Printf("Ошибка: неверная команда '%s'. Доступные команды: title, deadline\n", arg)
				continue
			}
			if arg == "title" {
				title := args[index+1]
				if title == "" {
					fmt.Println("Ошибка: название задачи не может быть пустым.")
					continue
				}
				for _, task := range tasks {
					if task.ID == id {
						if title == task.Title {
							fmt.Printf("Ошибка: задача уже имеет название '%s'\n", title)
							break
						}
						task.Title = title
						break
					}
				}
			} else {
				d := args[index+1] + " " + args[index+2]
				deadline, err := time.Parse("2006-01-02 15:04:05", d)
				if err != nil {
					fmt.Printf("Ошибка: неверный формат даты '%s'. Ожидалось: YYYY-MM-DD HH:MM:SS\n", d)
					continue
				}
				if deadline.Before(time.Now()) {
					fmt.Printf("Ошибка: дата '%s' уже прошла.\n", d)
					continue
				}
				for _, task := range tasks {
					if task.ID == id {
						if time.Time.Equal(deadline, task.Deadline) {
							fmt.Printf("Ошибка: задача уже имеет дедлайн '%s'\n", deadline)
							break
						}
						task.Deadline = deadline
						break
					}
				}
			}
			fmt.Println("Название и дедлайн задачи успешно изменены")

		case "delete":
			if len(args) != 1 {
				fmt.Println("Ошибка: неверное количество аргументов. Ожидался 1 аргумент: ID")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Ошибка: ожидалось целое число для аргумента ID. Получено: '%s'\n", args[0])
				continue
			}
			res := make([]*Task, len(tasks))
			for _, task := range tasks {
				if task.ID == id {
					fmt.Printf("Задача с ID %d удалена\n", id)
					continue
				}
				res = append(res, task)
			}
			if len(res) == len(tasks) {
				fmt.Printf("Элемента с ID %d не существует. Удаление не было выполнено\n", id)
				continue
			}
			tasks = res

		default:
			fmt.Printf("Неизвестная команда: '%s'. Наберите 'help' для справки.\n", command)
		}
	}
}
