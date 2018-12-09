package main
import (
	"time"
	"testing"
	"fmt"
	"reflect"
	"net/http"
	"net/http/httptest"
	"strconv"
	"net/url"
)

type TestCase struct {
	Client      *SearchClient
	Request     SearchRequest
	Response    *SearchResponse
	ExpectedErr error
}

func TestXml(t *testing.T) {
	ts := httptest.NewServer(SearchServer("unknown.xml"))

	item := TestCase{
		Client:   &SearchClient{AccessToken: "Hi"},
		Request:  SearchRequest{Limit: 0, Offset: 1, Query: "Ow", OrderField: "", OrderBy: 1},
		Response: nil,//&SearchResponse{Users: []User{}, NextPage: true },
		ExpectedErr: fmt.Errorf("SearhServer fatal error"),
	}
	cl       := item.Client
	cl.URL    = ts.URL
	req      := item.Request
	resp     := item.Response
	tst, err := cl.FindUsers(req)
	if !reflect.DeepEqual(tst, resp) && err != nil {
		t.Errorf("Error!!!")
	}

	ts.Close()
}

func TestBadURL(t *testing.T) {
	ts := httptest.NewServer(SearchServer("dataset.xml"))

	item := TestCase{
		Client:   &SearchClient{AccessToken: "Hi", URL: "lclhost"},
		Request:  SearchRequest{Limit: 0, Offset: 1, Query: "Ow", OrderField: "", OrderBy: 1},
		Response: nil,//&SearchResponse{Users: []User{}, NextPage: true },
		ExpectedErr: fmt.Errorf("unknown error"),
	}
	cl       := item.Client
	req      := item.Request
	resp     := item.Response
	tst, err := cl.FindUsers(req)
	if !reflect.DeepEqual(tst, resp) && err != nil {
		t.Errorf("Error!!!")
	}

	ts.Close()
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		time.Sleep(2*time.Second)
	}))

	item := TestCase{
		Client:   &SearchClient{AccessToken: "Hi", URL: ts.URL},
		Request:  SearchRequest{Limit: 0, Offset: 1, Query: "Ow", OrderField: "", OrderBy: 1},
		Response: nil,//&SearchResponse{Users: []User{}, NextPage: true },
		ExpectedErr: fmt.Errorf("timeout error"),
	}
	cl       := item.Client
	req      := item.Request
	resp     := item.Response
	tst, err := cl.FindUsers(req)
	if !reflect.DeepEqual(tst, resp) && err != nil {
		t.Errorf("Error!!!")
	}

	ts.Close()
}

func TestErrResJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("Yes No"))
	}))

	item := TestCase {
		Client:      &SearchClient{AccessToken: "Hi", URL: ts.URL},
		Request:     SearchRequest{Limit: 0, Offset: 1, Query: "Ow", OrderField: "", OrderBy: 1},
		Response:    nil,
		ExpectedErr: fmt.Errorf("cant unpack result json: "),
	}
	cl       := item.Client
	req      := item.Request
	resp     := item.Response
	tst, err := cl.FindUsers(req)
	if !reflect.DeepEqual(tst, resp) && err != nil {
		t.Errorf("Error!!!")
	}

	ts.Close()
}

func TestErrJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusBadRequest)
		js := []byte(`{`)
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(js)
	}))

	item := TestCase {
		Client:      &SearchClient{AccessToken: "Hi", URL: ts.URL},
		Request:     SearchRequest{Limit: 0, Offset: 1, Query: "Ow", OrderField: "_", OrderBy: 1},
		Response:    nil,
		ExpectedErr: fmt.Errorf("cant unpack error json: "),
	}
	cl       := item.Client
	req      := item.Request
	resp     := item.Response
	tst, err := cl.FindUsers(req)
	if !reflect.DeepEqual(tst, resp) && err != nil {
		t.Errorf("Error!!!")
	}

	ts.Close()
}

