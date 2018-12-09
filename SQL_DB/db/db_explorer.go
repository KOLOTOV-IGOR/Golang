// тут лежит тестовый код
// менять вам может потребоваться только коннект к базе
package main

import (
	"database/sql"
	"fmt"
	"strings"
	"net/http"
	"regexp"
	"strconv"
	"io/ioutil"
	"reflect"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

type Handler struct {
	DB       *sql.DB
	Tables   []string
	nameID   map[string]string
	Types    map[string]map[string]string
	ColNames map[string][]string
	Nullable map[string]map[string]string
}

var (
	re = regexp.MustCompile(`\/(\w+)\/*(\d*)$`)
)


func Convert(req *http.Request) (string, string) {
	lim  := req.FormValue("limit")
	off := req.FormValue("offset")
	_, err1 := strconv.Atoi(lim)
	_, err2 := strconv.Atoi(off)
	var limit  string = "5"
	var offset string = "0"
	if lim != "" && err1 == nil {
		limit = lim
	}
	if off != "" && err2 == nil {
		offset = off
	}
	return limit, offset
}

func checkTheBeing(tables []string, table string) bool {
	isHere := false
	for _, item := range tables {
		if item == table {
			isHere = true
		}
	}
	return isHere
}
///Struct Methods/////////////////////////////////////////////////////////////////////////////////////////
func (h *Handler) CheckTable(resp http.ResponseWriter, table string) bool {
	ResponseMap := map[string]interface{}{}
	if isHere := checkTheBeing(h.Tables, table); !isHere {
		resp.WriteHeader(http.StatusNotFound)
		ResponseMap["error"] = "unknown table"
		js, _ := json.Marshal(ResponseMap)
		resp.Write(js)
		return false
	}
	return true
}

func (h *Handler) FullGetProcessing(resp http.ResponseWriter, req *http.Request, table string) {
	db := h.DB
	colsNames := h.ColNames[table]
	vals := make([]interface{}, len(colsNames))
	Records  := map[string]interface{}{}
	FullResp := map[string]interface{}{}
	limit, offset := Convert(req)
//	fmt.Println(limit, offset)

	rows, err := db.Query("select * from " + table + " limit " + offset + ", " + limit)
	if err != nil {
		fmt.Println("Error!!!", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	ListResp := make([]map[string]interface{}, 0, len(colsNames))
	for i, _ := range colsNames {
		vals[i] = new(sql.RawBytes)
	}
	for rows.Next() {
		respMap  := map[string]interface{}{}
		err = rows.Scan(vals...)
		temp, _ := rows.ColumnTypes()
		for i, item := range vals {
		//	fmt.Printf("%s\n", temp[i].ScanType().Name()+" "+temp[i].Name())//DatabaseTypeName())
			stor := *item.(*sql.RawBytes)
			switch temp[i].ScanType().Name() {
			case "int32":
				t, err := strconv.Atoi(string(stor))
				respMap[temp[i].Name()] = t
				if err != nil {
					return
				}
			default:
				if stor == nil {
					respMap[temp[i].Name()] = nil
				} else {
					respMap[temp[i].Name()] = string(stor)
				}
			}
		}
		ListResp = append(ListResp, respMap)
	}
	//Preparation to send 
	Records["records"]   = ListResp
	FullResp["response"] = Records
	resp.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(FullResp)
	resp.Write(js)
	rows.Close()
	return
}

func (h *Handler) PartGetProcessing(resp http.ResponseWriter, req *http.Request, table, id string) {
	db := h.DB
	colsNames := h.ColNames[table]
	ResponseMap := map[string]interface{}{}
	FullResp := map[string]interface{}{}

	rows, err := db.Query("select * from " + table + " where " + h.nameID[table] + "=" + id)
	if err != nil {
		fmt.Println(err)
		return
	}
	vals := make([]interface{}, len(colsNames))
	Records  := map[string]interface{}{}
	for i, _ := range colsNames {
		vals[i] = new(sql.RawBytes)
	}
	respMap  := map[string]interface{}{}
	for rows.Next() {
		err = rows.Scan(vals...)
		temp, _ := rows.ColumnTypes()
		for i, item := range vals {
			stor := *item.(*sql.RawBytes)
			switch temp[i].ScanType().Name() {
			case "int32":
				t, err := strconv.Atoi(string(stor))
				respMap[temp[i].Name()] = t
				if err != nil {
					return
				}
			default:
				if stor == nil {
					respMap[temp[i].Name()] = nil
				} else {
					respMap[temp[i].Name()] = string(stor)
				}
			}
		}
	}
	//Preparation to send 
	if len(respMap) == 0 {
		resp.WriteHeader(http.StatusNotFound)
		ResponseMap["error"] = "record not found"
		js, _ := json.Marshal(ResponseMap)
		resp.Write(js)
		return
	}
	Records["record"]   = respMap
	FullResp["response"] = Records
	resp.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(FullResp)
	resp.Write(js)
	rows.Close()
	return
}

func (h * Handler) PutProcessing(resp http.ResponseWriter, req *http.Request, table string) {
	b, err := ioutil.ReadAll(req.Body)
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(m)
	ResponseMap := map[string]interface{}{}
	Records     := map[string]interface{}{}

	tmpP := make([]string, 0, len(m))
	tmpV := make([]interface{}, 0, len(m))
	ignored := ""
	for k, v := range m {
		if strings.Contains(k, "id") {
			ignored = k
			continue
		}
		isHere := false
		for _, it := range h.ColNames[table] {
			if it == k {
				isHere = true
			}
		}
		if !isHere {
			continue
		}
		tmpP = append(tmpP, k)
		tmpV = append(tmpV, v)
	}

	params := " ("+strings.Join(tmpP, ",")+") "
	values := " values("
	for _, _ = range tmpP {
		values += "?,"
	}
	values = values[:len(values)-1]
	stmt, _ := h.DB.Prepare("insert into " + table + params + values + ")")
	res, err := stmt.Exec(tmpV...)
	if err != nil {
		fmt.Println("ERROR!!!", err)
		return
	}
	lastID, _ := res.LastInsertId()

	Records[ignored] = lastID
	ResponseMap["response"] = Records
	js, _ := json.Marshal(ResponseMap)
	resp.Write(js)
	return
}

func (h * Handler) PostProcessing(resp http.ResponseWriter, req *http.Request, table, id string) {
	b, err := ioutil.ReadAll(req.Body)
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		return
	}
	ResponseMap := map[string]interface{}{}
	Records     := map[string]interface{}{}
	nameID      := h.nameID[table]
	nullable    := h.Nullable[table]
	types       := h.Types[table]
	//fmt.Println(nameID, nullable, types)
	tmpP := make([]string, 0, len(m))
	tmpV := make([]interface{}, 0, len(m))
	for k, v := range m {
		if strings.Contains(k, "id") {
			resp.WriteHeader(http.StatusBadRequest)
			respMap := map[string]interface{}{}
			respMap["error"] = "field " + k  +" have invalid type"
			js, _ := json.Marshal(respMap)
			resp.Write(js)
			return
		}
		if m[k] == nil && nullable[k] != "YES" {
			resp.WriteHeader(http.StatusBadRequest)
			respMap := map[string]interface{}{}
			respMap["error"] = "field " + k  +" have invalid type"
			js, _ := json.Marshal(respMap)
			resp.Write(js)
			return
		}
		dtype := reflect.ValueOf(v).Kind()
//		fmt.Println("REFL", dtype, "---", dtype.String())
		switch dtype.String() {
		case "string":
			if strings.Contains(types[k], "int") {
				resp.WriteHeader(http.StatusBadRequest)
				respMap := map[string]interface{}{}
				respMap["error"] = "field " + k  +" have invalid type"
				js, _ := json.Marshal(respMap)
				resp.Write(js)
				return
			}
			tmpP = append(tmpP, k)
			tmpV = append(tmpV, v)
		case "invalid":
			tmpP = append(tmpP, k)
			tmpV = append(tmpV, v)
		case "float64":
			if strings.Contains(types[k], "varchar") || types[k] == "text" {
				resp.WriteHeader(http.StatusBadRequest)
				respMap := map[string]interface{}{}
				respMap["error"] = "field " + k  +" have invalid type"
				js, _ := json.Marshal(respMap)
				resp.Write(js)
				return
			}
			tmpP = append(tmpP, k)
			tmpV = append(tmpV, v)
		}
	}

	//fmt.Printf("%#v\n", tmpV)
	values := ""
	for _, k := range tmpP {
		values = values + k + " = ?, "
	}
	values = values[:len(values)-2]
	//fmt.Println("update " + table + " set " + values + " where " + nameID + " = " + id)
	stmt, _ := h.DB.Prepare("update " + table + " set " + values + " where " + nameID + " = " + id)
	res, err := stmt.Exec(tmpV...)
	if err != nil {
		fmt.Println(err)
		return
	}
	r, _ := res.RowsAffected()
	Records["updated"] = r

	resp.WriteHeader(http.StatusOK)
	ResponseMap["response"] = Records
	js, _ := json.Marshal(ResponseMap)
	resp.Write(js)
	return
}

func (h *Handler) DeleteProcessing(resp http.ResponseWriter, req *http.Request, table, id string) {
	b, err := ioutil.ReadAll(req.Body)
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		return
	}

	nameID := h.nameID[table]
	res, err := h.DB.Exec("delete from " + table + " where " + nameID + " = ?", id)
	if err != nil {
		fmt.Println("ERROR!!!", err)
		return
	}

	resp.WriteHeader(http.StatusOK)
	respMap := map[string]interface{}{}
	rowaff, _ := res.RowsAffected()
	respMap["deleted"] = rowaff
	FullResp := map[string]interface{}{}
	FullResp["response"] = respMap
	js, _ := json.Marshal(FullResp)
	resp.Write(js)
}

/////The main method and function/////////////////////////////////////////////////////
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	group := re.FindStringSubmatch(req.URL.Path)
	table := ""
	id    := ""
	if group != nil {
		table = group[1]
		id    = group[2]
	}

	if req.URL.Path == "/" {
		ResponseMap := map[string]interface{}{}
		temp := map[string]interface{}{}
		resp.WriteHeader(http.StatusOK)
		temp["tables"] = h.Tables
		ResponseMap["response"] = temp
		js, _ := json.Marshal(ResponseMap)
		resp.Write(js)
		return
	}

	if isHere := h.CheckTable(resp, table); !isHere {//Is the needed table here?
		return
	}

	if id == "" {
		switch req.Method {
		case http.MethodGet:
			h.FullGetProcessing(resp, req, table)
			return
		case http.MethodPut:
			h.PutProcessing(resp, req, table)
			return
		}
	} else {
		switch req.Method {
		case http.MethodGet:
			h.PartGetProcessing(resp, req, table, id)
			return
		case http.MethodPost:
			h.PostProcessing(resp, req, table, id)
			return
		case http.MethodDelete:
			h.DeleteProcessing(resp, req, table, id)
			return
		}
	}
}

