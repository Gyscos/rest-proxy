package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type client struct {
	host string
}

func newClient(host string) *client {
	host = ensureHasPort(host, 80)
	return &client{host}
}

func (c *client) Connect(target string) error {
	target = ensureHasPort(target, 6666)

	conn, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}

	// Prepare communication codec
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	// Read our token
	var id string
	err = dec.Decode(&id)
	if err != nil {
		return err
	}
	// And print it
	fmt.Println(id)

	// Now, answer calls on the socket with the host server
	return c.serveHTTP(dec, enc)
}

func (c *client) serveHTTP(dec *gob.Decoder, enc *gob.Encoder) error {
	for {
		var url string
		err := dec.Decode(&url)
		if err != nil {
			return err
		}
		resp, err := http.Get("http://" + c.host + url)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		enc.Encode(body)
	}

	return nil
}
