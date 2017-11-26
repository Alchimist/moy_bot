package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"os"
	"gopkg.in/telegram-bot-api.v4"
	coinApi "github.com/miguelmota/go-coinmarketcap"
	"strings"
	"strconv"
	"github.com/Jeffail/gabs"
)



// не сохраняет конфигурацию в файле, только читает из config_bot.json
// не создаёт под нового пользователя его структуру в файле config_new.json
//
//
//
//
//

type Bot_config struct {
	TelegramBotToken	string 	`json:"TelegramBotToken"`
	ChatID				int64 	`json:"chatID"`
	Timer				int 	`json:"timer"`
}

// тип для извлечения списка коинов по юзерам
type interface_map_type map[string]interface{}

var (
	interface_data 		interface_map_type
	bot_config			Bot_config
	CoinList   			map[string]float64
	botToken   			map[string]interface{}
	chatID     			int64
	telegramBotToken 	string
	configFile 			string
	configFile_new 		string
	configFileBot 		string
	jsonParsed 			*gabs.Container
	jsonString 			[]byte
	timer 				int
	alarm 				string
	HelpMsg    = "Это простой мониторинг для подсчёта баланса криптовалюты. Он мониторит валюту по списку и выводит сумму в рублях и общий баланс\n" +
		"Список доступных комманд:\n" +
		"/coin_list - покажет список валюты в мониторинге и их курс \n" +
		"/coin_add [coin_name] [volume] - добавит коин в список мониторинга\n" +
		"/coin_del [coin_name] - удалит коин из списка мониторинга\n" +
		"/coin_timer [time] - таймер мониторинга\n" +
		"/help - отобразить это сообщение\n" +
		"\n"
)

