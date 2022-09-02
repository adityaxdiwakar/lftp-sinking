package sink

import (
	"bytes"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

func GetDirectoryListing(addr, user, pass, dir string) (*[]string, *[]string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, "22"), config)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	err = session.Run("ls -p " + dir)

	files := []string{}
	folders := []string{}

	for _, line := range strings.Split(b.String(), "\n") {
		if strings.HasSuffix(line, "/") {
			folders = append(folders, line[:len(line)-1])
		} else {
			files = append(files, line)
		}
	}

	return &folders, &files, nil
}