func NewDbExplorer(db *sql.DB) (*Handler, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		fmt.Println("SHOW ERROR: ", err)
		return nil, err
	}
	tables := []string{}
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		tables = append(tables, s)
	}
	rows.Close()

	colsNames := map[string][]string{}
	IDnames   := map[string]string{}
	Nullable  := map[string]map[string]string{}
	Types     := map[string]map[string]string{}
	for _, item := range tables {
		r, _ := db.Query("show full columns from " + item)
		c, _ := r.Columns()
		vals := make([]interface{}, len(c))
		for i, _ := range c {
			vals[i] = new(sql.RawBytes)
		}
		columns  := []string{}
		nullable := map[string]string{}
		types := map[string]string{}
		for r.Next() {
			r.Scan(vals...)
			columns = append(columns, string(*vals[0].(*sql.RawBytes)))
			//fmt.Printf("%s %s\n", string(*vals[0].(*sql.RawBytes)), string(*vals[3].(*sql.RawBytes)))
			nullable[string(*vals[0].(*sql.RawBytes))] = string(*vals[3].(*sql.RawBytes))
			types[string(*vals[0].(*sql.RawBytes))] = string(*vals[1].(*sql.RawBytes))
		}
		//fmt.Printf("%s\n", string(*vals[1].(*sql.RawBytes)))
		//fmt.Println(columns)
		colsNames[item] = columns
		Nullable[item]  = nullable
		Types[item]     = types
		for _, it := range colsNames[item] {
			if strings.Contains(it, "id") {
				IDnames[item] = it
			}
		}
		r.Close()
	}

	//fmt.Println(Nullable)
	handler := &Handler{
		DB:       db,
		Tables:   tables,
		nameID:   IDnames,
		Types:    Types,
		ColNames: colsNames,
		Nullable: Nullable,
	}

	return handler, nil
}

