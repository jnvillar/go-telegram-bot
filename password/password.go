package password

import "sync"

type PassWordManager struct {
	masterKeys sync.Map
	passwords  sync.Map
}

func New() *PassWordManager {
	return &PassWordManager{
		masterKeys: sync.Map{},
		passwords:  sync.Map{},
	}
}

func (p *PassWordManager) StoreMasterPassword(userID int, masterPassword string) {
	p.masterKeys.Store(userID, masterPassword)
}

func (p *PassWordManager) StorePassword(userID int, name, password string) {
	userPasswords, _ := p.passwords.LoadOrStore(userID, &sync.Map{})
	userPasswords.(*sync.Map).Store(name, password)
}

func (p *PassWordManager) LoadPassword(userID int, name string) (string, bool) {
	userPasswords, ok := p.passwords.Load(userID)
	if !ok {
		return "", false
	}
	pass, ok := userPasswords.(*sync.Map).Load(name)
	if !ok {
		return "", false
	}
	return pass.(string), true
}

func (p *PassWordManager) GetMasterPassword(userID int) (string, bool) {
	pass, found := p.masterKeys.Load(userID)
	if !found {
		return "", found
	}
	return pass.(string), true
}
