package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type student struct {
	Trash        int
	Shopping     int
	BigShopping  int
	MoneySpent   int
	FlatCleaning int
	RoomCleaning int
}

var chatId := 0 //chatId to add bot (int)
var adminChatId := 0 //admins chatId to manage group (int)
const botToken := "" //bot token here? you can take it from @BotFather 

var studMap = make(map[string]*student)

var shopButton = tgbotapi.KeyboardButton{
	Text: "/shopping",
}
var trashButton = tgbotapi.KeyboardButton{
	Text: "/throw_trash",
}
var bigShopButton = tgbotapi.KeyboardButton{
	Text: "/big_shopping",
}
var roomButton = tgbotapi.KeyboardButton{
	Text: "/room_cleaning",
}
var flatButton = tgbotapi.KeyboardButton{
	Text: "/flat_cleaning",
}

var listButton = tgbotapi.KeyboardButton{
	Text: "/list",
}

var wg sync.WaitGroup

func recurTimer() {
	save()
	time.AfterFunc(duration(), recurTimer)
}
func save() string {
	file, err := os.Create("duty_struct.json")
	if err != nil {
		return "Didnt save"
	}
	defer file.Close()
	buff, _ := json.Marshal(studMap)
	file.Write(buff)
	return "Saved succesfully"
}

func recov() string {
	file, err := os.Open("duty_struct.json")
	if err != nil {
		return "File doesnt exist"
	}
	defer file.Close()

	stat, _ := file.Stat()
	buff := make([]byte, int(stat.Size()))
	file.Read(buff)
	json.Unmarshal(buff, &studMap)
	return "recovered succesfully"
}

func duration() time.Duration {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 8, 0, 0, 0, t.Location())
	//n := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	if t.After(n) {
		n = n.Add(24 * time.Hour)
	}
	d := n.Sub(t)
	return d
}

var replyButtons = tgbotapi.ReplyKeyboardMarkup{Keyboard: [][]tgbotapi.KeyboardButton{{listButton, trashButton}, {shopButton, bigShopButton}, {roomButton, flatButton}}}

func respondTemplate(msg *tgbotapi.Message, cmd string) string {
	result := cmd + " "
	user := msg.From
	result += "was added to " + user.UserName
	return result
}

func addUser(studMap map[string]student, userName string) {
	studMap[userName] = *new(student)
}

func centrWord(word string, lineLen int) string {
	spaces := (lineLen - len(word)) / 2
	result := ""
	lastSpace := ""
	for i := 0; i < spaces; i++ {
		result += " "
	}
	if spaces*2 != lineLen-len(word) {
		lastSpace = " "
	}
	return lastSpace + result + word + result
}

func drawTable(respond *tgbotapi.MessageConfig) {
	for key, val := range studMap {
		respond.Text += "```\n===========================```\n"
		respond.Text += "  *\\" + key + "*" + "\n"
		respond.Text += "```\n"
		respond.Text += "---------------------------" + "\n"
		respond.Text += "shop|trsh|mt|rm|fl|  money  " + "\n"
		respond.Text += "---------------------------" + "\n"
		respond.Text += centrWord(strconv.Itoa(val.Shopping), 4) + "|"
		respond.Text += centrWord(strconv.Itoa(val.Trash), 4) + "|"
		respond.Text += centrWord(strconv.Itoa(val.BigShopping), 2) + "|"
		respond.Text += centrWord(strconv.Itoa(val.RoomCleaning), 2) + "|"
		respond.Text += centrWord(strconv.Itoa(val.FlatCleaning), 2) + "|"
		respond.Text += centrWord(strconv.Itoa(val.MoneySpent), 8) + "\n"
		respond.Text += "===========================" + "\n"
		respond.Text += "```"
		respond.Text += "\n"
	}
	respond.ParseMode = "MarkdownV2"
}

