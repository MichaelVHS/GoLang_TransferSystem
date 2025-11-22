package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)


type UserService interface {
	Register(name, password string, balance int) error
	Login(name, password string) error
	Logout()
	GetCurrentUser() *User
}

type TransferService interface {
	Transfer(toName string, amount int) error
	GetTransfersForUser(userName string) []Transfer
	PrintAllOtherUsers(currentUserName string)
}

type User struct {
	ID       int
	Name     string
	Password string
	Balance  int
}

type Transfer struct {
	FromName string
	ToName   string
	Amount   int
}

type System struct {
	users       map[int]*User
	transfers   []Transfer
	currentUser *User
	nextID      int
}


func (s *System) Register(name, password string, balance int) error {
	if name == "" || password == "" {
		return errors.New("имя и пароль не могут быть пустыми")
	}
	if balance < 0 {
		return errors.New("баланс не может быть отрицательным")
	}
	for _, u := range s.users {
		if u.Name == name {
			return errors.New("пользователь с таким именем уже существует")
		}
	}

	s.users[s.nextID] = &User{
		ID:       s.nextID,
		Name:     name,
		Password: password,
		Balance:  balance,
	}
	s.nextID++
	return nil
}

func (s *System) Login(name, password string) error {
	for _, u := range s.users {
		if u.Name == name {
			if u.Password == password {
				s.currentUser = u
				return nil
			}
			return errors.New("неверный пароль")
		}
	}
	return errors.New("пользователь не найден")
}

func (s *System) Logout() {
	s.currentUser = nil
}

func (s *System) GetCurrentUser() *User {
	return s.currentUser
}


func (s *System) Transfer(toName string, amount int) error {
	if s.currentUser == nil {
		return errors.New("сначала войдите в аккаунт")
	}
	if amount <= 0 {
		return errors.New("сумма должна быть больше нуля")
	}
	if s.currentUser.Name == toName {
		return errors.New("нельзя переводить самому себе")
	}

	var toUser *User
	for _, u := range s.users {
		if u.Name == toName {
			toUser = u
			break
		}
	}
	if toUser == nil {
		return errors.New("получатель не найден")
	}
	if s.currentUser.Balance < amount {
		return errors.New("недостаточно средств")
	}

	s.currentUser.Balance -= amount
	toUser.Balance += amount

	s.transfers = append(s.transfers, Transfer{
		FromName: s.currentUser.Name,
		ToName:   toName,
		Amount:   amount,
	})
	return nil
}

func (s *System) GetTransfersForUser(userName string) []Transfer {
	var result []Transfer
	for _, t := range s.transfers {
		if t.FromName == userName || t.ToName == userName {
			result = append(result, t)
		}
	}
	return result
}

func (s *System) PrintAllOtherUsers(currentUserName string) {
	fmt.Println("\nДоступные пользователи:")
	otherExists := false
	for _, u := range s.users {
		if u.Name != currentUserName {
			fmt.Printf("  %s (баланс: %d)\n", u.Name, u.Balance)
			otherExists = true
		}
	}
	if !otherExists {
		fmt.Println("Нет других пользователей.")
	}
}


func readLine(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}


func main() {
	sys := &System{
		users:     make(map[int]*User),
		transfers: []Transfer{},
		nextID:    1,
	}

	var userService UserService = sys
	var transferService TransferService = sys

	reader := bufio.NewReader(os.Stdin)
	end := false

	for {
		fmt.Println("\n--- Система переводов ---")
		currentUser := sys.GetCurrentUser()
		if currentUser != nil {
			fmt.Printf("Вы вошли как: %s\n", currentUser.Name)
		} else {
			fmt.Println("Вы не вошли в аккаунт.")
		}

		if currentUser != nil {
			fmt.Print(`1. Перевести деньги
2. Мой баланс
3. История переводов
4. Выйти из аккаунта
5. Выйти из программы
Выберите (1-5):`)
		} else {
			fmt.Print(`1. Зарегистрироваться
2. Войти
3. Выйти
Выберите (1-3): `)
		}

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if currentUser != nil {
			switch choice {
			case "1":
				transferService.PrintAllOtherUsers(currentUser.Name)
				toName := readLine("Кому перевести? (имя): ")
				amtStr := readLine("Сумма (целое число): ")
				amount, err := parseInt(amtStr)
				if err != nil {
					fmt.Println("Некорректная сумма (должно быть целое число)")
					continue
				}
				if err := transferService.Transfer(toName, amount); err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					fmt.Println("Перевод выполнен!")
				}

			case "2":
				fmt.Printf("\nВаш баланс: %d\n", currentUser.Balance)

			case "3":
				transfers := transferService.GetTransfersForUser(currentUser.Name)
				if len(transfers) == 0 {
					fmt.Println("У вас пока нет переводов.")
				} else {
					fmt.Println("\nВаши переводы:")
					for _, t := range transfers {
						if t.FromName == currentUser.Name {
							fmt.Printf("  Вы -> %s: %d\n", t.ToName, t.Amount)
						} else {
							fmt.Printf("  %s -> Вам: %d\n", t.FromName, t.Amount)
						}
					}
				}

			case "4":
				userService.Logout()
				fmt.Println("Вы вышли из аккаунта.")

			case "5":
				fmt.Println("До свидания!")
				end = true

			default:
				fmt.Println("Неверный выбор.")
			}
		} else {
			switch choice {
			case "1":
				name := readLine("Ваше имя: ")
				pass := readLine("Пароль: ")
				balStr := readLine("Начальный баланс (целое число): ")
				bal, err := parseInt(balStr)
				if err != nil {
					fmt.Println("Некорректный баланс (должен быть целым числом)")
					continue
				}
				if err := userService.Register(name, pass, bal); err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Регистрация успешна!")
				}

			case "2":
				name := readLine("Имя: ")
				pass := readLine("Пароль: ")
				if err := userService.Login(name, pass); err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Добро пожаловать!")
				}

			case "3":
				fmt.Println("До свидания!")
				end = true

			default:
				fmt.Println("Неверный выбор.")
			}
		}

		if end {
			break
		}
	}
}
