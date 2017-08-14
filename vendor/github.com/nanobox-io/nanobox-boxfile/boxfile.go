// Package boxfile contains logic for working with the nanobox boxfile
package boxfile

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Struct Boxfile contains information about the boxfile
type Boxfile struct {
	Raw    []byte                 // Raw bytes to be parsed into Boxfile.Parsed
	Parsed map[string]interface{} // Parsed is the parsed boxfile.yml
	Valid  bool                   // Valid is true if the boxfile.yml is valid
}

// New returns a boxfile object from raw data
func New(raw []byte) Boxfile {
	box := Boxfile{
		Raw:    raw,
		Parsed: make(map[string]interface{}),
	}
	box.parse()
	return box
}

// NewFromPath creates a new boxfile from a file instead of raw bytes
//
// Deprecated: Use NewFromFile instead, which diferentiates between a missing
// or invalid boxfile
func NewFromPath(path string) Boxfile {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		// fmt.Println("FAILED TO READ BOXFILE")
		return Boxfile{}
	}
	return New(raw)
}

// NewFromFile reads a file and returns a boxfile object. It differentiates
// between a missing boxfile.yml(err) or an invalid one(!.Valid)
func NewFromFile(path string) (*Boxfile, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read boxfile - %s", err.Error())
	}
	box := New(raw)
	return &box, err
}

func (self Boxfile) SaveToPath(path string) error {
	return ioutil.WriteFile(path, self.Raw, 0755)
}

// Node returns just a specific node from the boxfile. If the object is a sub hash,
// it returns a boxfile object. This allows Node to be chained if you know the data.
func (self Boxfile) Node(name string) (box Boxfile) {
	switch self.Parsed[name].(type) {
	case map[string]interface{}:
		box.Parsed = self.Parsed[name].(map[string]interface{})
		box.fillRaw()
		box.Valid = true
	case map[interface{}]interface{}:
		box.Parsed = convertMap(self.Parsed[name].(map[interface{}]interface{}))
		box.Valid = true
	default:
		box.Parsed = make(map[string]interface{})
		box.Valid = false
	}
	return
}

// Nodes allows the user to specify which types of nodes interested in
func (b Boxfile) Nodes(types ...string) (rtn []string) {
	if len(types) == 0 {
		for key, _ := range b.Parsed {
			rtn = append(rtn, key)
		}
		return
	}

	for i := range types {
		for key := range b.Parsed {
			nodeType := regexp.MustCompile(`\..+`).ReplaceAllString(key, "")
			switch types[i] {
			case "container":
				if nodeType == "web" ||
					nodeType == "worker" ||
					nodeType == "data" {
					rtn = append(rtn, key)
				}
			case "code":
				if nodeType == "web" || nodeType == "worker" {
					rtn = append(rtn, key)
				}
			case "web", "worker", "data":
				if nodeType == types[i] {
					rtn = append(rtn, key)
				}
			default:
				if key == types[i] {
					rtn = append(rtn, key)
				}
			}
		}
	}
	return
}

// String returns the raw boxfile as a string
func (self Boxfile) String() string {
	return string(self.Raw)
}

// Value returns the value of a boxfile node
func (b Boxfile) Value(name string) interface{} {
	return b.Parsed[name]
}

func (b Boxfile) StringSliceValue(name string) []string {
	rtn := []string{}
	switch b.Parsed[name].(type) {
	default:
		return []string{}
	case []string:
		return b.Parsed[name].([]string)
	case string:
		return strings.Split(b.Parsed[name].(string), ",")
	case interface{}:
		val, ok := b.Parsed[name].([]interface{})
		if ok {
			for _, key := range val {
				str, ok := key.(string)
				if ok {
					rtn = append(rtn, str)
				}
			}
		}
	case []interface{}:
		for _, key := range b.Parsed[name].([]interface{}) {
			str, ok := key.(string)
			if ok {
				rtn = append(rtn, str)
			}
		}
	}
	return rtn
}

