package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/DanilaNova/go-4-sprint-final/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

var (
	ErrIncorrectFormat = errors.New("неверный формат строки")

	ErrValueIsZeroOrLess    = errors.New("значение меньше или равно нулю")
	ErrStepsIsZeroOrLess    = fmt.Errorf("%w: шаги", ErrValueIsZeroOrLess)
	ErrDurationIsZeroOrLess = fmt.Errorf("%w: продолжительность", ErrValueIsZeroOrLess)
)

// parsePackage достаёт из строки количество шагов и время прогулки.
//
// Строка должна передаваться в формате "123,3h4m".
//
// # Типы ошибок:
//
// [ErrIncorrectFormat] |
// [ErrStepsIsZeroOrLess] |
// [ErrDurationIsZeroOrLess] |
// [strconv.NumError]
func parsePackage(data string) (int, time.Duration, error) {
	split := strings.Split(data, ",")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("%w, ожидалось: 123,3h4m", ErrIncorrectFormat)
	}

	steps, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, err
	} else if steps <= 0 {
		return 0, 0, ErrStepsIsZeroOrLess
	}

	walkTime, err := time.ParseDuration(split[1])
	if err != nil {
		return 0, 0, err
	} else if walkTime <= 0 {
		return 0, 0, ErrDurationIsZeroOrLess
	}

	return steps, walkTime, nil
}

// DayActionInfo выводит информацию о прогулке.
//
// Строка data должна иметь формат, описанный в [parsePackage].
// Вес передаётся в килограммах.
// Высота передаётся в метрах.
//
// В случае ошибки возращается пустая строка и ошибка отправляется в [log]
func DayActionInfo(data string, weight, height float64) string {
	steps, walkTime, err := parsePackage(data)
	if err != nil {
		log.Println(fmt.Errorf("ошибка в обработке информации о прогулке: %w", err))
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	spent, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkTime)
	if err != nil {
		log.Println(fmt.Errorf("ошибка в расчёте потраченных каллорий: %w", err))
		return ""
	}

	return fmt.Sprintf(`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`, steps, distance, spent)
}
