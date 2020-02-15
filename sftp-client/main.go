package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func main() {

	var user string
	// var password string
	var server string

	flag.StringVar(&user, "user", "root", "")
	// flag.StringVar(&password, "password", "", "")
	flag.StringVar(&server, "server", "10.0.2.9:22", "")
	flag.Parse()

	// if password == "" {
	//     fmt.Println("plz input server password")
	//     reader := bufio.NewReader(os.Stdin)
	//     line, _, err := reader.ReadLine()
	//     if err != nil {
	//         log.Fatal(err)
	//     }
	//     password = string(line)
	// }

	privateKeyPath := os.ExpandEnv("$HOME/.ssh/id_rsa")
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(privateKey)

	hostKeyCallback, err := kh.New(os.ExpandEnv("$HOME/.ssh/known_hosts"))
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	// var hostKey ssh.PublicKey
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// ssh.Password(password),
			ssh.PublicKeys(signer),
		},
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: hostKeyCallback,
	}
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// walk a directory
	w := sftp.Walk("/tmp")
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
	}

	// leave your mark
	f, err := sftp.Create("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("Hello world!")); err != nil {
		log.Fatal(err)
	}

	// check it's there
	fi, err := sftp.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)
}
