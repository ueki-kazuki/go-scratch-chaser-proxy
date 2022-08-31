package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/textproto"
)

const (
	GameSet = iota
	UpLeft
	Up
	UpRight
	Left
	Center
	Right
	DownLeft
	Down
	DownRight
)

const (
	Normal = iota
	ENEMY
	BLOCK
	ITEM
)

type Client struct {
	conn    *textproto.Conn
	port    int
	address net.IP
	name    string
	GameSet bool
}

func NewClient(name string, host string, port int) (*Client, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		ip := net.ParseIP(host)
		if ip == nil {
			return nil, errors.New("ParseIP Error")
		}
		ips[0] = ip
	}
	ip := ips[0]
	conn, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := conn.PrintfLine("%v", name); err != nil {
		log.Println(err)
		return nil, err
	}

	return &Client{
		conn:    conn,
		name:    name,
		address: ip.To4(),
		port:    port,
	}, nil
}

func (client *Client) Close() {
	if client.conn != nil {
		if err := client.conn.Close(); err != nil {
			log.Println(err)
		}
	}
}

func (client *Client) Order(command string) (string, error) {
	if err := client.conn.PrintfLine("%v", command); err != nil {
		return "", err
	}
	response, err := client.conn.ReadLine()
	if err != nil {
		return "", err
	}
	switch response[0] {
	case '0':
		log.Println("GameSet!!")
		client.conn.Close()
		client.GameSet = true
		return response, errors.New("GameSet")
	case '1':
		log.Printf("%v\n", response)
	default:
		log.Println("responce error")
		return response, errors.New("responce error")
	}
	if command != "gr" {
		if err := client.conn.PrintfLine(""); err != nil {
			return "", err
		}
	}
	return response, nil
}

func (client *Client) GetReady() (string, error) {
	log.Println("GetReady")
	res, err := client.conn.ReadLine()
	if err != nil {
		return "", err
	}
	if res[0] != '@' {
		log.Println("connection failed")
		client.conn.Close()
		return "", errors.New("connection failed")
	}
	return client.Order("gr")
}
