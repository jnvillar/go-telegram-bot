package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"

	// "github.com/odysseus/vigenere"
)

type PassWordManager struct {
	masterKeys   sync.Map
	passwords    sync.Map
	lastAccess   sync.Map
	renewMinutes int
}

func New() *PassWordManager {
	return &PassWordManager{
		masterKeys: sync.Map{},
		passwords:  sync.Map{},
		lastAccess: sync.Map{},
	}
}

func (p *PassWordManager) StoreMasterPassword(userID int, masterPassword string) {
	p.masterKeys.Store(userID, masterPassword)
	p.lastAccess.Store(userID, time.Now())
}

func (p *PassWordManager) getMasterPassword(userID int) (string, bool) {
	pass, found := p.masterKeys.Load(userID)
	if !found {
		return "", false
	}
	return pass.(string), true
}

func (p *PassWordManager) mustRenewMasterPassword(userID int) bool {
	v, found := p.lastAccess.Load(userID)
	if !found {
		return true
	}
	lastTime := v.(time.Time)
	lastTime.Add(time.Minute * time.Duration(p.renewMinutes))
	return lastTime.After(time.Now())
}

func (p *PassWordManager) StorePassword(userID int, name, password string) error {
	mustRenew := p.mustRenewMasterPassword(userID)
	if mustRenew {
		return errors.New("Hace falta la contraseña maestra")
	}
	masterPassword, found := p.getMasterPassword(userID)
	if !found {
		return errors.New("Hace falta la contraseña maestra")
	}
	encryptedPassword, err := encrypt(masterPassword, []byte(password))
	if err != nil {
		return err
	}
	userPasswords, _ := p.passwords.LoadOrStore(userID, &sync.Map{})
	userPasswords.(*sync.Map).Store(name, encryptedPassword)
	return nil
}

func (p *PassWordManager) LoadPassword(userID int, name string) (string, bool, error) {
	mustRenew := p.mustRenewMasterPassword(userID)
	if mustRenew {
		return "", false, errors.New("Hace falta la contraseña maestra")
	}
	masterPassword, found := p.getMasterPassword(userID)
	if !found {
		return "", false, errors.New("Hace falta la contraseña maestra")
	}
	userPasswords, ok := p.passwords.Load(userID)
	if !ok {
		return "", false, errors.New("No tenes guardadas contraseñas")
	}
	pass, ok := userPasswords.(*sync.Map).Load(name)
	if !ok {
		return "", false, errors.New("Constraseña no encontrada")
	}
	password, err := decrypt([]byte(masterPassword), pass.([]byte))
	if err != nil {
		return "", false, err
	}
	return string(password), true, nil
}

func encrypt(key string, text []byte) ([]byte, error) {
	for i:=0; i<20; i++ {
		key += key
	}
	byteKey := []byte(key)[0:32]
	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
