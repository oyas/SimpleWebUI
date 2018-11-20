package functions

import (
	"strings"
	"strconv"
)

type SharpFunc struct {

}

// 外部参照可能な名前(1文字目が大文字)
func (sf SharpFunc) Func1(args []string, body string) string {
	return "<p>" + "===" + body + "</p>"
}

func (sf SharpFunc) Func2(args []string, body string) string {
	pre := args[0]
	ret := "<p>"
	for _, v := range strings.Split(body, "\n") {
		ret += pre + " > " + v + "<br>"
	}
	ret += "</p>"
	return ret
}

func (sf SharpFunc) Func3(args []string, body string) string {
	num, _ := strconv.Atoi(args[0])
	ret := "<p>"
	for _, v := range strings.Split(body, "\n") {
		ret += strconv.Itoa(num) + " > " + v + "<br>"
		num += 1
	}
	ret += "</p>"
	return ret
}

func (sf SharpFunc) Input(args []string, body string) string {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	return  "<label>" + body + " " +
			"<input type=\"text\" class=\"form-control\" name=\"" + name + "\">" +
			"</label>"
}

func (sf SharpFunc) Submit(args []string, body string) string {
	name := "submit"
	if body != "" {
		name = body
	}
	ret := ""
	ret += "<button type=\"submit\" class=\"btn btn-default\">" + name + "</button>"
	return ret
}
