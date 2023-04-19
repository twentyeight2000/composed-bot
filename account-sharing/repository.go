package account_sharing

import (
	"errors"
	"time"
)

var ErrAccountNotExist = errors.New("account doesn't exist")
var ErrAccountExisted = errors.New("account already existed")
var ErrNotSupported = errors.New("function is not supported yet")

var _ IAccountSharingDB = (*InmemAccountSharingDB)(nil)

type IAccountSharingDB interface {
	AddUserID2Name(userID string, userName string) error
	AddAccount(ownerID string, userIDs, boosterIDs []string, defaultUserSessionDuration, defaultBoosterSessionDuration time.Duration) error
	AddAccountOwner(ownerID string) error
	AddAccountUserToOwnerAccount(userID string, ownerID string) error
	AddAccountBoosterToOwnerAccount(boosterID string, ownerID string) error
	SetOnlinePlayer(ownerID string, userID string) error
	ChangeUserSessionDuration(ownerID string, userSessionDuration time.Duration) error
	ChangeBoosterSessionDuration(ownerID string, userSessionDuration time.Duration) error

	CheckIsUser(ownerID string, userID string) bool
	GetUserNameFromID(userID string) (string, error)
	GetLoginInfo(ownerID string) (string, bool, error)
}

type InmemAccountSharingDB struct {
	userID2Name           map[string]string
	ownerID2SharedAccount map[string]*SharedAccount
}

func NewInmemAccountSharingDB() *InmemAccountSharingDB {
	return &InmemAccountSharingDB{
		userID2Name:           make(map[string]string),
		ownerID2SharedAccount: make(map[string]*SharedAccount),
	}
}

func (db *InmemAccountSharingDB) AddUserID2Name(userID string, userName string) error {
	db.userID2Name[userID] = userName
	return nil
}

func (db *InmemAccountSharingDB) GetUserNameFromID(userID string) (string, error) {
	return db.userID2Name[userID], nil
}

func (db *InmemAccountSharingDB) CheckIsUser(ownerID string, userID string) bool {
	sharedAccount, ok := db.ownerID2SharedAccount[ownerID]
	if !ok {
		return false
	}
	for _, id := range sharedAccount.boosterIDs {
		if userID == id {
			return true
		}
	}
	for _, id := range sharedAccount.userIDs {
		if userID == id {
			return true
		}
	}
	return false
}

func (db *InmemAccountSharingDB) AddAccount(ownerID string, userIDs, boosterIDs []string, userSessionDuration, boosterSessionDuration time.Duration) error {
	_, ok := db.ownerID2SharedAccount[ownerID]
	if ok {
		return ErrAccountExisted
	}
	account := &SharedAccount{
		ownerID:                   ownerID,
		userIDs:                   userIDs,
		boosterIDs:                boosterIDs,
		onlinePlayerID:            "",
		lastLogin:                 time.Unix(0, 0),
		defaultUserSessionTime:    userSessionDuration,
		defaultBoosterSessionTime: boosterSessionDuration,
	}
	db.ownerID2SharedAccount[ownerID] = account
	return nil
}

func (db *InmemAccountSharingDB) AddAccountOwner(ownerID string) error {
	_, ok := db.ownerID2SharedAccount[ownerID]
	if ok {
		return ErrAccountExisted
	}
	db.ownerID2SharedAccount[ownerID] = &SharedAccount{
		ownerID:                   ownerID,
		userIDs:                   make([]string, 0),
		boosterIDs:                make([]string, 0),
		onlinePlayerID:            "",
		lastLogin:                 time.Unix(0, 0),
		defaultUserSessionTime:    time.Duration(30) * time.Minute,
		defaultBoosterSessionTime: time.Duration(300) * time.Minute,
	}
	return nil
}

func (db InmemAccountSharingDB) SetOnlinePlayer(ownerID, userID string) error {
	sharedAccount, ok := db.ownerID2SharedAccount[ownerID]
	if !ok {
		return ErrAccountNotExist
	}
	sharedAccount.onlinePlayerID = userID
	sharedAccount.lastLogin = time.Now()
	return nil
}

func (db *InmemAccountSharingDB) isBooster(ownerID string, playerID string) bool {
	sharedAccount, ok := db.ownerID2SharedAccount[ownerID]
	if !ok {
		return false
	}
	for _, boosterID := range sharedAccount.boosterIDs {
		if playerID == boosterID {
			return true
		}
	}
	return false
}

func (db *InmemAccountSharingDB) GetLoginInfo(ownerID string) (string, bool, error) {
	sharedAccount, ok := db.ownerID2SharedAccount[ownerID]
	if !ok {
		return "", false, ErrAccountNotExist
	}

	duration := sharedAccount.defaultUserSessionTime
	if db.isBooster(ownerID, sharedAccount.onlinePlayerID) {
		duration = sharedAccount.defaultBoosterSessionTime
	}
	shouldLogin := false
	if time.Since(sharedAccount.lastLogin) > duration {
		shouldLogin = true
	}
	return sharedAccount.onlinePlayerID, shouldLogin, nil
}

func (db *InmemAccountSharingDB) AddAccountUserToOwnerAccount(userID string, ownerID string) error {
	return ErrNotSupported
}

func (db *InmemAccountSharingDB) AddAccountBoosterToOwnerAccount(boosterID string, ownerID string) error {
	return ErrNotSupported
}
func (db *InmemAccountSharingDB) ChangeUserSessionDuration(ownerID string, userSessionDuration time.Duration) error {
	return ErrNotSupported
}
func (db *InmemAccountSharingDB) ChangeBoosterSessionDuration(ownerID string, userSessionDuration time.Duration) error {
	return ErrNotSupported
}
