package misc

import (
	"bufio"
	//"fmt"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"log"
	"os"
	"strings"
)

var wpath_encryption = "/etc/wicd/encryption/templates/"

type CurType struct {
	Type      string
	Fields    [][]string
	Optional  [][]string
	Required  [][]string
	Protected [][]string
	Name      string
}

func NewCurType() *CurType {
	p := &CurType{}

	return p
}

func parse_field_ent(field_line string, field_type string) [][]string {
	if field_type == "" {
		field_type = "require"
	}
	var ret [][]string
	fields := strings.Split(field_line, " ")
	if len(fields)%2 != 0 {
		return ret
	} else {
		var t [][2]string
		for i := 0; i < len(fields); i += 2 {
			t = append(t, [2]string{fields[i], fields[i+1]})
		}

		for _, v := range t {

			if strings.HasPrefix(string(v[0]), "*") || strings.HasPrefix(string(v[1]), "*") == false {
				return ret
			}
			ret = append(ret, []string{string(v[0]), string(v[1])[1:]})
		}

	}
	return ret
}

func parse_ent(line, key string) string {
	line = strings.Replace(line, key, "", -1)
	line = strings.Replace(line, "=", "", -1)
	return strings.TrimSpace(line)
}

func parse_enc_templat(enctype string) *CurType {

	file, err := os.Open(wpath_encryption + enctype)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	cur_type := NewCurType()

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)

		if strings.HasPrefix(line, "name") && cur_type.Name == "" {
			cur_type.Name = parse_ent(line, "name")
		} else if strings.HasPrefix(line, "require") {

			cur_type.Required = parse_field_ent(parse_ent(line, "require"), "")

		} else if strings.HasPrefix(line, "optional") {

			cur_type.Optional = parse_field_ent(parse_ent(line, "optional"), "")

		} else if strings.HasPrefix(line, "protected") {

			cur_type.Protected = parse_field_ent(parse_ent(line, "protected"), "")

		} else if strings.HasPrefix(line, "----") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(cur_type.Name) == 0 {
		cur_type = nil
	}

	if len(cur_type.Required) == 0 {
		cur_type = nil
	}

	return cur_type

}

func LoadEncryptionMethods(wired bool) []*CurType {
	active_fname := "active"

	if wired == true {
		active_fname = "active_wired"
	} else {
		active_fname = "active"
	}

	enctypes, _ := UI.ReadLines(wpath_encryption + active_fname)

	var encryptionTypes []*CurType

	for _, v := range enctypes {

		c := parse_enc_templat(v)
    c.Type = v
		encryptionTypes = append(encryptionTypes, c)
	}

	return encryptionTypes

}

