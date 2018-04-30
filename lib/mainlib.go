package lib

import (
	"github.com/deweysasser/locksmith/connection"
	"github.com/deweysasser/locksmith/data"
	"reflect"
)

type MainLibrary struct {
	Path        string
	connections Library
	keys        Library
	accounts    Library
}

func init() {
	AddType(reflect.TypeOf(data.SSHKey{}))
	AddType(reflect.TypeOf(data.AWSKey{}))
	AddType(reflect.TypeOf(connection.SSHHostConnection{}))
	AddType(reflect.TypeOf(connection.FileConnection{}))
	AddType(reflect.TypeOf(connection.AWSConnection{}))
	AddType(reflect.TypeOf(data.SSHAccount{}))
	AddType(reflect.TypeOf(data.AWSAccount{}))
	AddType(reflect.TypeOf(data.AWSInstanceAccount{}))
	AddType(reflect.TypeOf(data.AWSIamAccount{}))
}

func (l *MainLibrary) Connections() Library {
	if l.connections == nil {
		clib := new(library)
		clib.Init(l.Path+"/connections", nil, nil)
		l.connections = clib
	}

	return l.connections
}

func (l *MainLibrary) Keys() Library {
	if l.keys == nil {
		klib := new(library)
		klib.Init(l.Path+"/keys", keyid, nil)
		l.keys = klib
	}

	return l.keys
}

func (l *MainLibrary) Accounts() Library {
	if l.accounts == nil {
		klib := new(library)
		klib.Init(l.Path+"/accounts", nil, nil)
		l.accounts = klib
	}

	return l.accounts
}

type keyReadError struct {
	Path string
}

func (e *keyReadError) Error() string {
	return "Error reading " + e.Path
}

func keyid(key interface{}) string {
	return string(key.(data.Key).Id())
}
