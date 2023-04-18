package main

import "fmt"

func main() {
	blackCell := true
	robot1Pos := -10
	robot2Pos := 10
	robot1Program := []string{"IF FLAG", "GOTO 4", "ML", "GOTO 1"} // программа первого робота
	robot2Program := []string{"IF FLAG", "GOTO 1", "MR", "GOTO 2"} // программа второго робота

	for !blackCell { // жду, пока роботы не найдут чёрную клетку
		if robot1Pos == robot2Pos { // если роботы встретились, то прерываю цикл
			break
		}

		for i := 0; i < 2; i++ { // выполняю по одной команде у каждого робота
			if i == 0 {
				command := robot1Program[robot1Pos]
				switch command {
				case "ML":
					robot1Pos++
				case "MR":
					robot1Pos--
				case "IF FLAG":
					if blackCell {
						robot1Pos++
					} else {
						robot1Pos += 2
					}
				case "GOTO N":
					n := int(command[len(command)-1] - '0') // парсим номер строки
					robot1Pos = n - 1                       // отнимаем 1, т.к. индексы в программе начинаются с 0
				}
			} else {
				command := robot2Program[robot2Pos]
				switch command {
				case "ML":
					robot2Pos++
				case "MR":
					robot2Pos--
				case "IF FLAG":
					if blackCell {
						robot2Pos += 2
					} else {
						robot2Pos++
					}
				case "GOTO N":
					n := int(command[len(command)-1] - '0')
					robot2Pos = n - 1
				}
			}
		}
		if robot1Pos == robot2Pos && robot1Pos%2 == 1 { // проверяю, нашли ли роботы чёрную клетку
			blackCell = true
		}
	}

	if blackCell {
		fmt.Println("Роботы встретились на чёрной клетке")
	} else {
		fmt.Println("Роботы не встретились")
	}
}
