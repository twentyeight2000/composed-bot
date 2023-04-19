package main

import account_sharing "github.com/twentyeight2000/composed-bot.git/account-sharing"

func main() {
	accountSharingBot := account_sharing.NewAccountSharingBot(
		account_sharing.NewInmemAccountSharingDB(),
		[]string{"1097993812376309884", "911525779488256040"},
	)
	accountSharingBot.Run()
}
