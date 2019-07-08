package main

import (
	//"encoding/json"
	"fmt"
	"os"
	"strconv"

	//"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/ini.v1"

	// MySQL (mail.ru)
	_ "github.com/go-sql-driver/mysql"
	"github.com/mailru/dbr"

	//m "github.com/ValidatorCenter/minter-go-sdk"

	api "github.com/MinterTeam/minter-node-go-api"
)

// пока данные будем хранить в памяти
var (
	dbMySQL *dbr.Connection
	mnt     *api.MinterNodeApi

	TgTokenAPI   string // Токен к API телеграма
	TgTimeUpdate int64  // Время в сек. обновления статуса
	DBAddress    string // DB_URI

	// TODO: мультиязычность
	HelpMsg = "Это простой мониторинг доступности мастерноды валидатора и краткая информация о ней.\n" +
		"Список доступных комманд:"
)

////////////////////////////////////////////////
//SQL::

// Очистка и создание таблицы для ...
func creatTabl_Acc_MySQl(db *dbr.Connection) {
	sess := db.NewSession(nil)
	////////////////////////////////////////////////////////////////////////////
	delMy_node_acc := `DROP TABLE IF EXISTS node_acc`
	sess.Exec(delMy_node_acc)
	schemaMy_node_acc := `
		   		CREATE TABLE node_acc (
					id INT NOT NULL AUTO_INCREMENT,
		   			address VARCHAR(128),
					comment VARCHAR(256),
		   			priority  INT UNSIGNED,
					coin VARCHAR(10),
		   			prc  INT UNSIGNED,
					PRIMARY KEY (id)
		   		)
		   		`
	sess.Exec(schemaMy_node_acc)
	fmt.Println("OK", "...очищена - node_acc")
	////////////////////////////////////////////////////////////////////////////
}

////////////////////////////////////////////////

// Сам мониторинг! как горутина!
func monitor(bot *tgbotapi.BotAPI) {
	// бесконечный цикл
	for {
		/*ReturnValid()

		for _, oneUser := range allUser {
			if !getStatusValid(oneUser.PubKey) && oneUser.Notification == true {
				//Алам!
				fmt.Println("NOOOOO! ", oneUser.UserName)
				// отправляем пользователю сообщение
				msg := tgbotapi.NewMessage(oneUser.ChatID, "Нода не в валидаторах!")
				bot.Send(msg)
			}
		}*/

		fmt.Printf("Пауза %dсек.... в этот момент лучше прерывать\n", TgTimeUpdate)
		time.Sleep(time.Second * time.Duration(TgTimeUpdate)) // пауза
	}
}

func main() {
	ConfFileName := "cmc0.ini"
	cmdClearDB := false

	// проверяем есть ли входной параметр/аргумент
	if len(os.Args) == 2 {
		if os.Args[1] == "new" {
			cmdClearDB = true
		} else {
			ConfFileName = os.Args[1]
		}
	}
	fmt.Printf("INI=%s\n", ConfFileName)

	// INI
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, ConfFileName)
	if err != nil {
		fmt.Println("Ошибка загрузки INI файла:", err.Error())
		return
	} else {
		fmt.Println("...данные с INI файла = загружены!")
	}
	//secMN := cfg.Section("masternode")
	//MnAddress = secMN.Key("ADDRESS").String()
	secDB := cfg.Section("database")
	DBAddress = secDB.Key("ADDRESS").String()
	//netMN := cfg.Section("network")
	//CoinMinter = netMN.Key("COINNET").String()
	secTG := cfg.Section("telegram")
	TgTokenAPI = secTG.Key("TOKEN").String()
	_TgTimeUpdate, err := strconv.Atoi(secTG.Key("TIMEUPDATE").String())
	if err != nil {
		fmt.Println(err)
		TgTimeUpdate = 60
	}
	TgTimeUpdate = int64(_TgTimeUpdate)

	//////////////////////////////////////////////////////
	// DB:: MySQL
	dbMySQL, err = dbr.Open("mysql", DBAddress, nil)
	if err != nil {
		fmt.Println("Ошибка соединения с БД:", err)
		os.Exit(1)
	}
	defer dbMySQL.Close()

	fmt.Println(time.Now())

	if cmdClearDB == true {
		// очистка и создание таблиц в базе MySQL
		creatTabl_Acc_MySQl(dbMySQL)
		return
	}

	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(TgTokenAPI)
	if err != nil {
		fmt.Println("Ошибка соединения с Telegram:", err.Error())
		return
	}

	bot.Debug = true
	fmt.Printf("Авторизован: %s\n", bot.Self.UserName)

	// Некие еще предварительные действия перед запуском:
	// TODO: ....

	// БОТ:
	// в отдельном потоке запускаем функцию мониторинга
	go monitor(bot)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := ""
		if update.Message == nil {
			continue
		}
		/*if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}*/

		// логируем от кого какое сообщение пришло
		fmt.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {

		// выводим информацию о боте
		case "start":
			reply = HelpMsg
		case "help":
			reply = HelpMsg

		// задать язык пользователя
		case "i18n":
			// TODO: реализация

		// выводим информацию о мастерноде(валидаторе!)
		case "node_info":
			// TODO: реализация

		// добавить мастерноду в список мониторинга UPD: она же и ->изменить pubkey у мастерноды
		case "node_add":
			// TODO: реализация

		// удаление мастерноды
		case "node_del":
			// TODO: реализация

		// изменить статус уведомления да/нет
		case "notification":
			// TODO: реализация

		// обработка пропуска блоков да/нет, при ДА будет отслеживание и остановка ноды
		case "node_stoping":
			// TODO: реализация

		// задать параметр - через сколько подряд блоков будет отправлена команда, по умолчанию: 3
		case "amnt_blocks":
			// TODO: реализация

		// задать подписанную команду остановки ноды!
		case "rlp_stoping":
			// TODO: реализация

			// Прочие команды для бота:
			// TODO: case "....":

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err = bot.Send(msg)
		if err != nil {
			fmt.Println("Ошибка отправки сообщения:", err)
		}
	}
}
