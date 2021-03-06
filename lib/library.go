package lib

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deweysasser/locksmith/data"
	"github.com/deweysasser/locksmith/output"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
)

var TypeMap = make(map[string]reflect.Type)

func AddType(p reflect.Type) {
	TypeMap[p.Name()] = p
}

type IdStringer interface {
	IdString() string
}

type Deserializer func(string, []byte) (interface{}, error)
type IdFunction func(interface{}) string

type library struct {
	Path         string
	deserializer Deserializer
	idfunc       IdFunction
	cache        map[string]interface{}
	cacheLoaded  bool
	// If we want to make store *NOT* hit disk, then uncomment and implement
	//changes map[string]interface{}
}

type Library interface {
	// Flush all active objects to disk
	Flush() error
	// Store the given data as the ID
	Store(object interface{}) error
	// Get the ID of the given object
	Id(object interface{}) string
	// Fetch the data given by ID from the disk
	Fetch(id string) (interface{}, error)
	// Delete the object with the given ID from the disk
	Delete(id string) error
	// Delete the object given
	DeleteObject(o interface{}) error
	// List the objects
	List() chan interface{}
	// Print the cache, for debugging purposes
	PrintCache()
}

func (l *library) Init(path string, idfunc IdFunction, deserializer Deserializer) {
	l.Path = path
	l.idfunc = idfunc
	l.deserializer = deserializer
}

func (l *library) deserialize(id string, bytes []byte) (interface{}, error) {
	switch {
	case l.deserializer != nil:
		return l.deserializer(id, bytes)
	default:
		o := make(map[string]interface{})
		e := json.Unmarshal(bytes, &o)
		if t, ok := o["Type"]; ok {
			if strT, ok := t.(string); ok { // it's a string
				if p, ok := TypeMap[strT]; ok { // it's in the type map
					no := reflect.New(p).Interface()
					e := json.Unmarshal(bytes, &no)
					return no, e
				} else {
					panic(fmt.Sprint("Type ", t, " not in type map"))
				}
			} else {
				panic(fmt.Sprint("Type is not a string"))
			}
		} else {
			panic(fmt.Sprint("Object for ID '", id, "' has no type"))
		}
		return o, e
	}
}

/* Return the primary identifier for this object
 */
func (l *library) Id(o interface{}) string {
	//fmt.Printf("type is %s\n", reflect.TypeOf(o))

	if i, ok := o.(data.Ider); ok {
		return string(i.Id())
	}

	if i, ok := o.(IdStringer); ok {
		return i.IdString()
	}

	if s, ok := o.(fmt.Stringer); ok {
		return hashString(s.String())
	}
	return hash(toJson(o))
}

/** Return the set of identifiers used by this object.  Each identifier must be unique to this object, but there (obviously) can be many.
 */
func (l *library) ids(o interface{}) chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		if i, ok := o.(data.Identiferser); ok {
			for _, id := range i.Identifiers() {
				c <- string(id)
			}
		} else {
			c <- l.Id(o)
		}
	}()
	return c
}

func hashString(s string) string {
	return hash([]byte(s))
}

func hash(s []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(s))
}

func toJson(o interface{}) []byte {
	bytes, e := json.MarshalIndent(o, " ", " ")
	//check(e)
	if e != nil {
		panic(e)
	}
	return bytes

}

func (l *library) PrintCache() {
	l.Load()
	if l.cache == nil {
		output.Debug("Cache is nil")
	}
	output.Debug("Printing cache")
	for k, v := range l.cache {
		output.Debug(k, "=", v)
	}
}

func (l *library) pathOfObject(o interface{}) string {
	return l.pathOfId(l.Id(o))
}

func (l *library) pathOfId(s string) string {
	return fmt.Sprintf("%s/%s.json", l.Path, sanitize(s))
}

func (l *library) Store(o interface{}) error {
	_, e := os.Stat(l.Path)
	if e != nil {
		e = os.MkdirAll(l.Path, 0777)
		if e != nil {
			return e
		}
	}

	path := l.pathOfObject(o)
	//fmt.Println("Writing to " , path)
	bytes, e := json.MarshalIndent(o, " ", " ")
	if e != nil {
		return e
	}
	if e = ioutil.WriteFile(path, bytes, 0666); e == nil {
		l.addToCache(o)
	} else {
		return errors.New(fmt.Sprint("Error storing ", path))
	}
	return e
}

func (l *library) addToCache(o interface{}) {
	if l.cache == nil {
		l.cache = make(map[string]interface{})
	}
	for id := range l.ids(o) {
		l.cache[id] = o
	}
}

func (l *library) Fetch(id string) (interface{}, error) {
	l.Load()
	if l.cache != nil {
		if v, ok := l.cache[id]; ok {
			return v, nil
		}
	}
	path := l.pathOfId(id)
	return l.fetchFrom(id, path)
}

func sanitize(path string) string {
	re := regexp.MustCompile(`\W+`)
	return re.ReplaceAllString(path, "_")

}

func (l *library) fetchFrom(id, path string) (interface{}, error) {

	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}

	o, e := l.deserialize(id, bytes)

	if e == nil {
		//l.cache[Id] = o
	} else {
		return nil, errors.New(fmt.Sprint("Failed to read key in ", path))
	}

	//fmt.Printf("Read %s\n", o)
	return o, e
}

func (l *library) Flush() error {
	return nil
}

func (l *library) Delete(id string) error {
	path := l.pathOfId(id)

	if o, e := l.fetchFrom(id, path); e == nil {
		for i := range l.ids(o) {
			delete(l.cache, i)
		}
	} else {
		return errors.New(fmt.Sprint("Did not find key on disk at ", path, " to delete"))
	}

	return os.Remove(path)
}

func (l *library) DeleteObject(o interface{}) error {
	id := l.Id(o)
	output.Debug("Deleting object with id ", id)
	return l.Delete(id)
}

func (l *library) Load() {
	if l.cache == nil {
		l.cache = make(map[string]interface{})
	}
	if !l.cacheLoaded {
		l.cacheLoaded = true
		for o := range l.List() {
			l.addToCache(o)
		}
	}
}

func (l *library) List() (c chan interface{}) {
	output.Debug("Listing objects from", l.Path)
	c = make(chan interface{})

	_, e := os.Stat(l.Path)

	if e != nil {
		close(c)
		return
	}

	files, e := ioutil.ReadDir(l.Path)

	//fmt.Println("Reading files in ", lib.Path)

	if e != nil {
		close(c)
		return
	}

	go readFiles(l, files, c)
	return
}

func readFiles(lib *library, files []os.FileInfo, c chan interface{}) {
	defer close(c)

	for _, f := range files {
		path := lib.Path + "/" + f.Name()
		//fmt.Println("Reading from ", path)
		o, e := lib.fetchFrom("", path)

		if e == nil {
			//fmt.Println("Enqueuing ", o)
			c <- o
		}

	}
}