func TestUnknownErr(t *testing.T) {
	ts := httptest.NewServer(SearchServer("dataset.xml"))

	item := TestCase {
		Client:      &SearchClient{AccessToken: "Hi", URL: ts.URL},
		Request:     SearchRequest{Limit: 1, Offset: 1, Query: "", OrderField: "", OrderBy: 1},
		Response:    nil,
		ExpectedErr: fmt.Errorf("unknown bad request!"),
	}
	cl       := item.Client
	req      := item.Request
	resp     := item.Response
	tst, err := cl.FindUsers(req)
	if !reflect.DeepEqual(tst, resp) && err != nil {
		t.Errorf("Error!!!")
	}

	ts.Close()
}

func TestAdd1(t *testing.T) {
	ts := httptest.NewServer(SearchServer("dataset.xml"))
		client  = &http.Client{Timeout: time.Second}
		searcherParams := url.Values{}

		searcherParams.Add("limit", "afa")
		searcherParams.Add("offset", "1")
		searcherParams.Add("query", "da")
		searcherParams.Add("order_field", "ad")
		searcherParams.Add("order_by", strconv.Itoa(0))

		searcherReq, _ := http.NewRequest("GET", ts.URL+"?"+searcherParams.Encode(), nil)
		searcherReq.Header.Add("AccessToken", "Hi")

		resp, _ := client.Do(searcherReq)
		
		if resp.StatusCode != http.StatusInternalServerError {
			t.Error(resp.StatusCode, " != ", http.StatusInternalServerError)
		}

	ts.Close()	
}

func TestAdd2(t *testing.T) {
	ts := httptest.NewServer(SearchServer("dataset.xml"))
		client  = &http.Client{Timeout: time.Second}
		searcherParams := url.Values{}

		searcherParams.Add("limit", "1")
		searcherParams.Add("offset", "sfs")
		searcherParams.Add("query", "da")
		searcherParams.Add("order_field", "ad")
		searcherParams.Add("order_by", strconv.Itoa(0))

		searcherReq, _ := http.NewRequest("GET", ts.URL+"?"+searcherParams.Encode(), nil)
		searcherReq.Header.Add("AccessToken", "Hi")

		resp, _ := client.Do(searcherReq)
		
		if resp.StatusCode != http.StatusInternalServerError {
			t.Error(resp.StatusCode, " != ", http.StatusInternalServerError)
		}

	ts.Close()	
}

func TestAdd3(t *testing.T) {
	ts := httptest.NewServer(SearchServer("dataset.xml"))
		client  = &http.Client{Timeout: time.Second}
		searcherParams := url.Values{}

		searcherParams.Add("limit", "1")
		searcherParams.Add("offset", "0")
		searcherParams.Add("query", "da")
		searcherParams.Add("order_field", "ad")
		searcherParams.Add("order_by", "ry")

		searcherReq, _ := http.NewRequest("GET", ts.URL+"?"+searcherParams.Encode(), nil)
		searcherReq.Header.Add("AccessToken", "Hi")

		resp, _ := client.Do(searcherReq)
		
		if resp.StatusCode != http.StatusInternalServerError {
			t.Error(resp.StatusCode, " != ", http.StatusInternalServerError)
		}

	ts.Close()	
}

