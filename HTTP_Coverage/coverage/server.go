package main

import (
	"net/http"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"encoding/xml"
	"encoding/json"
	"sort"
)

type Users struct {
	List    []Person `xml:"row"`
}

type Person struct {
	ID      int    `xml:"id"`
	Age     int    `xml:"age"`
	Name    string `xml:"first_name"`
	Surname string `xml:"last_name"`
	About   string `xml:"about"`
	Gender  string `xml:"gender"`
}

type ErrorMs struct {
	Error string
}

func Sorting(ordFd string, howSort int, temp []Person, resp http.ResponseWriter) {
	switch ordFd {
	case "":
		fallthrough
	case "Name":
		switch howSort {
		case 1: sort.Slice(temp, func(i,j int) bool {
				NameI := temp[i].Name + temp[i].Surname
				NameJ := temp[j].Name + temp[j].Surname
				return NameI < NameJ
			})
		case -1: sort.Slice(temp, func(i,j int) bool {
				NameI := temp[i].Name + temp[i].Surname
				NameJ := temp[j].Name + temp[j].Surname
				return NameI > NameJ
			})
		case 0:
			break
		}
	case "ID":
		switch howSort {
		case 1: sort.Slice(temp, func(i,j int) bool {
				return temp[i].ID < temp[j].ID
			})
		case -1: sort.Slice(temp, func(i,j int) bool {
				return temp[i].ID > temp[j].ID
			})
		case 0:
			break
		}
	case "Age":
		switch howSort {
		case 1: sort.Slice(temp, func(i,j int) bool {
				return temp[i].Age < temp[j].Age
			})
		case -1: sort.Slice(temp, func(i,j int) bool {
				return temp[i].Age > temp[j].Age
			})
		case 0:
			break
		}
	
	}
	js, _ := json.Marshal(&temp)
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(js)
}

func SearchServer(filename string) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		xmlFile, err := os.Open(filename)
		byteVal, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		users := new(Users)
		xml.Unmarshal(byteVal, users)
		xmlFile.Close()
		acsTok := req.Header.Get("AccessToken")
		if acsTok == "" {
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}
		limit,  err := strconv.Atoi(req.FormValue("limit"))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		offset, err := strconv.Atoi(req.FormValue("offset"))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}	
		query := req.FormValue("query")
		ordFd := req.FormValue("order_field")
		ordBy, err := strconv.Atoi(req.FormValue("order_by"))
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		if query == ""  && ordFd == "" {
			errJs := &ErrorMs{ Error: "Unknown Error" }
			js, _ := json.Marshal(errJs)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json")
			resp.Write(js)
			return
		}

		if ordFd != "ID" && ordFd != "Age" && ordFd != "Name"  && ordFd != "" {
			errJs := &ErrorMs{ Error: "ErrorBadOrderField" }
			js, _ := json.Marshal(errJs)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json")
			resp.Write(js)
			return
		}

		temp := make([]Person, 0, 150)
		if query == "" {
			temp = users.List
			if len(temp) >= limit && offset < limit {
				temp = temp[offset:limit]
				Sorting(ordFd, 1, temp, resp)
			}
			return
		}

		for _, val := range users.List {
			if strings.Contains(val.Name + " " + val.Surname, query) ||
				strings.Contains(val.About, query) {
				temp = append(temp, val)
			}
		}

		if len(temp) >= limit && offset < limit {
			temp = temp[offset:limit]
		}
		Sorting(ordFd, ordBy, temp, resp)
	})
}



