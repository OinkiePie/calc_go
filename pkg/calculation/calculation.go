package main

import (
	"errors"
	// "fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	// Удаляем пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	// Преобразуем выражение в постфиксное
	postfixExpression, err := infixToPostfix(expression)
	if err != nil {
		return 0, err
	}

	// Вычисляем значение постфиксного выражения
	result, err := evaluatePostfix(postfixExpression)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Преобразование из инфиксного выражения в постфиксное
func infixToPostfix(expression string) (string, error) {
	stack := []rune{} 
	postfix := "" 
	for _, char := range expression { // Проходим по всем символам в инфиниксном выражении
		if isOperand(char) { // Если символ - операнд 
			postfix += string(char) // Добавляем его в постфиксное выражение
		} else if char == '(' { // Если символ - открывающая скобка
			stack = append(stack, char) // Добавляем её в стек
		} else if char == ')' { // Если символ - закрывающая скобка
			for len(stack) > 0 && stack[len(stack)-1] != '(' { // Переносим операторы из стека в постфиксное выражение, пока не встретим открывающую скобку

				postfix += string(stack[len(stack)-1]) 
				stack = stack[:len(stack)-1] 
			}
			if len(stack) == 0 { // Если стек пуст, то скобки не сбалансированы
				return "", errors.New("Неправильное количество скобок")
			}
			stack = stack[:len(stack)-1] // Удаляем открывающую скобку из стека
		} else { // Если символ - оператор
			for len(stack) > 0 && precedence(char) <= precedence(stack[len(stack)-1]) { // Переносим операторы с более высоким или равным приоритетом из стека в постфиксное выражение
				postfix += string(stack[len(stack)-1]) 
				stack = stack[:len(stack)-1] 
			}
			stack = append(stack, char)
		}
	}
	for len(stack) > 0 { // Переносим оставшиеся операторы из стека в постфикс
		postfix += string(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return postfix, nil // Возвращаем постфиксное выражение
}

// Вычисление постфиксного выражения
func evaluatePostfix(expression string) (float64, error) {
	stack := []float64{}
	for _, char := range expression { //Проходим по каждому символу в постфиксном выражении
		if isOperand(char) { // 3. Если символ - число
			num, err := strconv.ParseFloat(string(char), 64) 
			if err != nil { 
				return 0, err
			}
			stack = append(stack, num)
		} else { // 7. Если символ - оператор
			if len(stack) < 2 { // Проверяем, достаточно ли операндов для операции
				return 0, errors.New("Недостаточно операндов для операции")
			}
			operand2 := stack[len(stack)-1] // Извлекаем второй операнд с конца стека
			operand1 := stack[len(stack)-2] // Извлекаем первый операнд 
			stack = stack[:len(stack)-2] // Удаляем оба операнда из стека
			result := calculate(operand1, operand2, char) // Вычисляем результат операции
			stack = append(stack, result) // Добавляем результат в стек
		}
	}
	if len(stack) != 1 { // В стеке должно остаться ровно одно число - результат выражения
		return 0, errors.New("Неверное количество операндов")
	}
	return stack[0], nil 
}

// Вычисление операции
func calculate(operand1, operand2 float64, operator rune) float64 {
	switch operator {
	case '+':
		return operand1 + operand2
	case '-':
		return operand1 - operand2
	case '*':
		return operand1 * operand2
	case '/':
		if operand2 == 0 {
			return 0
		}
		return operand1 / operand2
	default:
		return 0
	}
}

// Проверка, является ли символ операндом
func isOperand(char rune) bool {
	return char >= '0' && char <= '9'
}

// Проверка, является ли символ буквой
func isLetter(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}

// Приоритет оператора
func precedence(operator rune) int {
	switch operator {
	case '(', ')':
		return 0
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return -1
	}
}

// func main() {
// 	expression := "2 * (5 - 2) + 4 * 7 - (2 + 3)"
// 	result, err := Calc(expression)
// 	if err != nil {
// 		fmt.Println("Ошибка:", err)
// 	} else {
// 		fmt.Println("Результат:", result)
// 	}
// }