func TestFindUsers( t *testing.T) {
	cases := []TestCase {
		TestCase{
			Client:   &SearchClient{AccessToken: "Hi"},
			Request:  SearchRequest{Limit: 0, Offset: 1, Query: "Ow", OrderField: "", OrderBy: 1},
			Response: &SearchResponse{Users: []User{}, NextPage: true },
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 4, Offset: 3, Query: "B", OrderField: "", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{22, "Beth", 31, "Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n", "female"}, {19, "Bell", 26, "Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n", "male"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 4, Offset: 3, Query: "B", OrderField: "", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{19, "Bell", 26, "Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n", "male"}, {22, "Beth", 31, "Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n", "female"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 4, Offset: 3, Query: "B", OrderField: "Name", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{22, "Beth", 31, "Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n", "female"}, {19, "Bell", 26, "Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n", "male"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 4, Offset: 3, Query: "B", OrderField: "Name", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{19, "Bell", 26, "Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n", "male"}, {22, "Beth", 31, "Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n", "female"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 5, Offset: 4, Query: "A", OrderField: "ID", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{15, "Allison", 21, "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n", "male"}, {11, "Gilmore", 32, "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n", "male"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 5, Offset: 4, Query: "", OrderField: "ID", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"}, {5, "Beulah", 30, "Enim cillum eu cillum velit labore. In sint esse nulla occaecat voluptate pariatur aliqua aliqua non officia nulla aliqua. Fugiat nostrud irure officia minim cupidatat laborum ad incididunt dolore. Fugiat nostrud eiusmod ex ea nulla commodo. Reprehenderit sint qui anim non ad id adipisicing qui officia Lorem.\n", "female"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 2, Offset: 1, Query: "K", OrderField: "Age", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{34, "Kane", 34, "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n", "male"}, {32, "Christy", 40, "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n", "female"},}, NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 2, Offset: 1, Query: "K", OrderField: "Age", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{32, "Christy", 40, "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n", "female"},{34, "Kane", 34, "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n", "male"},},NextPage: false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:   &SearchClient{AccessToken: "Hi"},
			Request:  SearchRequest{Limit: 10, Offset: 1, Query: "Ow", OrderField: "Name", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:   &SearchClient{AccessToken: "Hi"},
			Request:  SearchRequest{Limit: 10, Offset: 1, Query: "Ow", OrderField: "Age", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:   &SearchClient{AccessToken: "Hi"},
			Request:  SearchRequest{Limit: 10, Offset: 1, Query: "Ow", OrderField: "ID", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:   &SearchClient{AccessToken: "Hi"},
			Request:  SearchRequest{Limit: 10, Offset: 1, Query: "Ow", OrderField: "", OrderBy: 1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: -10, Offset: 1, Query: "Ow", OrderField: "Age", OrderBy: -1},
			Response:    nil,//&SearchResponse{},
			ExpectedErr: fmt.Errorf("limit must be > 0"),
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 110, Offset: 0, Query: "Ow", OrderField: "Age", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 110, Offset: 0, Query: "Ow", OrderField: "ID", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 110, Offset: 0, Query: "Ow", OrderField: "Name", OrderBy: -1},
			Response: &SearchResponse{Users: []User{{4, "Owen", 30, "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", "male"},}, NextPage: false },
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "Hi",},
			Request:     SearchRequest{Limit: 10, Offset: -1, Query: "Ow", OrderField: "Age", OrderBy: -1},
			Response:    nil,//&SearchResponse{},
			ExpectedErr: fmt.Errorf("offset must be > 0"),
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "",},
			Request:     SearchRequest{Limit: 10, Offset: 1, Query: "Ow", OrderField: "Age", OrderBy: -1},
			Response:    nil,//&SearchResponse{},
			ExpectedErr: fmt.Errorf("Bad AccessToken"),
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "hi",},
			Request:     SearchRequest{Limit: 10, Offset: 1, Query: "Ow", OrderField: "5", OrderBy: -1},
			Response:    nil,//&SearchResponse{},
			ExpectedErr: fmt.Errorf("ErrorBadOrderField"),
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "hi",},
			Request:     SearchRequest{Limit: 10, Offset: 0, Query: "ZZZZ", OrderField: "ID", OrderBy: 0},
			Response:    &SearchResponse{[]User{}, false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "hi",},
			Request:     SearchRequest{Limit: 10, Offset: 0, Query: "ZZZZ", OrderField: "Age", OrderBy: 0},
			Response:    &SearchResponse{[]User{}, false},
			ExpectedErr: nil,
		},
		TestCase{
			Client:      &SearchClient{AccessToken: "hi",},
			Request:     SearchRequest{Limit: 10, Offset: 0, Query: "ZZZZ", OrderField: "Name", OrderBy: 0},
			Response:    &SearchResponse{[]User{}, false},
			ExpectedErr: nil,
		},
	}

	ts := httptest.NewServer(SearchServer("dataset.xml"))
	for _, item := range cases {
		cl       := item.Client
		cl.URL    = ts.URL
		req      := item.Request
		resp     := item.Response
		tst, err := cl.FindUsers(req)
		if !reflect.DeepEqual(tst, resp) && err == nil {
			t.Error(tst, " doesn't equal ", resp)
		}
	}
	ts.Close()
}


