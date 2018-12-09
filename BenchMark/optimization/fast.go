package main

import (
	"io"
	"strconv"
	"bufio"
	"bytes"
	"os"
	"strings"
)

//easyjson:json
type DataJson struct {
	Browsers []string
	Email    string
	Name     string
}

func MyWriter(buf *bytes.Buffer, i int, name, email string) {
	buf.WriteByte('[')
	buf.WriteString(strconv.Itoa(i))
	buf.WriteByte(']')
	buf.WriteByte(' ')
	buf.WriteString(name)
	buf.WriteByte(' ')
	buf.WriteByte('<')
	buf.WriteString(email)
	buf.WriteByte('>')
	buf.WriteString("\n")
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := bytes.Buffer{}
	//seenBrowsers := []string{}
	seenBrowsers := map[string]bool{}//string{}
	uniqueBrowsers := 0
	buf.WriteString("found users:\n")

	var i int = 0
	user := DataJson{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; ok {
					continue
				}
				seenBrowsers[browser] = true
				uniqueBrowsers++
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; ok {
					continue
				}
				seenBrowsers[browser] = true
				uniqueBrowsers++
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		temp := strings.Split(user.Email, "@")
		email := strings.Join(temp, " [at] ")
		MyWriter(&buf, i, user.Name, email)
		buf.WriteTo(out)
		i++
	}

	buf.WriteString("\nTotal unique browsers ")
	buf.WriteString(strconv.Itoa(uniqueBrowsers))
	buf.WriteString("\n")
	buf.WriteTo(out)
}

