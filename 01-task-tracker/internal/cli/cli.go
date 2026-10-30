package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	service "github.com/verdex4/go-course-ozon/task-tracker/internal/service"
)

func Run() {
	service := service.NewService()
	scanner := bufio.NewScanner(os.Stdin)
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
			fmt.Println("delete <id> - удаляет задачу")

		case "exit":
			fmt.Println("Завершение работы...")
			return

		case "add":
			if len(args) != 3 {
				fmt.Println("Ошибка: команда add требует 2 аргумента: title и deadline в формате YYYY-MM-DD HH:MM:SS")
				continue
			}
			title := args[0]
			deadline, err := time.Parse("2006-01-02 15:04:05", args[1]+" "+args[2])
			if err != nil {
				fmt.Printf("Ошибка: неверный формат даты '%s'.\n", args[1]+" "+args[2])
				continue
			}
			err = service.AddTask(title, deadline)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Задача добавлена")

		case "list":
			if len(args) != 0 {
				fmt.Println("Ошибка: команда list не требует аргументов")
				continue
			}
			lst := service.ListTasks()
			if len(lst) == 0 {
				fmt.Println("Список задач пуст.")
				continue
			}
			for _, task := range lst {
				fmt.Printf("ID: %d; Title: %s; Status: %s; Deadline: %s\n", task.ID, task.Title, task.Status, task.Deadline.Format("2006-01-02 15:04:05"))
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
			task, err := service.FindTask(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("ID: %d; Title: %s; Status: %s; Deadline: %s\n", task.ID, task.Title, task.Status, task.Deadline.Format("2006-01-02 15:04:05"))

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
			err = service.ChangeStatus(id, status)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Статус задачи успешно изменен на '%s'\n", status)

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
				service.ChangeTitle(id, title)
				index = 3
			} else {
				d := args[2] + " " + args[3]
				deadline, err := time.Parse("2006-01-02 15:04:05", d)
				if err != nil {
					fmt.Printf("Ошибка: неверный формат даты '%s'. Ожидалось: YYYY-MM-DD HH:MM:SS\n", d)
					continue
				}
				service.ChangeDeadline(id, deadline)
				index = 4
			}
			if index == len(args) {
				fmt.Printf("%s успешно изменен\n", args[1])
				continue
			}

			if len(args) == index+1 {
				fmt.Println("Ошибка: неверный формат второй команды. Отсутствует значение для изменения.")
			}
			secondArg := args[index]
			if arg == secondArg {
				fmt.Printf("Ошибка: вторая команда '%s' повторяет первую.\n", secondArg)
				continue
			}
			if secondArg != "title" && secondArg != "deadline" {
				fmt.Printf("Ошибка: неверная вторая команда: '%s'. Доступные команды: title, deadline\n", secondArg)
				continue
			}
			if secondArg == "title" {
				title := args[index+1]
				service.ChangeTitle(id, title)
			} else {
				d := args[index+1] + " " + args[index+2]
				deadline, err := time.Parse("2006-01-02 15:04:05", d)
				if err != nil {
					fmt.Printf("Ошибка: неверный формат даты '%s'. Ожидалось: YYYY-MM-DD HH:MM:SS\n", d)
					continue
				}
				service.ChangeDeadline(id, deadline)
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
			deleted := service.DeleteTask(id)
			if deleted {
				fmt.Printf("Задача с ID %d успешно удалена\n", id)
			} else {
				fmt.Printf("Задача с ID %d не найдена и не была удалена\n", id)
			}

		default:
			fmt.Printf("Неизвестная команда: '%s'. Наберите 'help' для справки.\n", command)
		}
	}
}
