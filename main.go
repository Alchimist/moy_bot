package main

import (
	"encoding/json"
	"flag"
	"fmt"
<<<<<<< HEAD
//	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"crypto/tls"
	"gopkg.in/telegram-bot-api.v4"
//	"strconv"
//	"strings"
//	"encoding/xml"
=======
	"io/ioutil"
	"log"
	"time"
	"os"
	"gopkg.in/telegram-bot-api.v4"
	coinApi "github.com/miguelmota/go-coinmarketcap"
	"strings"
	"strconv"
>>>>>>> new_features
)

//
//Site list map  "URL" - status
//Site status
//0 - never checked
//1 - timeout
//200 - ok
//other statuses - crit

var (
<<<<<<< HEAD
	SiteList   map[string]int
=======
	SiteList   map[string]float64
>>>>>>> new_features
	botToken   map[string]interface{}
	chatID     int64
	telegramBotToken string
	configFile string
	configFileBot string
<<<<<<< HEAD
	HelpMsg    = "Это простой мониторинг доступности сайтов. Он обходит сайты в списке и ждет что он ответит 200, если возвращается не 200 или ошибки подключения, то бот пришлет уведомления в групповой чат\n" +
		"Список доступных комманд:\n" +
		"/site_list - покажет список сайтов в мониторинге и их статусы (про статусы ниже)\n" +
		"/site_add [url] - добавит url в список мониторинга\n" +
		"/site_del [url] - удалит url из списка мониторинга\n" +
		"/help - отобразить это сообщение\n" +
		"\n" +
		"У сайтов может быть несколько статусов:\n" +
		"0 - никогда не проверялся (ждем проверки)\n" +
		"1 - ошибка подключения \n" +
		"200 - ОК-статус" +
		"все остальные http-коды считаются некорректными"
)

func init() {
	SiteList = make(map[string]int)

=======
	timer int
	alarm string
	HelpMsg    = "Это простой мониторинг для подсчёта баланса криптовалюты. Он мониторит валюту по списку и выводит сумму в рублях и общий баланс\n" +
		"Список доступных комманд:\n" +
		"/coin_list - покажет список валюты в мониторинге и их курс \n" +
		"/coin_add [coin_name] [volume] - добавит коин в список мониторинга\n" +
		"/coin_del [coin_name] - удалит коин из списка мониторинга\n" +
		"/coin_timer [time] - таймер мониторинга\n" +
		"/help - отобразить это сообщение\n" +
		"\n"
)

func init() {
	SiteList = make(map[string]float64)
	timer = 60 // время по умолчанию
>>>>>>> new_features
//	file, _ := os.Open("config_bot.json")
//	decoder := json.NewDecoder(file)
//	configuration := Config_bot{}
//	err := decoder.Decode(&configuration)
//	if err != nil {
//		log.Panic(err)
//	}
//	fmt.Println(configuration.TelegramBotToken)

	flag.StringVar(&configFileBot, "config_bot", "config_bot.json", "config file bot")
	flag.StringVar(&configFile, "config", "config.json", "config file")
//	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
//	flag.Int64Var(&chatID, "chatid", 0, "chatId to send messages")

	flag.Parse()

	load_list()

	telegramBotToken = botToken["TelegramBotToken"].(string) // "400069657:AAHldU0VZ7ZSfTSU55jnYtJpVnSdvgAqiyM"//
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}

//	chatID = -263587509
	chatID = int64(botToken["chatID"].(float64))
	if chatID == 0 {
		log.Print("-chatid is required")
		os.Exit(1)
	}

}

func send_notifications(bot *tgbotapi.BotAPI) {
	for site, status := range SiteList {
		if status != 200 {
<<<<<<< HEAD
			alarm := fmt.Sprintf("CRIT - %s ; status: %v", site, status)
=======
			alarm := fmt.Sprintf("CRIT - %s ; status: %.0f", site, status)
>>>>>>> new_features
			bot.Send(tgbotapi.NewMessage(chatID, alarm))
		}
	}
}

func save_list() {
	data, err := json.Marshal(SiteList)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		panic(err)
	}
}

func load_list() {
	data, err := ioutil.ReadFile(configFile)
	databot, err1 := ioutil.ReadFile(configFileBot)

	if err != nil {
		log.Printf("No such file - starting without config")
		return
	}

	if err1 != nil {
		log.Printf("No such file - starting without config bot")
		return
	}

	if err = json.Unmarshal(data, &SiteList); err != nil {
		log.Printf("Cant read file - starting without config")
		return
	}

	if err = json.Unmarshal(databot, &botToken); err != nil {
		log.Printf("Cant read file - starting without configbot")
		return
	}

//	fmt.Println(databot)

<<<<<<< HEAD
	fmt.Printf("тип: %T\n", botToken["TelegramBotToken"])
	fmt.Printf("тип: %T\n", int64(botToken["chatID"].(float64)))
=======
//	fmt.Printf("тип: %T\n", botToken["TelegramBotToken"])
//	fmt.Printf("тип: %T\n", int64(botToken["chatID"].(float64)))
>>>>>>> new_features

	log.Printf(string(data))
}

