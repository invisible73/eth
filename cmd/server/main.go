package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"

	"github.com/invisible73/eth/pkg/services/eth"
	"github.com/invisible73/eth/pkg/services/transactions"
	txrepo "github.com/invisible73/eth/pkg/services/transactions/postgresql"
	"github.com/invisible73/eth/pkg/types"
	_ "github.com/lib/pq"
)

const (
	sendEth = "SendEth"
	getLast = "GetLast"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client, err := eth.NewClient(os.Getenv("GETH_URL"))
	if err != nil {
		log.Fatal(err)
	}

	repo, err := txrepo.New(db)
	if err != nil {
		log.Fatal(err)
	}

	txservice, err := transactions.New(db, client, repo)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", "localhost:"+os.Getenv("LISTEN_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn, txservice)
	}
}

func handleRequest(c net.Conn, txservice transactions.Service) {
	buf, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		c.Write([]byte(err.Error()))
		c.Close()
		return
	}

	s := strings.Split(buf, " ")

	switch strings.TrimSpace(s[0]) {
	case sendEth:
		if len(s) != 4 {
			c.Write([]byte("wrong number of arguments\n"))
			break
		}

		amount := new(big.Int)
		amount, ok := amount.SetString(strings.TrimRight(s[3], "\r\n"), 10)
		if !ok {
			c.Write([]byte("wrong value for amount field\n"))
			break
		}

		tx, err := txservice.Send(types.Transaction{From: s[1], To: s[2], Amount: amount})
		if err != nil {
			c.Write([]byte(fmt.Sprintf("error while sending transaction: %s\n", err)))
			break
		}

		log.Printf("transaction: %s\n", tx)

		d, _ := json.Marshal(tx)

		c.Write(d)
	case getLast:
		data, err := txservice.GetLast()
		if err != nil {
			c.Write([]byte(fmt.Sprintf("error: %s\n", err)))
		}

		d, err := json.Marshal(data)
		if err != nil {
			c.Write([]byte(fmt.Sprintf("error: %s\n", err)))
		}

		c.Write(d)
	default:
		c.Write([]byte(fmt.Sprintf("undefined command: %s\n", s[0])))
	}

	c.Close()
}