func LoadConfiguration(file string) Bot_config {
	var config Bot_config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func init() {
	CoinList = make(map[string]float64)
	timer = 1 // время по умолчанию
//	file, _ := os.Open("config_bot.json")
//	decoder := json.NewDecoder(file)
//	configuration := Config_bot{}
//	err := decoder.Decode(&configuration)
//	if err != nil {
//		log.Panic(err)
//	}
//	fmt.Println(configuration.TelegramBotToken)

//	flag.StringVar(&configFileBot, "config_bot", "config_bot.json", "config file bot")
//	flag.StringVar(&configFile, "config", "config.json", "config file")
	flag.StringVar(&configFile, "config_new", "config_new.json", "config file")
//	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
//	flag.Int64Var(&chatID, "chatid", 0, "chatId to send messages")

	flag.Parse()
// парсим файл с валютами и пользователями config_new.json
	jsonString, _ = ioutil.ReadFile(configFile)
	jsonParsed, _ = gabs.ParseJSON([]byte(jsonString))

	// загрузка нового конфига config_new.json
	bot_config = LoadConfiguration("config_bot.json")

	load_list()

	log.Printf("telegramBotToken: %s\n", bot_config.TelegramBotToken)
	log.Printf("ChatID: %d\n", bot_config.ChatID)
//	log.Printf("config: %s\n", bot_config)
//	telegramBotToken = botToken["TelegramBotToken"].(string) // "400069657:AAHldU0VZ7ZSfTSU55jnYtJpVnSdvgAqiyM"//

	telegramBotToken = bot_config.TelegramBotToken

	if telegramBotToken == "" {
		log.Print("TelegramBotToken is required")
		os.Exit(1)
	}

//	chatID = int64(botToken["chatID"].(float64))
	chatID = bot_config.ChatID
//	chatID = -263587509
	if chatID == 0 {
		log.Print("chatID is required")
		os.Exit(1)
	}

	timer = bot_config.Timer
	if timer == 0 {
		log.Print("таймер не установлен, выставляется по умолчанию")
		timer = 1
	}

}

func send_notifications(bot *tgbotapi.BotAPI) {
	for site, status := range CoinList {
		if status != 200 {
			alarm := fmt.Sprintf("CRIT - %s ; status: %.0f", site, status)
			bot.Send(tgbotapi.NewMessage(chatID, alarm))
		}
	}
}

func save_list() {

//	configFile, err := os.fi

// сохраняем конфиг с коинами config_new.json
	/*
	j, err := json.Marshal(interface_data)
	err = ioutil.WriteFile(configFile_new, j, 0644)
	if err != nil {
		panic(err)
	}
	*/

	err1 := ioutil.WriteFile(configFile, []byte(jsonParsed.String()), 0644)
	if err1 != nil {
		log.Println("Config file seve error: %s", err1)
	}
	/*

	data, err := json.Marshal(CoinList)
	if err != nil {
		panic(err)
	}
	*/

	/*
	data, err := json.Marshal(interface_data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		panic(err)
	}
	*/

}

func load_list() {
//	data, err := ioutil.ReadFile(configFile)
//	databot, err1 := ioutil.ReadFile(configFileBot)

///	data_new, err2 := ioutil.ReadFile(configFile_new)
///	if err2 != nil {
///		log.Printf("No such file - starting without config_new")
///		return
///	}
	// чтение из файла списка коинов по юзерам
///	if err := json.Unmarshal(data_new, &interface_data); err != nil {
///		log.Printf("Cant read file - starting without config_new")
///	}

//	fmt.Println(databot)
//	fmt.Printf("тип: %T\n", botToken["TelegramBotToken"])
//	fmt.Printf("тип: %T\n", int64(botToken["chatID"].(float64)))
//	log.Printf(string(data_new))
}

var summ float64

func monitor(bot *tgbotapi.BotAPI) {

	// в вечном цикле обходим список урлов раз в 5 мин и сохраняем статус в глобальный map CoinList
	for {
		// сохраняем текущий статус CoinList в файл configFile
		var summ float64 = 0
		alarm = ""

		// ******************** Реализация со старым однопользовательским конфигом
		/*
		for coin, _ := range CoinList {
			// Get info about coin
			coinInfo, err := coinApi.GetCoinData(coin)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf(" %s: ($%.0f) %5.2f\n", coin, coinInfo.PriceUsd, float64(coinInfo.PriceRub) * float64(CoinList[coin]))
				alarm = fmt.Sprintf(" %s: ($%.0f) %5.2f\n", coin, coinInfo.PriceUsd, coinInfo.PriceRub * CoinList[coin]) + alarm
//				bot.Send(tgbotapi.NewMessage(chatID, alarm))
			}
			summ = (coinInfo.PriceRub * CoinList[coin]) + summ
		}
		fmt.Printf("Total: %5.0f руб.\n", summ)
		alarm = fmt.Sprintf("Total: %5.0f\n", summ) + alarm
		bot.Send(tgbotapi.NewMessage(chatID, alarm))
*/
		// **************************************************************************

		/*
		// **************************************************************************
		var price_coin float64
		for user_name := range interface_data { // перебираем все секции с именами пользователей
			coins := interface_data[user_name].(interface{})
			summ = 0
			m := coins.(map[string]interface{})
			f := m["coins"] // берём из конфига секцию "coins" где хранятится количество валюты текущего пользователя
			for coin_name := range f.(map[string]interface{}) { // перебираем коины пользователя
				// Get info about coin
				coinInfo, err := coinApi.GetCoinData(coin_name)
				if err != nil {
					log.Println(err)
				} else {
					price_coin = f.(map[string]interface{})[coin_name].(float64)
					fmt.Printf(" %s: ($%.0f) %0.2f\n", coin_name, coinInfo.PriceUsd, float64(coinInfo.PriceRub) * price_coin)
					alarm = fmt.Sprintf(" %s: ($%.0f) %0.2f\n", coin_name, coinInfo.PriceUsd, price_coin * float64(coinInfo.PriceRub)) + alarm
					summ = (float64(coinInfo.PriceRub) * price_coin) + summ
				}
			}
			fmt.Printf("%s total: %5.0f руб.\n", user_name, summ)
			alarm = fmt.Sprintf("%s total: %5.0f\n", user_name, summ) + alarm
			bot.Send(tgbotapi.NewMessage(chatID, alarm))
		}


		// шлем нотификации
//		send_notifications(bot)
		fmt.Printf("timer %d\n", time.Duration(timer))
		time.Sleep(time.Minute * time.Duration(timer))
//		time.Sleep(time.Second * 180)
	}
		// **************************************************************************
*/

		save_list()
		value1, _ := jsonParsed.S("coins").ChildrenMap()
		//	fmt.Printf("value: %s\n", value1)
		for user_name := range value1 {
			fmt.Printf("value: %s\n", user_name)
			children1, _ := jsonParsed.S("coins").S(user_name).ChildrenMap()
			summ = 0
			alarm = ""
			for coin_name, coin_volume := range children1 {
				// Get info about coin
				coinInfo, err := coinApi.GetCoinData(coin_name)
				if err != nil {
					log.Println(err)
				} else {
					if coin_volume.Data().(float64) != 0 { // если объём равен нулю, значит валюта было добавлена, а потом удалена. Её не считаем
						fmt.Printf(" %s: ($%.0f) %0.2f\n", coin_name, coinInfo.PriceUsd, float64(coinInfo.PriceRub) * coin_volume.Data().(float64))
						alarm = fmt.Sprintf(" %s: ($%.2f) %0.2f\n", coin_name, coinInfo.PriceUsd, coin_volume.Data().(float64) * float64(coinInfo.PriceRub)) + alarm
						summ = (float64(coinInfo.PriceRub) * coin_volume.Data().(float64)) + summ
					}
				}
//				fmt.Printf("валюта: %s, объём: %f, сумма: %.2f\n", coin_name, coin_volume.Data().(float64), coin_volume.Data().(float64) * coinInfo.PriceRub)
			}
			fmt.Printf("%s total: %5.0f руб.\n", user_name, summ)
			alarm = fmt.Sprintf("%s total: %5.0f\n", user_name, summ) + alarm
			bot.Send(tgbotapi.NewMessage(chatID, alarm))
		}
		time.Sleep(time.Minute * time.Duration(timer))
	}
}


func main() {

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Config file: %s", configFile)
	log.Printf("Config file: %s", configFileBot)
	log.Printf("ChatID: %v", chatID)
	log.Printf("Starting monitoring thread")

	go monitor(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprint("буду мониторить: ", CoinList)))

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		reply := "Не понимаю о чём вы. См. /help"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "coin_list":
			reply = "Ваша валюта:\n"
			coins, _ := jsonParsed.S("coins").S(update.Message.From.UserName).ChildrenMap()
			for coin, volume := range coins {
				if volume.Data().(float64) != 0 {
					fmt.Printf("%s: %f\n", coin, volume.Data().(float64))
					reply = reply + fmt.Sprintf("%s: %f\n", coin, volume.Data().(float64))
				}
			}

		case "coin_add":
			reply = ""
				str := strings.Split(update.Message.CommandArguments(), " ")
//			var str1 string

			// Get info about coin
			coinInfo, err := coinApi.GetCoinData(str[0])
			log.Printf("coinInfo %s\n", coinInfo)
			if err != nil {
				log.Println(err)
				reply = "Нет такой валюты"
			} else {
				if len(str) > 1 {
					reply = str[1]
					fl, _ := strconv.ParseFloat(str[1], 64)
					jsonParsed.Set(fl, "coins", update.Message.From.UserName, str[0])
					fmt.Println(jsonParsed.String())
					log.Printf("объёмом %f, пользователю %s, добавлена валюта %s\n", fl, update.Message.From.UserName, str[0])
					reply = fmt.Sprintf("объёмом %f, пользователю %s, добавлена валюта %s\n", fl, update.Message.From.UserName, str[0])
				} else {
					reply = "мало аргументов"
				}
			}
//			reply = "Site added to monitoring list"
		case "coin_timer":
			str := strings.Split(update.Message.CommandArguments(), " ")
			//timer, _ = strconv.ParseFloat(str[0], 64)
			timer, _ = strconv.Atoi(str[0])
			reply = fmt.Sprintf("timer chenged to %d", timer)

		case "coin_del":
			if jsonParsed.Exists("coins", update.Message.From.UserName, update.Message.CommandArguments()) {
				jsonParsed.Set(0.0, "coins", update.Message.From.UserName, update.Message.CommandArguments())
				reply = "коин "+update.Message.CommandArguments()+" удалён"
			} else {
				reply = "ошибка удаления"
			}

/*
			coins := interface_data[update.Message.From.UserName].(interface{})
			m := coins.(map[string]interface{})
			f := m["coins"] // берём из конфига секцию "coins" где хранятится количество валюты текущего пользователя
			delete(f.(map[string]interface{}), update.Message.CommandArguments())
			reply = "Site deleted from monitoring list"
*/
		case "help":
			reply = HelpMsg
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
