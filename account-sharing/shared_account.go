package account_sharing

import "time"

type SharedAccount struct {
	ownerID                   string
	userIDs                   []string
	boosterIDs                []string
	onlinePlayerID            string
	lastLogin                 time.Time
	defaultUserSessionTime    time.Duration
	defaultBoosterSessionTime time.Duration
}
