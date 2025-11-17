package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func NewSystem() *System {
	return &System{
		users:     make(map[int]*User),
		transfers: []Transfer{},
		nextID:    1,
	}
}

func (s *System) findUserByName(name string) *User {
	for _, u := range s.users {
		if u.Name == name {
			return u
		}
	}
	return nil
}

func (s *System) Register(name, password string, balance int) error {
	if name == "" || password == "" {
		return errors.New("имя и пароль не могут быть пустыми")
	}
	if balance < 0 {
		return errors.New("баланс не может быть отрицательным")
	}
	if s.findUserByName(name) != nil {
		return errors.New("пользователь с таким именем уже существует")
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
	user := s.findUserByName(name)
	if user == nil {
		return errors.New("пользователь не найден")
	}
	if user.Password != password {
		return errors.New("неверный пароль")
	}
	s.currentUser = user
	return nil
}

func (s *System) Logout() {
	s.currentUser = nil
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

	toUser := s.findUserByName(toName)
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

func (s *System) PrintAllUsers() {
	if len(s.users) == 0 {
		fmt.Println("Нет пользователей.")
		return
	}
	fmt.Println("\nДоступные пользователи:")
	for _, u := range s.users {
		if s.currentUser != nil && u.ID == s.currentUser.ID {
			continue
		}
		fmt.Printf("  %s (баланс: %d)\n", u.Name, u.Balance)
	}
}

func (s *System) PrintUserTransfers(userName string) {
	var userTransfers []Transfer
	for _, t := range s.transfers {
		if t.FromName == userName || t.ToName == userName {
			userTransfers = append(userTransfers, t)
		}
	}

	if len(userTransfers) == 0 {
		fmt.Println("У вас пока нет переводов.")
		return
	}

	fmt.Println("\nВаши переводы:")
	for _, t := range userTransfers {
		if t.FromName == userName {
			fmt.Printf("  Вы -> %s: %d\n", t.ToName, t.Amount)
		} else {
			fmt.Printf("  %s -> Вам: %d\n", t.FromName, t.Amount)
		}
	}
}

func (s *System) PrintBalances() {
	fmt.Printf("\nВаш баланс: %d\n", s.currentUser.Balance)
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
	sys := NewSystem()
	reader := bufio.NewReader(os.Stdin)
	end := false
	for {
		fmt.Println("\n--- Система переводов ---")
		if sys.currentUser != nil {
			fmt.Printf("Вы вошли как: %s\n", sys.currentUser.Name)
		} else {
			fmt.Println("Вы не вошли в аккаунт.")
		}

		if sys.currentUser != nil {
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

		if sys.currentUser != nil {
			switch choice {
			case "1":
				sys.PrintAllUsers()
				if len(sys.users) <= 1 {
					fmt.Println("Некому переводить.")
					continue
				}
				toName := readLine("Кому перевести? (имя): ")
				amtStr := readLine("Сумма (целое число): ")
				amount, err := parseInt(amtStr)
				if err != nil {
					fmt.Println("Некорректная сумма (должно быть целое число)")
					continue
				}
				if err := sys.Transfer(toName, amount); err != nil {
					fmt.Printf("Ошибка: %v\n", err)
				} else {
					fmt.Println("Перевод выполнен!")
				}

			case "2":
				sys.PrintBalances()

			case "3":
				sys.PrintUserTransfers(sys.currentUser.Name)

			case "4":
				sys.Logout()
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
				if err := sys.Register(name, pass, bal); err != nil {
					fmt.Print("Ошибка: ", err)
				} else {
					fmt.Println("Регистрация успешна!")
				}

			case "2":
				name := readLine("Имя: ")
				pass := readLine("Пароль: ")
				if err := sys.Login(name, pass); err != nil {
					fmt.Print("Ошибка: ", err)
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
