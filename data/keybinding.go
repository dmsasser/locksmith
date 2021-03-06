package data

import (
	"fmt"
	"errors"
	"encoding/base64"
)


type Fetcher interface {
	Fetch(id ID) (Key, error)
}


/** Where a Key is bound on an account
 */
type BindingLocation string

const (
	FILE                      BindingLocation = "FILE"
	AUTHORIZED_KEYS           BindingLocation = "AUTHORIZED_KEYS"
	AWS_CREDENTIALS           BindingLocation = "CREDENTIALS"
	INSTANCE_ROOT_CREDENTIALS BindingLocation = "INSTANCE ROOT"
)

type KeyBindingImpl struct {
	KeyID ID
	//AccountID ID `json:",omitempty"`
	Location BindingLocation `json:",omitempty"`
	Name     string          `json:",omitempty"`
}

type KeyBinding interface {
	Describe(keylib Fetcher) (s string, key interface{})
	// TODO:  this should move into a speicfic binding type
	GetSshLine(keylib Fetcher) (string, error)
}

// Describe returns a key binding description and the key described
func (k *KeyBindingImpl) Describe(keylib Fetcher) (s string, key interface{}) {
	if k.Name != "" {
		s = k.Name + " = "
	}

	if key, err := keylib.Fetch(k.KeyID); err != nil {
		s = fmt.Sprintf("%s%s", s, "Unknown key "+k.KeyID)
	} else {
		s = fmt.Sprintf("%s%s", s, key)
	}

	return
}

func (k *KeyBindingImpl) GetSshLine(keylib Fetcher) (string, error){
	if key, err := keylib.Fetch(k.KeyID); err != nil {
		return "", err
	} else {
		if sshKey, ok := key.(*SSHKey) ; !ok{
			return "", errors.New(fmt.Sprint("Key ", key, " is not an SSH key"))
		} else {
			Key2 := sshKey.PublicKey.Key
			return fmt.Sprintf("%s %s %s",
				Key2.Type(),
				base64.StdEncoding.EncodeToString(Key2.Marshal()),
				sshKey.Comments.StringArray()[0]), nil
		}
	}
}