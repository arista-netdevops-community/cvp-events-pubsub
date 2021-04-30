package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	server   string
	ssl      bool
	verify   bool
	username string
	password string
)

func ReturnToken() {
	server_ip := strings.Split(server, ":")
	authUrl := "https://" + server_ip[0] + "/cvpservice/login/authenticate.do"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !verify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", authUrl, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	var f map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		fmt.Println(err)
	}
	//token := (f["sessionId"].(string))
	//fmt.Println(f["sessionId"].(string))
	//return token
	err = ioutil.WriteFile("../config/token.txt", []byte(f["sessionId"].(string)), 0644)
	if err != nil {
		fmt.Println(err)
	}
	if ssl {
		certFile, err := os.Create("../config/cvp.crt")
		if err != nil {
			fmt.Println(err)
		}
		w := bufio.NewWriter(certFile)
		certs := resp.TLS.PeerCertificates
		for _, cert := range certs {
			err := pem.Encode(w, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			})
			if err != nil {
				fmt.Println(err)
			}
			if err != nil {
				fmt.Println(err)
			}
			w.Flush()
		}
	}

}
