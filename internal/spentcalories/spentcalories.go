package spentcalories

import (
	"errors"
	"fmt"
	"log"
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

// parseTraining парсит строку, переводит данные
// из строки в соответствующие типы и возвращает эти значения.
func parseTraining(data string) (int, string, time.Duration, error) {

	// парсим строку
	parseData := strings.Split(data, ",")

	// проверяем кол-во данных в полученном слайсе
	if len(parseData) != 3 {
		return 0, "", 0, errors.New("Недостаточно данных в строке")
	}

	// преобразуем строку с количеством шагов в int,
	// затем проверяем их количество
	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, "", 0, errors.New("Не удалось преобразовать кол-во шагов")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("Количество шагов должно быть больше 0")
	}

	// преобразуем строку со временем в time.Duration
	duration, err := time.ParseDuration(parseData[2])
	if err != nil {
		return 0, "", 0, errors.New("Не удалось получить время")
	}

	return steps, parseData[1], duration, nil
}

// distance принимает количество шагов и рост
// пользователя в метрах, а возвращает
// дистанцию в километрах
func distance(steps int, height float64) float64 {

	// вычисляем длину шага на основе роста
	stepLength := height * stepLengthCoefficient

	// вычисляем дистанцию
	distance := (float64(steps) * stepLength) / float64(mInKm)

	return distance
}

// meanSpeed вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	// проверяем кол-во времени
	if duration <= 0 {
		return 0
	}
	// вычисляем дистанцию с помощью distance()
	distKm := distance(steps, height)

	// переводим продолжительность в часы функцией из пакета "time"
	hours := duration.Hours()

	// возвращаем среднюю скорость
	return distKm / hours
}

// RunningSpentCalories принимает данные пользователя и
// вычисляет кол-во потраченных каллорий во время бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// производим проверки
	if steps <= 0 {
		return 0, errors.New("Количество шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0, errors.New("Вес должен быть больше 0")
	}

	if height <= 0 {
		return 0, errors.New("Рост должен быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("Время должно быть больше 0")
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
	if steps <= 0 {
		return 0, errors.New("Количество шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0, errors.New("Вес должен быть больше 0")
	}

	if height <= 0 {
		return 0, errors.New("Рост должен быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("Время должно быть больше 0")
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
		log.Println(err)
		return "", err
	}
	// проводим остальные проверки полученных данных
	if weight <= 0 {
		return "", errors.New("Вес должен быть больше 0")
	}

	if height <= 0 {
		return "", errors.New("Рост должен быть больше 0")
	}

	switch train {
	case "Бег":
		dist := distance(steps, height)

		mSpeed := meanSpeed(steps, height, duration)

		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			train, duration, dist, mSpeed, calories), nil
	case "Ходьба":
		dist := distance(steps, height)

		mSpeed := meanSpeed(steps, height, duration)

		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			train, duration, dist, mSpeed, calories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}
