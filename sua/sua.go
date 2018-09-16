// MIT License
//
// Copyright (c) 2018 Pablo Ignacio Lalloni
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sua

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/xmlpath.v2"

	"github.com/lalloni/afip/cuit"
)

var (
	userid   = xmlpath.MustCompile("TODO")
	usercuil = xmlpath.MustCompile("TODO")
)

type User struct {
	ID     string
	CUIL   uint64
	Email  string
	Legajo string
}

type Login struct {
	User           User
	Authentication time.Time
	Expiration     time.Time
	Reception      time.Time
	Groups         []string
	ServiceData    []byte
}

type Config struct {
	AllowUnsigned      bool
	AllowUntrusted     bool
	AllowExpired       bool
	AllowServices      []string
	TrustedCertificate *x509.Certificate
}

func Authenticate(token io.Reader, signature io.Reader, config *Config) (*Login, error) {
	now := time.Now()
	tokenbs, err := base64dec(token)
	if err != nil {
		return nil, errors.Wrap(err, "decoding token")
	}
	signaturebs, err := base64dec(signature)
	if err != nil {
		return nil, errors.Wrap(err, "decoding signature")
	}
	// Sólo se valida firma si viene una o si no se permite ingresar sin firma.
	// Esto permite configuraciones de prueba que permiten ingresar solo con token.
	if signature != nil || !config.AllowUnsigned {
		err = config.TrustedCertificate.CheckSignature(x509.MD5WithRSA, tokenbs, signaturebs)
		if err != nil {
			return nil, errors.Wrap(err, "checking token signature")
		}
	}
	root, err := xmlpath.Parse(bytes.NewReader(tokenbs))
	if err != nil {
		return nil, errors.Wrap(err, "parsing token xml")
	}
	login := Login{Reception: now}
	if err := setstring(&login.User.ID, root, userid); err != nil {
		return nil, errors.Wrap(err, "user id")
	}
	if err := setcuit(&login.User.CUIL, root, usercuil); err != nil {
		return nil, errors.Wrap(err, "user cuit")
	}
	return &login, nil
}

func base64dec(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, r))
}

func setstring(dst *string, node *xmlpath.Node, path *xmlpath.Path) error {
	s, err := getstring(node, path)
	if err != nil {
		return err
	}
	*dst = s
	return nil
}

func setcuit(dst *uint64, node *xmlpath.Node, path *xmlpath.Path) error {
	s, err := getstring(node, path)
	if err != nil {
		return err
	}
	num, err := cuit.Parse(s)
	if err != nil {
		return err
	}
	*dst = num
	return nil
}

func setuint64(dst *uint64, node *xmlpath.Node, path *xmlpath.Path) error {
	s, err := getstring(node, path)
	if err != nil {
		return err
	}
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*dst = num
	return nil
}

func getstring(node *xmlpath.Node, path *xmlpath.Path) (string, error) {
	s, ok := path.String(node)
	if !ok {
		return "", errors.New("not found")
	}
	return s, nil
}
