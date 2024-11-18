package game

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/SobolevTim/secret-santa-bot/internal/database"
)

const (
	gameCodePrefix    = "santa-game:"
	randomPartLength  = 17
	letters           = "abcdefghijklmnopqrstuvwxyz"
	digits            = "0123456789"
	specialCharacters = "!@#$%^&*"
	allCharacters     = letters + digits + specialCharacters
	requiredDigits    = 3
	requiredSpecials  = 2
)

func GenerateGameCode(service database.Service) (string, error) {
	var code string
	for {
		// Создаем слайс для случайных символов
		randomPart := make([]rune, randomPartLength)

		// Добавляем минимум 3 цифры
		for i := 0; i < requiredDigits; i++ {
			randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
			randomPart[i] = rune(digits[randIndex.Int64()])
		}

		// Добавляем минимум 2 специальных символа
		for i := requiredDigits; i < requiredDigits+requiredSpecials; i++ {
			randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(specialCharacters))))
			randomPart[i] = rune(specialCharacters[randIndex.Int64()])
		}

		// Заполняем оставшиеся символы
		for i := requiredDigits + requiredSpecials; i < randomPartLength; i++ {
			randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allCharacters))))
			randomPart[i] = rune(allCharacters[randIndex.Int64()])
		}

		// Перемешиваем символы вручную
		for i := range randomPart {
			j, _ := rand.Int(rand.Reader, big.NewInt(int64(randomPartLength)))
			randomPart[i], randomPart[j.Int64()] = randomPart[j.Int64()], randomPart[i]
		}

		// Формируем полный код
		code = fmt.Sprintf("%s%s", gameCodePrefix, string(randomPart))

		exists, err := service.CheakToken(code)

		if err != nil {
			return "", fmt.Errorf("ошибка при проверки токена на уникальность в GenerateGameCode: %w", err)
		}

		// Если код уникален, выходим из цикла
		if !exists {
			break
		}
	}
	return code, nil
}