func delData(userName, action, amount string) (string, error) {
	if studMap[userName] == nil {
		return "Incorrect username", errors.New("Incorrect username")
	}
	num, err := strconv.Atoi(amount)
	if err != nil {
		return "Incorrect amount to delete", errors.New("Incorrect amount to delete")
	}
	respondText := ""

	switch action {
	case "shopping":
		studMap[userName].Shopping -= num
		respondText = amount + " successfully deleted from shopping"
	case "big_shopping":
		studMap[userName].BigShopping -= num
		respondText = amount + " successfully deleted from big shopping"
	case "throw_trash":
		studMap[userName].Trash -= num
		respondText = amount + " successfully deleted from trash throws"
	case "room_cleaning":
		studMap[userName].RoomCleaning -= num
		respondText = amount + " successfully deleted from room cleanings"
	case "flat_cleaning":
		studMap[userName].FlatCleaning -= num
		respondText = amount + " successfully deleted from flat cleanings"
	case "money":
		studMap[userName].MoneySpent -= num
		respondText = amount + " successfully deleted from money"
	default:
		return "Incorrect action to delete", errors.New("Incorrect action to delete")
	}
	return respondText, nil
}

func messageHandler(msg tgbotapi.Message, userName string) (tgbotapi.MessageConfig, error) {
	cmd := msg.Command()
	args := msg.CommandArguments()
	respond := tgbotapi.NewMessage(msg.Chat.ID, "")
	if msg.Chat.ID != chatId && msg.Chat.ID != adminCHatId {
		respond.Text = "Bot works only in allowed chats"
		return respond, errors.New("Not allowed chat")
	}
	if cmd == "" {
		respond.Text = "Incorrect command name"
		return respond, errors.New("Incorrect command name")
	}
	if studMap[userName] == nil {
		studMap[userName] = new(student)
	}
	switch cmd {
	case "start":
		respond.ReplyMarkup = replyButtons
		respond.Text = "Hi"
		return respond, nil
	case "shopping":
		studMap[userName].Shopping++
		respond.Text = respondTemplate(&msg, cmd)
	case "big_shopping":
		studMap[userName].BigShopping++
		respond.Text = respondTemplate(&msg, cmd)
	case "throw_trash":
		studMap[userName].Trash++
		respond.Text = respondTemplate(&msg, cmd)
	case "room_cleaning":
		studMap[userName].RoomCleaning++
		respond.Text = respondTemplate(&msg, cmd)
	case "flat_cleaning":
		studMap[userName].FlatCleaning++
		respond.Text = respondTemplate(&msg, cmd)
	case "money":
		money, err := strconv.Atoi(args)
		if err != nil || args == "" {
			respond.Text = "Incorrect money value"
			return respond, errors.New("Incorrect money value")
		}
		studMap[userName].MoneySpent += money
		respond.Text = respondTemplate(&msg, cmd)
	case "list":
		drawTable(&respond)
	case "del":
		if msg.Chat.ID != adminChatId {
			respond.Text = "Only admin can delete"
			return respond, errors.New("Only admin can delete")
		}
		argsArr := strings.Split(args, " ")
		if len(argsArr) != 3 {
			respond.Text = "Del accepts only 3 arguments:\ndel ...user_name ...action ...amount"
			return respond, errors.New("Incorrect args for del")
		}
		respond.Text, _ = delData(argsArr[0], argsArr[1], argsArr[2])
	case "save":
		if msg.Chat.ID != adminChatId {
			respond.Text = "Only admin can save"
			return respond, errors.New("Only admin can save")
		}
		respond.Text = save()
	case "recover":
		if msg.Chat.ID != adminChatId {
			respond.Text = "Only admin can recover"
			return respond, errors.New("Only admin can recover")
		}
		respond.Text = recov()

	default:
		respond.Text = "Incorrect command name"
		return respond, errors.New("Incorrect command name")
	}

	return respond, nil
}

func main() {

	time.AfterFunc(duration(), recurTimer)
	wg.Add(1)
	// do normal task here

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		msg, _ := messageHandler(*update.Message, update.Message.From.UserName)
		//log.Println(err.Error())

		bot.Send(msg)

	}
	wg.Wait()

}
