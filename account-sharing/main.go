package account_sharing

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "MTA5Nzk5MTE5Mjg5NTAzNzUyMg.GLn_yJ.ogHcxxmwuaQ_Ay61MIrScIupysAALQzUAUSwf4", "Bot token")
	flag.Parse()
}

var players = []string{"Solo", "Kieu", "Leo", "28"}
var currentPlayer string
var lastLogin time.Time
var discordID2User = map[string]string{
	"921928977671671909": "Kieu",
	"403220436919517194": "Leo",
	"466605197880328193": "Solo",
	"511880106999021588": "28",
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	currentPlayer = "Solo"
	lastLogin = time.Now()
	fmt.Println("dg = ", dg, " err = ", err, " token = ", Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func goodToSendMsg(channelId string) bool {
	allowedChannels := []string{"1034382495338205204", "1097993812900585544"}
	for _, channel := range allowedChannels {
		if channelId == channel {
			return true
		}
	}
	return false
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !goodToSendMsg(m.ChannelID) {
		return
	}
	// if m.ChannelID != "1097993812900585544" {
	// 	return
	// }
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	fmt.Println("authorID = ", m.Author.ID)
	fmt.Println("channelId = ", m.ChannelID)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "check" {
		s.ChannelMessageSend(m.ChannelID, currentPlayer+" đã đăng nhập vào lúc "+lastLogin.Format("Mon Jan 2 15:04:05 MST 2006"))
		return
	}
	requestedPlayer := discordID2User[m.Author.ID]
	action := m.Content
	if action != "login" && action != "logout" {
		return
	}
	if requestedPlayer == "" {
		s.ChannelMessageSend(m.ChannelID, "Người chơi chưa được đăng kí")
		return
	}
	if action == "logout" {
		if currentPlayer == requestedPlayer {
			s.ChannelMessageSend(m.ChannelID, currentPlayer+" đã đăng xuất vào lúc "+lastLogin.Format("Mon Jan 2 15:04:05 MST 2006"))
			currentPlayer = ""
		} else {
			s.ChannelMessageSend(m.ChannelID, "có đăng nhập đâu mà đòi log out? :pepethink:")
		}
		return
	}
	if requestedPlayer == "Solo" {
		currentPlayer = requestedPlayer
		lastLogin = time.Now()
		s.ChannelMessageSend(m.ChannelID, "Mời ngài Solo vào tryhard, thằng "+currentPlayer+" đang chơi thì kệ nó")
		return
	}

	if currentPlayer == "" {
		currentPlayer = requestedPlayer
		lastLogin = time.Now()
		s.ChannelMessageSend(m.ChannelID, "K có ai onl cả, vào đi :enter: :enter:")
		return
	}
	if currentPlayer == "Solo" {
		s.ChannelMessageSend(m.ChannelID, "Ngài Solo đang tryhard :chingchong:")
		return
	}
	if time.Since(lastLogin) > time.Minute*30 {
		s.ChannelMessageSend(m.ChannelID, currentPlayer+" đã đăng nhập được hơn 30 phút, vào mà sút nó ra")
		currentPlayer = requestedPlayer
		lastLogin = time.Now()
		return
	}
	if time.Since(lastLogin) < time.Minute*30 {
		if currentPlayer == requestedPlayer {
			s.ChannelMessageSend(m.ChannelID, "cho thằng "+currentPlayer+" chơi 30 phút :chongnanh:")
		} else {
			s.ChannelMessageSend(m.ChannelID, "để thằng "+currentPlayer+" chơi thêm lúc nữa, chưa đến 30 phút")
		}
		return
	}

}
