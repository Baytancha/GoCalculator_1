package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//var Syntax_Error = fmt.Errorf("Строка не является математической операцией")

type InvalidSymbol struct{}
type NegativeRoman struct{}
type MixedNumericSystems struct{}
type InvalidExpression struct{}

func (e InvalidSymbol) Error() string {
	return fmt.Sprintln("Error: введен некорректный символ")
}

func (e NegativeRoman) Error() string {
	return fmt.Sprintln("Error: в римской системе нет нуля и отрицательных чисел")
}

func (e MixedNumericSystems) Error() string {
	return fmt.Sprintln("Error: используются разные системы счисления")
}

func (e InvalidExpression) Error() string {
	return fmt.Sprintln("Error: Строка не является математической операцией")
}

func ArabicToRoman(num int) (result string) {

	charset := map[int]string{
		1:   "I",
		5:   "V",
		10:  "X",
		50:  "L",
		100: "C",
	}
	//5 это обшее количество слагаемых: 1,5,10,50,100 (с числами больше 100 не работаем по условиям задачи)
	for i, divisor := 0, 100; i < 5; i++ {

		if num == 0 {
			break
		}

		if num/1 == 9 { // 9
			result = result + charset[1]
			result = result + charset[10]
			num -= 9
		} else if num/10 == 9 { // 90...99
			result = result + charset[10]
			result = result + charset[100]
			num -= 90
		} else if num/100 == 9 { // 900...999 (не актуально)
			result = result + charset[100]
			result = result + charset[1000]
			num -= 900
		} else if num/1 == 4 { //4
			result = result + charset[1]
			result = result + charset[5]
			num -= 4
		} else if num/10 == 4 { /// 40...49
			result = result + charset[10]
			result = result + charset[50]
			num -= 40
		} else if num/100 == 4 { //400...499 (не актуально)
			result = result + charset[100]
			result = result + charset[500]
			num -= 400
		} else {
			for num-divisor >= 0 { //разложение числа на слагаемые конкретного разряда
				result = result + charset[divisor]
				num -= divisor
			}
		}

		if i%2 == 0 { //определение нового слагаемого
			divisor /= 2
		} else {
			divisor /= 5
		}
	}
	return result
}

func RomanToArabic(rom string) (x int, e error) {
	switch rom {
	case "I":
		return 1, nil
	case "II":
		return 2, nil
	case "III":
		return 3, nil
	case "IV":
		return 4, nil
	case "V":
		return 5, nil
	case "VI":
		return 6, nil
	case "VII":
		return 7, nil
	case "VIII":
		return 8, nil
	case "IX":
		return 9, nil
	case "X":
		return 10, nil
	default:
		return 0, InvalidSymbol{}
	}

}

// производится конвертация из римской в арабску, вычисления а также обработка ошибок
func Calculate(num1 string, num2 string, op string) (res int, err error) {

	//конвертация в массив unicode знаков и передача первого элемента массива
	if !unicode.IsDigit([]rune(num1)[0]) && unicode.IsDigit([]rune(num2)[0]) {
		return 0, MixedNumericSystems{}
	}

	if unicode.IsDigit([]rune(num1)[0]) && !unicode.IsDigit([]rune(num2)[0]) {
		return 0, MixedNumericSystems{}
	}

	if !unicode.IsDigit([]rune(num1)[0]) && !unicode.IsDigit([]rune(num2)[0]) {

		//проверка на валидность символов
		val1, err := RomanToArabic(num1)
		if err != nil {
			fmt.Println(err)
			panic(err) //не  перебрасываем в место вызова а прекращаем тут же
		}

		val2, err := RomanToArabic(num2)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		switch string(op) {
		case "+":
			res = add(val1, val2)
		case "-":
			if res = subtract(val1, val2); res <= 0 { //запрет вычитания отриц чисел
				return 0, NegativeRoman{}
			}
		case "*":
			res = multiply(val1, val2)
		case "/":
			res = divide(val1, val2)
		default:
			err = InvalidSymbol{} //критическая ошибка: введен неверный символ
			fmt.Println(err)
			panic(err)
		}
	} else {

		// конвертация str в int
		val1, _ := strconv.Atoi(num1)
		val2, _ := strconv.Atoi(num2)

		//fmt.Println(val1, val2)
		switch string(op) {
		case "+":
			res = add(val1, val2)
		case "-":
			res = subtract(val1, val2)
		case "*":
			res = multiply(val1, val2)
		case "/":
			res = divide(val1, val2)
		default:
			err = InvalidSymbol{}
			fmt.Println(err)
			panic(err)
		}

	}
	return res, nil
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) int {
	if b != 0 {
		return a / b
	}
	return 0
}

func ValidateInput(num1, num2, op *string, expression string) (err error) {
	//создаем массив из строк, убирая пробелы
	oper_count := strings.Fields(expression)
	//если в слайсе не три элемента, то выражение не валидно
	if len(oper_count) != 3 {
		return InvalidExpression{}
	}
	*num1 = oper_count[0]
	*op = oper_count[1]
	*num2 = oper_count[2]
	return nil

}

func main() {

	fmt.Println("Нажми '!' чтобы прекратить программу")
	for {

		fmt.Println("Input: ")

		var expression string
		var num1, num2, operation string

		//bufio для ввода строки с пробелами
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			expression = scanner.Text()
			//fmt.Println(expression)
		}

		if expression == "!" { //программа завершается если за вводом знака следует \r\n
			fmt.Println("Завершение программы")
			break
		}

		err := ValidateInput(&num1, &num2, &operation, expression)
		if err != nil {
			panic(err)
		}

		result, err := Calculate(num1, num2, operation)
		if err != nil {
			panic(err)
		} else {
			//если num1 и num2 содержат римские символы (иных после проверки содержать не могут)
			//конвертация строки в массив unicode знаков
			if !unicode.IsDigit([]rune(num1)[0]) && !unicode.IsDigit([]rune(num2)[0]) {
				fmt.Println("Output:")
				fmt.Println(ArabicToRoman(result))
				fmt.Println()
			} else {
				fmt.Println("Output:")
				fmt.Println(result)
				fmt.Println()
			}
		}
	}

}