func (b Boxfile) StringValue(name string) string {
	switch b.Parsed[name].(type) {
	default:
		return ""
	case string:
		return b.Parsed[name].(string)
	case bool:
		return strconv.FormatBool(b.Parsed[name].(bool))
	case int:
		return strconv.Itoa(b.Parsed[name].(int))
	case float32:
		return strconv.FormatFloat(b.Parsed[name].(float64), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(b.Parsed[name].(float64), 'f', -1, 64)
	}
}

func (b Boxfile) VersionValue(name string) string {
	switch b.Parsed[name].(type) {
	default:
		return ""
	case string:
		return b.Parsed[name].(string)
	case int:
		return strconv.Itoa(b.Parsed[name].(int)) + ".0"
	case float32:
		return strconv.FormatFloat(b.Parsed[name].(float64), 'f', 1, 32)
	case float64:
		return strconv.FormatFloat(b.Parsed[name].(float64), 'f', 1, 64)
	}
}

func (b Boxfile) IntValue(name string) int {
	switch b.Parsed[name].(type) {
	default:
		return 0
	case string:
		i, _ := strconv.Atoi(b.Parsed[name].(string))
		return i
	case bool:
		if b.Parsed[name].(bool) == true {
			return 1
		}
		return 0
	case int:
		return b.Parsed[name].(int)
	}
}

func (b Boxfile) BoolValue(name string) bool {
	switch b.Parsed[name].(type) {
	default:
		return false
	case string:
		boo, _ := strconv.ParseBool(b.Parsed[name].(string))
		return boo
	case bool:
		return b.Parsed[name].(bool)
	case int:
		return (b.Parsed[name].(int) != 0)
	}
}

// Merge puts a new boxfile data ontop of your existing boxfile
func (self *Boxfile) Merge(box Boxfile) {
	for key, val := range box.Parsed {
		switch self.Parsed[key].(type) {
		case map[string]interface{}, map[interface{}]interface{}:
			sub := self.Node(key)
			sub.Merge(box.Node(key))
			self.Parsed[key] = sub.Parsed
		default:
			self.Parsed[key] = val
		}
	}
}

// MergeProc drops a procfile into the existing boxfile
func (self *Boxfile) MergeProc(box Boxfile) {
	for key, val := range box.Parsed {
		self.Parsed[key] = map[string]interface{}{"exec": val}
	}
}

// Adds any missing storage nodes that are implied in the web => network_dirs but not
// explicitly placed inside the root as a nfs node
func (self *Boxfile) AddStorageNode() {
	for _, node := range self.Nodes() {
		name := regexp.MustCompile(`\d+`).ReplaceAllString(node, "")
		if (name == "web" || name == "worker") && self.Node(node).Value("network_dirs") != nil {
			found := false
			for _, storage := range self.Node(node).Node("network_dirs").Nodes() {
				found = true
				if !self.Node(storage).Valid {
					self.Parsed[storage] = map[string]interface{}{"found": true}
				}
			}

			// if i dont find anything but they did have a network_dirs.. just try adding a new one
			if !found {
				if !self.Node("nfs1").Valid {
					self.Parsed["nfs1"] = map[string]interface{}{"found": true}
				}
			}
		}
	}
}

func (self Boxfile) Equal(other Boxfile) bool {
	return string(self.Raw) == string(other.Raw)
}

// fillRaw is used when a boxfile is create from an existing boxfile and we want to
// see what the raw would look like
func (b *Boxfile) fillRaw() {
	b.Raw, _ = yaml.Marshal(b.Parsed)
}

// parse takes raw data and converts it to a map structure
func (b *Boxfile) parse() {
	// if im given no data it is not a valid boxfile
	if len(b.Raw) == 0 {
		b.Valid = false
		return
	}

	err := yaml.Unmarshal(b.Raw, &b.Parsed)
	if err != nil {
		b.Valid = false
	} else {
		b.Valid = true
	}
	b.ensureValid()
}

// ensureValid ensures the parsed nodes contain valid maps
func (b *Boxfile) ensureValid() {
	if b.Valid {
		for _, node := range b.Nodes() {
			box := b.Node(node)
			if box.Valid {
				box.ensureValid()
				// b.Parsed["image"] = map["image":"postgresql:9.6"]
				b.Parsed[node] = box.Parsed
			}
		}
	}
}

// recursive function that converts map[interface{}]interface
// to a map[string]interface{}
// this is required since json can parse something into the first
// but then cannot put it back into json
func convertMap(in map[interface{}]interface{}) map[string]interface{} {
	rtn := make(map[string]interface{})
	for key, val := range in {
		// if the value is a map of interface interfaces
		// make sure to convert them as well
		var newValue interface{}
		switch val.(type) {
		case []interface{}:
			newValue = convertArray(val.([]interface{}))
		case map[interface{}]interface{}:
			newValue = convertMap(val.(map[interface{}]interface{}))
		default:
			newValue = val
		}

		switch key.(type) {
		case string:
			rtn[key.(string)] = newValue
		}
	}

	return rtn
}

// convert any sub values in an array that may be a map of interface interfaces
func convertArray(in []interface{}) []interface{} {
	rtn := []interface{}{}
	for _, val := range in {
		var newValue interface{}
		switch val.(type) {
		case []interface{}:
			newValue = convertArray(val.([]interface{}))
		case map[interface{}]interface{}:
			newValue = convertMap(val.(map[interface{}]interface{}))
		default:
			newValue = val
		}
		rtn = append(rtn, newValue)
	}
	return rtn
}
