package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат тренировки")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, err
	}

	return steps, strings.TrimSpace(parts[1]), duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / mInKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s.\nКоличество шагов: %d.\nДистанция: %.2f км.\nСкорость: %.2f км/ч.\nПотрачено калорий: %.2f.\n",
		activity,
		steps,
		distance(steps, height)/mInKm,
		meanSpeed(steps, height, duration),
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}

	return weight * meanSpeed(steps, height, duration) * duration.Hours(), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}

	return weight * meanSpeed(steps, height, duration) * duration.Hours() * walkingCaloriesCoefficient, nil
}
