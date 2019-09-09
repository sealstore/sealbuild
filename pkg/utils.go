package pkg

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/wonderivan/logger"
	"golang.org/x/crypto/ssh"
)

const oneMBByte = 1024 * 1024

func AddrReformat(host string) string {
	if strings.Index(host, ":") == -1 {
		host = fmt.Sprintf("%s:22", host)
	}
	return host
}

func ReturnCmd(host, cmd string) string {
	session, _ := Connect(User, Passwd, PrivateKeyFile, host)
	defer session.Close()
	b, _ := session.CombinedOutput(cmd)
	return string(b)
}

func GetFileSize(url string) int {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]GetFileSize is error %s", err.Error())
		}
	}()
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return int(resp.ContentLength)
}

func WatchFileSize(host, filename string, size int) {
	t := time.NewTicker(3 * time.Second) //every 3s check file
	defer t.Stop()
	for {
		select {
		case <-t.C:
			length := ReturnCmd(host, "ls -l "+filename+" | awk '{print $5}'")
			length = strings.Replace(length, "\n", "", -1)
			length = strings.Replace(length, "\r", "", -1)
			lengthByte, _ := strconv.Atoi(length)
			if lengthByte == size {
				t.Stop()
			}
			lengthFloat := float64(lengthByte)
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", lengthFloat/oneMBByte), 64)
			logger.Alert("transfer total size is:", value, "MB")
		}
	}
}

//Cmd is
func Cmd(host string, cmd string) []byte {
	logger.Info(host, "    ", cmd)
	session, err := Connect(User, Passwd, PrivateKeyFile, host)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]Error create ssh session failed %s", err.Error())
		}
	}()
	if err != nil {
		panic(1)
	}
	defer session.Close()

	b, err := session.CombinedOutput(cmd)
	logger.Debug("command result is:", string(b))
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]Error exec command failed %s", err.Error())
		}
	}()
	if err != nil {
		panic(1)
	}
	return b
}

//Copy is
func Copy(host, localFilePath, remoteFilePath string) {
	sftpClient, err := SftpConnect(User, Passwd, PrivateKeyFile, host)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]scpCopy: %s", err.Error())
		}
	}()
	if err != nil {
		panic(1)
	}
	defer sftpClient.Close()
	srcFile, err := os.Open(localFilePath)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]scpCopy: %s", err.Error())
		}
	}()
	if err != nil {
		panic(1)
	}
	defer srcFile.Close()

	dstFile, err := sftpClient.Create(remoteFilePath)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]scpCopy: %s", err.Error())
		}
	}()
	if err != nil {
		panic(1)
	}
	defer dstFile.Close()
	buf := make([]byte, 100*oneMBByte) //100mb
	totalMB := 0
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		length, _ := dstFile.Write(buf[0:n])
		totalMB += length / oneMBByte
		logger.Alert("transfer total size is:", totalMB, "MB")
	}
}
func readFile(name string) string {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		logger.Error(err)
		return ""
	}

	return string(content)
}
func sshAuthMethod(passwd, pkFile string) ssh.AuthMethod {
	var am ssh.AuthMethod
	if pkFile != "" {
		pkData := readFile(pkFile)
		pk, _ := ssh.ParsePrivateKey([]byte(pkData))
		am = ssh.PublicKeys(pk)
	} else {
		am = ssh.Password(passwd)
	}
	return am
}

//Connect is
func Connect(user, passwd, pkFile, host string) (*ssh.Session, error) {
	auth := []ssh.AuthMethod{sshAuthMethod(passwd, pkFile)}
	config := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: time.Duration(5) * time.Minute,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := AddrReformat(host)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}

	return session, nil
}

//SftpConnect  is
func SftpConnect(user, passwd, pkFile, host string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, sshAuthMethod(passwd, pkFile))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = AddrReformat(host)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func SendPackage(host, url, fileName string) {
	//only http
	isHttp := strings.HasPrefix(url, "http")
	downloadCmd := ""
	if isHttp {
		downloadParam := ""
		if strings.HasPrefix(url, "https") {
			downloadParam = "--no-check-certificate"
		}
		downloadCmd = fmt.Sprintf(" wget %s -O %s", downloadParam, fileName)
	}
	remoteCmd := fmt.Sprintf("cd /root &&  %s %s ", downloadCmd, url)
	localFile := fmt.Sprintf("/root/%s", fileName)
	if isHttp {
		go WatchFileSize(host, localFile, GetFileSize(url))
		Cmd(host, remoteCmd)
	} else {
		Copy(host, url, localFile)
	}
}
