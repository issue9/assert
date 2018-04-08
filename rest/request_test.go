// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

type bodyTest struct {
	ID int `json:"id" xml:"id"`
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
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := json.Unmarshal(bs, b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b.ID++
			bs, err = json.Marshal(b)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Add("content-Type", "application/json;charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			w.Write(bs)
			return
		}

		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	w.WriteHeader(http.StatusNotFound)
})

func TestRequest_buildPath(t *testing.T) {
	a := assert.New(t)
	srv := NewServer(t, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	req := srv.NewRequest(http.MethodGet, "/get")
	a.NotNil(req)
	a.Equal(req.buildPath(), "/get")

	req.Param("id", "1").Query("page", "5")
	a.Equal(req.buildPath(), "/get?page=5")

	req = srv.NewRequest(http.MethodGet, "/users/{id}/orders/{oid}")
	a.NotNil(req)
	a.Equal(req.buildPath(), "/users/{id}/orders/{oid}")
	req.Param("id", "1").Param("oid", "2").Query("page", "5")
	a.Equal(req.buildPath(), "/users/1/orders/2?page=5")
}
