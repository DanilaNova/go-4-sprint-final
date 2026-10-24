package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	. "github.com/DanilaNova/go-4-sprint-final/pkg"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining достаёт из строки информацию о тренировке.
//
// Строка передаётся в формате "123,Тип,3h4m".
// Типом активности может быть "Ходьба" или "Бег"
//
// # Типы ошибок:
//
// [ErrStepsIsZeroOrLess] |
// [ErrDurationIsZeroOrLess] |
// ошибки [strconv.Atoi] |
// ошибки [time.ParseDuration]
func parseTraining(data string) (int, string, time.Duration, error) {
	split := strings.Split(data, ",")
	if len(split) != 3 {
		return 0, "", 0, fmt.Errorf("%w, ожидалось: 123,Тип,3h4m", ErrIncorrectFormat)
	}

	steps, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, ErrStepsIsZeroOrLess
	}

	activity := split[1]

	duration, err := time.ParseDuration(split[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, ErrDurationIsZeroOrLess
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

// meanSpeed расчитывает среднюю арифметическую скорость в км/ч
//
// # Высота передаётся в метрах
//
// Если продолжительность равна нулю, возвращается 0
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo выводит информацию о тренировке
//
// Вес передаётся в килограммах
// Высота передаётся в метрах
//
// # Типы ошибок:
//
// [ErrUnknownActivity] |
// ошибки [parseTraining] |
// ошибки [WalkingSpentCalories] |
// ошибки [RunningSpentCalories]
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch strings.ToLower(activity) {
	case "ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", ErrUnknownActivity
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	activity_distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`, activity, duration.Hours(), activity_distance, speed, calories), nil
}

// RunningSpentCalories расчитывает потраченные калории при беге
//
// Вес передаётся в килограммах
// Высота передаётся в метрах
//
// # Типы ошибок
//
// [ErrStepsIsZeroOrLess] |
// [ErrWeightIsZeroOrLess] |
// [ErrHeightIsZeroOrLess] |
// [ErrDurationIsZeroOrLess]
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrStepsIsZeroOrLess
	}
	if weight <= 0 {
		return 0, ErrWeightIsZeroOrLess
	}
	if height <= 0 {
		return 0, ErrHeightIsZeroOrLess
	}
	if duration <= 0 {
		return 0, ErrDurationIsZeroOrLess
	}

	speed := meanSpeed(steps, height, duration)

	return weight * speed * duration.Hours(), nil
}

// WalkingSpentCalories расчитывает потраченные калории при ходьбе
//
// Вес передаётся в килограммах
// Высота передаётся в метрах
//
// # Типы ошибок
//
// ошибки [RunningSpentCalories]
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	return calories * walkingCaloriesCoefficient, err
}
