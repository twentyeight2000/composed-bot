package account_sharing

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var players = []string{"Solo", "Kieu", "Leo", "28"}
var currentPlayer string
var lastLogin time.Time
var discordID2User = map[string]string{
	"921928977671671909": "Kieu",
	"403220436919517194": "Leo",
	"466605197880328193": "Solo",
	"511880106999021588": "28",
}

type AccountSharingBot struct {
	db                IAccountSharingDB
	workingChannelIDs []string
	discordSession    *discordgo.Session
	defaultOwnerID    string
}

func NewAccountSharingBot(db IAccountSharingDB, workingChannelIDs []string) *AccountSharingBot {
	token := os.Getenv("BOT_TOKEN")
	discordSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	return &AccountSharingBot{
		db:                db,
		workingChannelIDs: workingChannelIDs,
		discordSession:    discordSession,
		defaultOwnerID:    "511880106999021588",
	}
}

func (a *AccountSharingBot) inWorkingChannel(channelID string) bool {
	inWorkingChannel := false
	for _, cid := range a.workingChannelIDs {
		if cid == channelID {
			inWorkingChannel = true
		}
	}
	return inWorkingChannel
}

func (a *AccountSharingBot) Login(requestedUser, ownerID, channelID string) error {
	if !a.inWorkingChannel(channelID) {
		return errors.New("not woring channel")
	}
	onlinePlayerID, shouldLogin, err := a.db.GetLoginInfo(ownerID)
	if err != nil {
		return err
	}
	if onlinePlayerID == "" {
		err := a.db.SetOnlinePlayer(ownerID, requestedUser)
		if err != nil {
			return err
		}
		a.discordSession.ChannelMessageSend(channelID, "K có ai onl cả, vào đi :enter: :enter:")
		return nil
	}
	onlinePlayerName, _ := a.db.GetUserNameFromID(onlinePlayerID)
	if shouldLogin {
		a.db.SetOnlinePlayer(ownerID, requestedUser)
		a.discordSession.ChannelMessageSend(channelID, fmt.Sprintf("Thằng %s online lâu vl, kick nó ra", onlinePlayerName))
	} else {
		a.discordSession.ChannelMessageSend(channelID, fmt.Sprintf("Cho thằng %s chơi thêm lúc nữa", onlinePlayerName))
	}
	return nil
}

func (a *AccountSharingBot) Logout(userID, ownerID, channelID string) error {
	if !a.inWorkingChannel(channelID) {
		return errors.New("not woring channel")
	}
	onlinePlayerID, _, err := a.db.GetLoginInfo(ownerID)
	if err != nil {
		return err
	}
	onlinePlayerName, _ := a.db.GetUserNameFromID(onlinePlayerID)
	if onlinePlayerID != userID {
		a.discordSession.ChannelMessageSend(channelID, fmt.Sprintf("Có login đâu mà đòi logout :pepethink:"))
	} else {
		a.discordSession.ChannelMessageSend(channelID, fmt.Sprintf("%s đã logout", onlinePlayerName))
	}
	return nil
}

func (a *AccountSharingBot) Run() {
	a.discordSession.AddHandler(a.doMessageCreate())

	// In this example, we only care about receiving message events.
	a.discordSession.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := a.discordSession.Open()
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
	a.discordSession.Close()
}

func (a *AccountSharingBot) doMessageCreate() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !a.inWorkingChannel(m.ChannelID) {
			return
		}
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "login" {
			a.Login(m.Author.ID, a.defaultOwnerID, m.ChannelID)
		} else {
			a.Logout(m.Author.ID, a.defaultOwnerID, m.ChannelID)
		}
	}
}
