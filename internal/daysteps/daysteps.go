package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage принимает строку с данными
// разделяет их в слайс и переводит в нужные значения
// затем возвращет эти значения
func parsePackage(data string) (int, time.Duration, error) {
	// преобразуем полученные данные в слайс
	parseData := strings.Split(data, ",")
	
	// проверяем длину слайса
	if len(parseData) != 2 {
		return 0, 0, errors.New("Недостаточно данных для преобразования")
	}
	// преобразуем шаги в int
	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, 0, errors.New("Ошибка преобразования шагов")
	}
	// проверяем количество шагов
	if steps <= 0 {
		return 0, 0, errors.New("Количество шагов должно быть больше 0")
	}
	// преобразуем шаги
	duration, err := time.ParseDuration(parseData[1])
	if err != nil {
		return 0, 0, errors.New("Не удалось преобразовать время")
	}

	return steps, duration, nil
}

// DayActionInfo парсит строку с данными с помощью
// функции parsePackage, вычисляет дистанцию в
// километрах и кол-во потраченных каллорий
// Возвращает строку с полученной информацией
func DayActionInfo(data string, weight, height float64) string {

	// парсим строку на кол-во шагов и время прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// проверяем кол-во шагов
	if steps <= 0 {
		return ""
	}
	// вычисляем дистанцию и калории
	distInKm := (float64(steps) * stepLength) / float64(mInKm)  

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\n
		Вы сожгли %.2f ккал.", steps, distInKm, calories)
}
