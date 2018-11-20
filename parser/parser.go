package parser

import (
	"strings"
	"fmt"
	"reflect"
	"./functions"
)

type PaseData struct {
	a int
	b int
	str string
}

func check_char(s string, i int, c byte) bool {
	if len(s) > i {
		return s[i] == c
	}else{
		return false
	}
}

func Parse(str string) string {
	endterm := map[int]string {
		0: "",
		1: "</p>",
	}

	var ret string
	var last int
	var mode int

	func_name := ""
	args := ""
	body := ""

	lines := strings.Split(str, "\n")
	lines = append(lines, "#")
	pos := 0
	for pos < len(lines) {
		v := lines[pos]
		i := pos
		tr := strings.Trim(v, " ")
		tr = strings.Trim(tr, "\t")
		fmt.Println(mode, ";", i, v)

		if mode == 0 {
			func_name = ""
			args = ""
			body = ""

			ret += "\n"

			switch {
			case tr == "$$$$$$$":
				// brank
				ret += endterm[last]
				last = 0

			case check_char(v, 0, '#'):
				for i, c := range v[1:] {
					if c == '(' {
						v = v[i+2:]
						mode = 1
						break
					}
					if c == ' ' || c == '\t' || c == '#' {
						body = v[i+2:]
						break
					}
					func_name += string(c)
				}
				if mode == 0 {
					if func_name == "" {
						ret += v
					}else{
						ret += eval(func_name, "", body)
					}
				}

			default:
				//if last == 0 {
				//	ret += "<p>"
				//	last = 1
				//} else {
				//	ret += "<br>"
				//}
				ret += v
			}
		}
		if mode == 1 {
			idx := strings.Index(v, ")")
			if idx >= 0 {
				args += v[:idx]
				if check_char(v, idx + 1, ':') {
					body += v[idx+2:]
					mode = 2
				}else{
					body = strings.Trim(v[idx+1:], " ")
					ret += eval(func_name, args, body)
					mode = 0
				}
			}else{
				args += v
			}
		}else if mode == 2 {
			body += "\n"
			switch {
			case tr == "" && len(v) < 4:
				body += ""
			case v[0] == '\t':
				body += v[1:]
			case len(v) >= 4 && v[:4] == "    ":
				body += v[4:]
			default:
				body = strings.Trim(body, "\n")
				ret += eval(func_name, args, body)
				mode = 0
				continue
			}
		}

		pos += 1
	}

	ret = ret[:len(ret)-2]
	fmt.Println(ret)

	return ret
}

func eval(func_name string, arg_string string, body string) string {
	args := strings.Split(arg_string, ",")
	for i := range args {
		args[i] = strings.Trim(args[i], " ")
	}
	fmt.Println("eval: ", func_name, args, body)
	ret := ""

	if func_name == "" {
		// h1
		ret += "<h1>" + body + "</h1>"
	}else{
		// Call a Sharp Function
		var sf functions.SharpFunc
		func_name = strings.ToUpper(string(func_name[0])) + func_name[1:]    // make sure that first character is upper case
		function := reflect.ValueOf(&sf).MethodByName(func_name)
		if function.IsValid() {
			out := function.Call([]reflect.Value{reflect.ValueOf(args), reflect.ValueOf(body)})
			if s, ok := out[0].Interface().(string); ok {
				ret += s
			}
		}
	}

	return ret + "\n"
}
