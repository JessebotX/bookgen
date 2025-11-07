package html

import (
	"errors"
	"html/template"
)

var (
	ErrNotFloat = errors.New("arguments must be rational numbers")
	ErrNotInt   = errors.New("arguments must allow conversion into an integer (float/decimal numbers are truncated i.e. 3.9 => 3)")
)

var TplFuncs = template.FuncMap{
	"add":   add,
	"sub":   sub,
	"mul":   mul,
	"div":   div,
	"inc":   inc,
	"dec":   dec,
	"float": convFloat,
	"int":   convInt,
}

func convInt(num any) (int, error) {
	i, ok := num.(int)
	if ok {
		return i, nil
	}

	f, ok := num.(float64)
	if ok {
		return int(f), nil
	}

	return -1, ErrNotInt
}

func convFloat(num any) (float64, error) {
	f, ok := num.(float64)
	if ok {
		return f, nil
	}

	i, ok := num.(int)
	if ok {
		return float64(i), nil
	}

	return -1, ErrNotFloat
}

func add(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result += n
	}
	return result, nil
}

func sub(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result -= n
	}
	return result, nil
}

func mul(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result *= n
	}
	return result, nil
}

func div(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result /= n
	}
	return result, nil
}

func inc(num any) (float64, error) {
	n, err := convFloat(num)
	if err != nil {
		return -1, err
	}

	return n + 1, nil
}

func dec(num any) (float64, error) {
	n, err := convFloat(num)
	if err != nil {
		return -1, err
	}

	return n - 1, nil
}