<<<<<<< HEAD
func monitor(bot *tgbotapi.BotAPI) {

=======
var summ float64

func monitor(bot *tgbotapi.BotAPI) {

//	var summ float64
	/*
	// параметр для http клиента чтобы принимал все ssl сертификаты
>>>>>>> new_features
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

<<<<<<< HEAD
=======
	// важно указать таймаут http соединения, иначе он вечно будет висеть и мы не увидим
	// проблемы если сервер просто упал

>>>>>>> new_features
	var httpclient = &http.Client{
		Timeout: time.Second * 10,
		Transport: tr,
	}
<<<<<<< HEAD

	for {
		save_list()
		for site, _ := range SiteList {
			response, err := httpclient.Get(site)
=======
*/
	// в вечном цикле обходим список урлов раз в 5 мин и сохраняем статус в глобальный map SiteList
	for {
		// сохраняем текущий статус SiteList в файл configFile
		save_list()
		var summ float64
		alarm = ""
		for site, _ := range SiteList {
/*			response, err := httpclient.Get(site)
>>>>>>> new_features
			if err != nil {
				SiteList[site] = 1
				log.Printf("Status of %s: %s", site, "1 - Connection refused")
			} else {
				log.Printf("Status of %s: %s", site, response.Status)
				SiteList[site] = response.StatusCode
			}
<<<<<<< HEAD
		}
		send_notifications(bot)
		time.Sleep(time.Minute * 5)
	}
}

func main() {
=======
*/			// Get info about coin
			coinInfo, err := coinApi.GetCoinData(site)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf(" %s: ($%.0f) %5.0f\n", site, coinInfo.PriceUsd, coinInfo.PriceRub * SiteList[site])
				alarm = fmt.Sprintf(" %s: ($%.0f) %5.0f\n", site, coinInfo.PriceUsd, coinInfo.PriceRub * SiteList[site]) + alarm
//				bot.Send(tgbotapi.NewMessage(chatID, alarm))
			}
			summ = (coinInfo.PriceRub * SiteList[site]) + summ
		}
		fmt.Printf("Total: %5.0f руб.\n", summ)
		alarm = fmt.Sprintf("Total: %5.0f\n", summ) + alarm
		bot.Send(tgbotapi.NewMessage(chatID, alarm))

		/*
		// Get info about coin
		coinInfo, err := coinApi.GetCoinData("bitcoin")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf(" BTC: %5.0f\n", coinInfo.PriceRub * 0.05086186)
		}

		var summ = coinInfo.PriceRub * 0.05086186

		// Get info about coin
		coinInfo, err = coinApi.GetCoinData("monero")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf(" XMR: %5.0f\n", coinInfo.PriceRub * 0.243)
		}

		summ = (coinInfo.PriceRub * 0.243) + summ
		// Get info about coin
		coinInfo, err = coinApi.GetCoinData("ethereum")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf(" ETH: %5.0f\n", coinInfo.PriceRub * 0.1)
		}
		summ = (coinInfo.PriceRub * 0.1) + summ

		// Get info about coin
		coinInfo, err = coinApi.GetCoinData("dash")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("DASH: %5.0f\n", coinInfo.PriceRub * 0.12556303)
		}
		summ = (coinInfo.PriceRub * 0.12556303) + summ
		fmt.Printf("Total: %5.0f руб.\n", summ)
*/
		// шлем нотификации
//		send_notifications(bot)
		time.Sleep(time.Minute * time.Duration(timer))
//		time.Sleep(time.Second * 180)
	}
}


func main() {
/*
	// Get global market data
	marketInfo, err := coinApi.GetMarketData()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(marketInfo)
	}
*/
/*
	// Get info about coin
	coinInfo, err := coinApi.GetCoinData("monero")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(coinInfo.PriceRub)
	}

	// Get info about coin
	coinInfo, err = coinApi.GetCoinData("ethereum")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(coinInfo.PriceRub)
	}

	// Get info about coin
	coinInfo, err = coinApi.GetCoinData("dash")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(coinInfo.PriceRub)
	}
*/

>>>>>>> new_features
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}
<<<<<<< HEAD

=======
/*
>>>>>>> new_features
	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Config file: %s", configFile)
	log.Printf("Config file: %s", configFileBot)
	log.Printf("ChatID: %v", chatID)
	log.Printf("Starting monitoring thread")
<<<<<<< HEAD
	go monitor(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprint("Я живой; вот сайты которые буду мониторить: ", SiteList)))
=======
*/
	go monitor(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprint("буду мониторить: ", SiteList)))
>>>>>>> new_features

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
<<<<<<< HEAD
		case "site_list":
			sl, _ := json.Marshal(SiteList)
			reply = string(sl)

		case "site_add":
			SiteList[update.Message.CommandArguments()] = 0
			reply = "Site added to monitoring list"

		case "site_del":
=======
		case "coin_list":
			sl, _ := json.Marshal(SiteList)
			reply = string(sl)

		case "coin_add":
			str := strings.Split(update.Message.CommandArguments(), " ")
//			var str1 string
			if str[1] != "" {
				reply = str[1]
				f, _ := strconv.ParseFloat(str[1], 64)
				SiteList[str[0]] = f
//				str[0] = f
				log.Printf("%f\n", f)
			} else {
				reply = "фиг вам"
			}
//			reply = "Site added to monitoring list"
		case "coin_timer":
			str := strings.Split(update.Message.CommandArguments(), " ")
			//timer, _ = strconv.ParseFloat(str[0], 64)
			timer, _ = strconv.Atoi(str[0])
			reply = "timer chenged"

		case "coin_del":
>>>>>>> new_features
			delete(SiteList, update.Message.CommandArguments())
			reply = "Site deleted from monitoring list"
		case "help":
			reply = HelpMsg
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
