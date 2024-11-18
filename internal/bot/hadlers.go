package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/SobolevTim/secret-santa-bot/internal/cache"
	"github.com/SobolevTim/secret-santa-bot/internal/database"
	"github.com/SobolevTim/secret-santa-bot/internal/game"
	"github.com/mymmrac/telego"
)

const (
	participantStart            = "await_user_status"      // Статус ожидания подтвержения участия в игре с введеным токеном
	awaitParticipantName        = "await_user_name"        // Статус ожидания имени участиника
	awaitParticipantPreferences = "await_user_preferences" // Статус ожидания предпочтений от участиника
	organizerGameName           = "await_name_game"        // Статус ожидания названия игры
	organizerGameDescription    = "await_game_description" // Статус ожидания описания игры
)

func (b *Bot) handleMessage(msg *telego.Message, service *database.Service, cache *cache.Cache) {
	userID := msg.From.ID
	state, _, cacheOK := cache.Get(userID)
	if strings.HasPrefix(msg.Text, "/") {
		b.handleCommand(msg, service, cache)
	} else if strings.HasPrefix(msg.Text, "santa-game:") {
		b.handleParticipantGame(msg, service, cache)
	} else if cacheOK {
		switch state {
		case participantStart:
			b.handleParticipantIsStart(msg, cache)
		case awaitParticipantName:
			b.hadleInsertParticipantName(msg, service, cache)
		case awaitParticipantPreferences:
			b.handleEndFillParticipant(msg, service, cache)
		case organizerGameName:
			b.handleOrganizerGameDescription(msg, cache)
		case organizerGameDescription:
			b.handleOrganizerGameCreate(msg, service, cache)
		}
	}
}

func (b *Bot) handleCommand(msg *telego.Message, service *database.Service, cache *cache.Cache) {
	switch msg.Text {
	case "/start":
		b.handleStart(msg)
	case "/help":
		b.handleHelp(msg)
	case "/cancel":
		b.handleCancel(msg, cache)
	case "/mygames":
		b.handleMyGames(msg, service)
	case "/game":
		b.handleGame(msg, cache)
	// добавить /giftee,
	default:
		b.SendMessage(msg.Chat.ID, "Неизвестная команда.")
	}
}

func (b *Bot) handleStart(msg *telego.Message) {
	message := fmt.Sprintf("Привет, %s! Я бот для организации тайного санты.\n\nЕсли у тебя есть код участиника, который тебе прислал организатор - отправь мне код!", msg.From.FirstName)
	b.SendMessage(msg.Chat.ID, message)
	message = "Если ты хочешь запустить свою собственную игру - используй команду /game"
	b.SendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleHelp(msg *telego.Message) {
	message := "Я бесплатный бот, для игры \"Тайный Санта\".\nДля участия в игре - отправь мне код. Код тебе должен сообщить организатор игры.\n\nМои основные команды:\n"
	cmd := "/cancel - отменяет ввод. Используй, если застрял и не понимаешь что хочет бот\n/game - организовать свою собственную игру, в которую можно будет пригласить друзей или коллег.\n/mygames - посмотреть список игр, в которых ты участвуешь."
	b.SendMessage(msg.Chat.ID, message+cmd)
}

func (b *Bot) handleCancel(msg *telego.Message, cache *cache.Cache) {
	message := "Ты использовал команду отмены ввода данных\nВоспользуйся /help если не знаешь что делать дальше."
	b.SendMessage(msg.Chat.ID, message)
	cache.ClearUser(msg.Chat.ID)
}

func (b *Bot) handleMyGames(msg *telego.Message, service *database.Service) {
	games, err := service.GetGamesParticipant(database.Participant{UserID: msg.From.ID})
	if err != nil {
		b.SendMessage(msg.Chat.ID, "Ой! 😱\nЧто-то пошло не так! Попробуйте, пожалуйста, еще раз чуть позже!")
		log.Printf("ERROR: %v", err)
		return
	}
	var message string
	if games == nil {
		message = "Пока ты не участвуешь ни в каких новых играх. Возможно, что в игре уже была запущена жеребьевка. Если ты хочешь узнать кто твой подопечный в запущенной игре - используй команду /giftee"
	} else {
		message = "Вот список игр, в которых ты участвуешь и в которых пока еще не было жеребьевки:"
		for _, game := range games {
			//token, name, description
			text := fmt.Sprintf("\n\ntoken игры: %s\nНазвание игры: %s\nОписание игры: %s", game.Token, game.Name, game.Description)
			message += text
		}
	}
	b.SendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleGame(msg *telego.Message, cache *cache.Cache) {
	message := "Ты создаешь свою игру \"Тайный Санта\".\nДля начала, давай придумаем название (например, \"Офисный Тайный Санта\") или можно указать название отдела/службы, чтобы коллегам было проще ориентироваться\n\nДля отмены используй команду /cancel"
	cache.Set(msg.From.ID, organizerGameName, "")
	b.SendMessage(msg.From.ID, message)
}

func (b *Bot) handleParticipantGame(msg *telego.Message, service *database.Service, cache *cache.Cache) {
	token := msg.Text
	ok, game, err := service.CheckGameCodeParticipant(database.Game{Token: token})
	if err != nil {
		b.SendMessage(msg.Chat.ID, "Ой! 😱\nЧто-то пошло не так! Попробуйте, пожалуйста, еще раз чуть позже!")
		log.Printf("ERROR: %v", err)
		return
	}
	if !ok {
		b.SendMessage(msg.Chat.ID, "Хмм 🤔 Похоже что игры с таким кодом не существует!\nНапиши организатору, который прислал этот код, возможно он пересоздал игру или скопировал код не полностью.")
		return
	}
	message := fmt.Sprintf("Игра найдена, ура!🎉\nНазвание игры %s\nОписание игры: %s", game.Name, game.Description)
	b.SendMessage(msg.From.ID, message)
	message = "Если все верно и ты хочешь принять участие именно в этой игре в Тайного Санту, напиши мне: Начинаем\n\nЕсли ты не хочешь присоединяться к этой игре, напиши мне: Отмена"
	cache.Set(msg.From.ID, participantStart, token)
	b.SendMessage(msg.From.ID, message)
}

func (b *Bot) handleParticipantIsStart(msg *telego.Message, cache *cache.Cache) {
	_, token, _ := cache.Get(msg.From.ID)
	text := strings.ToLower(strings.TrimSpace(msg.Text))
	var message string
	switch text {
	case "начинаем":
		cache.Set(msg.From.ID, awaitParticipantName, token)
		message = "Отлично! Теперь давай запишем твои данные\nНапиши мне свое имя!\nТвое имя увидит Тайный Санта.\nРекомендую использовать свое реальное имя, либо псевдоним, который знают все участиники вашей игры в Тайного Санту, иначе Санте будет очень сложно подобрать тебе индивудальный подарок."
	case "отмена":
		cache.ClearUser(msg.From.ID)
		message = "Очень жаль, что ты не хочешь принять участие в игре в Тайного Санту.\nЕсли передумаешь - отправь мне код участиника заново.\n\nЕсли хочешь создать свою игру - используй команду /game"
	default:
		message = "Не понял твое сообщение! 😢\nЧтобы принять участие в игре в Тайного Санту - напиши: Начинаем\nЕсли не хочешь присоединять к игре - напиши: Отмена"
	}
	b.SendMessage(msg.From.ID, message)
}

func (b *Bot) hadleInsertParticipantName(msg *telego.Message, service *database.Service, cache *cache.Cache) {
	p := database.Participant{
		UserID:    msg.From.ID,
		Username:  msg.From.Username,
		FirstName: msg.From.FirstName,
		LastName:  msg.From.LastName,
		Name:      msg.Text,
	}
	_, token, _ := cache.Get(msg.From.ID)
	err := service.CreateParticipant(p, token)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "Ой! 😱\nЧто-то пошло не так! Попробуйте, пожалуйста, еще раз чуть позже!")
		log.Printf("ERROR: %v", err)
		return
	}
	message := fmt.Sprintf("Отлично, %s - твое имя\nА теперь давай запишем твои пожелания к подарку\n\nИмя и пожелания до начала жеребьёвки еще можно будет поменять. Чуть позже я подскажу как это сделать.", msg.Text)
	cache.Set(msg.From.ID, awaitParticipantPreferences, token)
	b.SendMessage(msg.From.ID, message)

}

