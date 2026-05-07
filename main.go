package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/spf13/afero"
)

const (
	publicHost = "10.0.2.2"
)

type Driver struct {
	root       string
	username   string
	password   string
	listenAddr string
}

func (d *Driver) GetSettings() (*ftpserver.Settings, error) {
	return &ftpserver.Settings{
		ListenAddr: d.listenAddr,
		PublicHost: publicHost,

		// Equivalent to Python range(30000, 30010):
		// ports 30000 through 30009.
		PassiveTransferPortRange: ftpserver.PortRange{
			Start: 30000,
			End:   30009,
		},

		Banner: "XP FTP Server",
	}, nil
}

func (d *Driver) ClientConnected(cc ftpserver.ClientContext) (string, error) {
	slog.Info("client connected", "remote", cc.RemoteAddr())
	return "Welcome to XP FTP Server", nil
}

func (d *Driver) ClientDisconnected(cc ftpserver.ClientContext) {
	slog.Info("client disconnected", "remote", cc.RemoteAddr())
}

func (d *Driver) AuthUser(cc ftpserver.ClientContext, username, password string) (ftpserver.ClientDriver, error) {
	if username != d.username || password != d.password {
		return nil, fmt.Errorf("invalid username or password")
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), d.root)
	return fs, nil
}

func (d *Driver) GetTLSConfig() (*tls.Config, error) {
	return nil, nil
}

func main() {
	username := flag.String("user", "xp", "FTP username")
	password := flag.String("pass", "xp", "FTP password")
	listenAddr := flag.String("addr", "127.0.0.1:2121", "Address to listen on")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-user USER] [-pass PASS] [-addr ADDR] ROOT\n", os.Args[0])
		os.Exit(2)
	}

	root, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(root)
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsDir() {
		log.Fatalf("root is not a directory: %s", root)
	}

	fmt.Println("serving root", root)
	fmt.Println("listening on", *listenAddr)
	fmt.Println()
	fmt.Printf("connect to server at: ftp://%s:%s@%s\n", *username, *password, publicHost+":2121")
	fmt.Println()

	driver := &Driver{
		root:       root,
		username:   *username,
		password:   *password,
		listenAddr: *listenAddr,
	}
	server := ftpserver.NewFtpServer(driver)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
