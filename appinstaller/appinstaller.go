package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	//"os/exec"
	"time"
	"log"

	"github.com/zyxar/argo/rpc"
	
	"database/sql"
	_"github.com/mattn/go-sqlite3"
)

var (
	rpcc               rpc.Client
	rpcSecret          string
	rpcURI             string
	launchLocal        bool
	errParameter       = errors.New("invalid parameter")
	errNotSupportedCmd = errors.New("not supported command")
	errInvalidCmd      = errors.New("invalid command")
)

func init() {
	flag.StringVar(&rpcSecret, "secret", "", "set --rpc-secret for aria2c")
	flag.StringVar(&rpcURI, "uri", "ws://localhost:6800/jsonrpc", "set rpc address")
	flag.BoolVar(&launchLocal, "launch", false, "launch local aria2c daemon")
}

func InitSql() {
	db, err := sql.Open("sqlite3", "foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS tasks (
                                        id integer PRIMARY KEY,
                                        gid text NOT NULL,
                                        title text NOT NULL,
                                        file  text NOT NULL,
                                        type  text NOT NULL,
                                        status text,
                                        totalLength text,
                                        completedLength text,
                                        fav text
                                    );
                                    `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	sqlStmt = `
	 CREATE TABLE IF NOT EXISTS warehouse (
                                        id integer PRIMARY KEY,
                                        title text NOT NULL,
                                        file  text NOT NULL,
                                        type  text NOT NULL
                                        );
                                        `
	_,err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n",err,sqlStmt)
		return
	}
	
	sqlStmt = `
	 CREATE TABLE IF NOT EXISTS wifi      (
                                        id integer PRIMARY KEY,
                                        essid text NOT NULL,
                                        pass  text NOT NULL
                                        );
                                        `
	_,err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n",err,sqlStmt)
		return
	}
	
}

func main() {
	flag.Parse()


	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "usage: app start\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}
	
	InitSql()

	var err error
	rpcc, err = rpc.New(context.Background(), rpcURI, rpcSecret, time.Second, AppNotifier{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer rpcc.Close()
	
	for {
		select{
		}
	}
}

