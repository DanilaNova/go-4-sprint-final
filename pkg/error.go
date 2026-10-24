package error

import (
	"errors"
	"fmt"
)

// Тест "неизвестный_тип_тренировки_-_проверка_текста_ошибки" из "spentcalories_test.go" подразумевает конкретный текст ошибки.
// Остальные ошибки выводятся на русском языке для единообразия.
var (
	ErrIncorrectFormat = errors.New("неверный формат строки")
	ErrUnknownActivity = errors.New("неизвестный тип тренировки")

	ErrValueIsZeroOrLess    = errors.New("значение меньше или равно нулю")
	ErrStepsIsZeroOrLess    = fmt.Errorf("%w: шаги", ErrValueIsZeroOrLess)
	ErrWeightIsZeroOrLess   = fmt.Errorf("%w: вес", ErrValueIsZeroOrLess)
	ErrHeightIsZeroOrLess   = fmt.Errorf("%w: высота", ErrValueIsZeroOrLess)
	ErrDurationIsZeroOrLess = fmt.Errorf("%w: продолжительность", ErrValueIsZeroOrLess)
)
