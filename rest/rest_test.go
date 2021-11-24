// SPDX-License-Identifier: MIT

package rest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/issue9/assert/v2"
)

type bodyTest struct {
	XMLName struct{} `json:"-" xml:"root"`
	ID      int      `json:"id" xml:"id"`
}

var h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/get" {
		w.WriteHeader(http.StatusCreated)
		return
	}

	if r.URL.Path == "/body" {
		if r.Header.Get("content-type") == "application/json" {
			b := &bodyTest{}
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := json.Unmarshal(bs, b); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b.ID++
			bs, err = json.Marshal(b)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("content-Type", "application/json;charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			w.Write(bs)
			return
		}

		if r.Header.Get("content-type") == "application/xml" {
			b := &bodyTest{}
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := xml.Unmarshal(bs, b); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b.ID++
			bs, err = xml.Marshal(b)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("content-Type", "application/xml;charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			w.Write(bs)
			return
		}

		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	w.WriteHeader(http.StatusNotFound)
})

func TestBuildHandler(t *testing.T) {
	a := assert.New(t, false)

	h := BuildHandler(a, 201, "body", map[string]string{"k1": "v1"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	a.Equal(w.Code, 201)
}

var raw = []*struct {
	req, resp string
}{
	{
		req: `GET /get HTTP/1.1
Host: {host}

`,
		resp: `HTTP/1.1 201

`,
	},
	{
		req: `POST /body HTTP/1.1
Host: {host}
Content-Type: application/json
Content-Length: 8

{"id":5}

`,
		resp: `HTTP/1.1 201
Content-Type: application/json;charset=utf-8

{"id":6}

`,
	},
	{
		req: `DELETE /body?page=5 HTTP/1.0
Host: {host}
Content-Type: application/xml
Content-Length: 23

<root><id>6</id></root>

`,
		resp: `HTTP/1.0 201
Content-Type: application/xml;charset=utf-8

<root><id>7</id></root>

`,
	},
}

func TestServer_RawHTTP(t *testing.T) {
	a := assert.New(t, true)
	s := NewServer(a, h, nil)

	for _, item := range raw {
		req := strings.Replace(item.req, "{host}", s.URL(), 1)
		s.RawHTTP(req, item.resp)
	}
}

func TestRawHandler(t *testing.T) {
	a := assert.New(t, true)
	host := "http://localhost:88"

	for _, item := range raw {
		req := strings.Replace(item.req, "{host}", host, 1)
		RawHandler(a, h, req, item.resp)
	}
}