func (b *Bot) handleEndFillParticipant(msg *telego.Message, service *database.Service, cache *cache.Cache) {
	_, token, _ := cache.Get(msg.From.ID)
	err := service.UpdateParticipantPrefences(database.Participant{GiftPreferences: msg.Text}, token)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "Ой! 😱\nЧто-то пошло не так! Попробуйте, пожалуйста, еще раз чуть позже!")
		log.Printf("ERROR: %v", err)
		return
	}
	message := fmt.Sprintf("%s - твои пожелания к подарку. Уверен, что Таный Санта их обязательно учтет", msg.Text)
	cache.ClearUser(msg.From.ID)
	b.SendMessage(msg.From.ID, message)
	message = "Теперь остается только ждать, когда организатор запустит игру.\nЕсли хочешь посмотреть список доступных команды - используй /help, там есть информация о том, как изменить имя и пожелания к подарку, а также создать свою собственную игру."
	b.SendMessage(msg.From.ID, message)
}

func (b *Bot) handleOrganizerGameDescription(msg *telego.Message, cache *cache.Cache) {
	cache.Set(msg.From.ID, organizerGameDescription, msg.Text)
	message := fmt.Sprintf("%s - название твоей игры.\nТеперь пришли мне описание для игры. Рекомендую прописать правила для участников, например сумму, на которую нужно купить подарок и время начало жеребьевки\n\nЕсли передумал - используй /cancel", msg.Text)
	b.SendMessage(msg.From.ID, message)
}

func (b *Bot) handleOrganizerGameCreate(msg *telego.Message, service *database.Service, cache *cache.Cache) {
	_, name, _ := cache.Get(msg.From.ID)
	token, err := game.GenerateGameCode(*service)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "Ой! 😱\nЧто-то пошло не так! Попробуйте, пожалуйста, еще раз чуть позже!")
		log.Printf("ERROR: %v", err)
		return
	}
	g := database.Game{
		Token:       token,
		OrganizerID: msg.From.ID,
		Name:        name,
		Description: msg.Text,
	}
	err = service.CreateGame(g)
	if err != nil {
		b.SendMessage(msg.Chat.ID, "Ой! 😱\nЧто-то пошло не так! Попробуйте, пожалуйста, еще раз чуть позже!")
		log.Printf("ERROR: %v", err)
		return
	}
	message := fmt.Sprintf("Ты создал игру с названием: %s\nОписанием: %s\nЧтобы участники смогли присоединиться к твоей игре - отправь им этот код:\n%s", g.Name, g.Description, g.Token)
	cache.ClearUser(msg.From.ID)
	b.SendMessage(msg.From.ID, message)
}
