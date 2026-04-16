package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// validateInputs общая функция для проверки введённых данных
func validateInputs(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return errors.New("время должно быть больше 0")
	}
	return nil
}

// parseTraining парсит строку, переводит данные
// из строки в соответствующие типы и возвращает эти значения.
func parseTraining(data string) (int, string, time.Duration, error) {

	// парсим строку
	parseData := strings.Split(data, ",")

	// проверяем кол-во данных в полученном слайсе
	if len(parseData) != 3 {
		return 0, "", 0, fmt.Errorf("недостаточно данных: ожидается 3 значения, получено %d", len(parseData))
	}

	// преобразуем строку с количеством шагов в int,
	// затем проверяем их количество
	steps, err := strconv.Atoi(strings.TrimSpace(parseData[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов: %w", err)
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("Количество шагов должно быть больше 0")
	}
	// проверяем тип тренировки
	trainType := strings.TrimSpace(parseData[1])
	if trainType == "" {
		return 0, "", 0, errors.New("тип тренировки не указан")
	}

	// преобразуем строку со временем в time.Duration
	duration, err := time.ParseDuration(strings.TrimSpace(parseData[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать время: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("время должно быть больше 0")
	}

	return steps, trainType, duration, nil
}

// getStepLenght вычисляет длину шага
func getStepLength(height float64) float64 {
	if height > 0 {
		return height * stepLengthCoefficient
	}
	return lenStep // средняя длина шага, если рост не указан
}

// distance расчитывает дистанцию в километрах
func distance(steps int, height float64) float64 {

	stepLength := getStepLength(height)
	return (float64(steps) * stepLength) / float64(mInKm)
}

// meanSpeed вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}
	distKm := distance(steps, height)

	hours := duration.Hours()

	return distKm / hours
}

// RunningSpentCalories принимает данные пользователя и
// вычисляет кол-во потраченных каллорий во время бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// производим проверки с помощью validateInputs()
	if err := validateInputs(steps, weight, height, duration); err != nil {
		return 0, err
	}

	// вычисляем среднюю скорость с помощью функции meanSpeed()
	mSpeed := meanSpeed(steps, height, duration)

	// переводим время в минуты для расчета каллорий
	durationInMinutes := duration.Minutes()

	// расчитываем каллории
	calories := (weight * mSpeed * durationInMinutes) / float64(minInH)

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// проверки с помощью validateInputs()
	if err := validateInputs(steps, weight, height, duration); err != nil {
		return 0, err
	}

	// вычисляем среднюю скорость с помощью функции meanSpeed()
	mSpeed := meanSpeed(steps, height, duration)

	// переводим время в минуты для расчета каллорий
	durationInMinutes := duration.Minutes()

	// расчитываем каллории
	calories := ((weight * mSpeed * durationInMinutes) / float64(minInH)) * walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo принимает данные о тренировке и пользователе
// и выводит подробную информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {

	// получаем значения из строки
	steps, train, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	// проводим остальные проверки полученных данных
	if weight <= 0 {
		return "", errors.New("Вес должен быть больше 0")
	}

	if height <= 0 {
		return "", errors.New("Рост должен быть больше 0")
	}
	// общие расчеты
	dist := distance(steps, height)
	mSpeed := meanSpeed(steps, height, duration)

	switch train {
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			train, duration.Hours(), dist, mSpeed, calories), nil
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			train, duration.Hours(), dist, mSpeed, calories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}
