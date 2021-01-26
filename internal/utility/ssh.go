package utility

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	gexpect "github.com/ThomasRooney/gexpect"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

var sshuser string
var sshpassword string
var sshport string
var sshAddress string
var deployUser string
var promptRE = regexp.MustCompile(`password`)

func RemoteSSHRun(addr string, port string, cmd string) string {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	//client, err := ssh.Dial("tcp", "172.16.155.137:22", &ssh.ClientConfig{
	fmt.Println(addrPort)
	//client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		Timeout:         time.Second * 10,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth: []ssh.AuthMethod{ssh.Password("^Two^Ten=1024$")},
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	//ce(err, "new session")
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
	if err := session.Run(cmd); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
	return b.String()

}
func CheckSSHPASS() {
	cmds := "yum list installed|grep sshpass"
	ret := RunCommand(cmds)
	fmt.Println("------ show check sshpass result-------")
	fmt.Println(ret)
	if ret == "" {
		fmt.Println("---- need to install sshpass ----- yum install sshpass")
		os.Exit(3)
	} else {
		fmt.Println("show", ret)
	}
	return
}
func SshCopyRoot(deployUser string, sshpassword string) {
	//command := ""
	// cannot su user here
	//cmd := fmt.Sprintf("/usr/bin/ssh-copy-id %s", sshAddress)
	//cmd := fmt.Sprintf("sudo -iu %s ssh-copy-id %s@%s", deployUser, deployUser, sshAddress)
	fmt.Println("allowing root to remote user with free ssh accessing")
	fmt.Println("but you need to install sshpass")
	cmds := fmt.Sprintf("sshpass -p %s ssh-copy-id -o StrictHostKeyChecking=no %s@%s", sshpassword, deployUser, sshAddress)
	fmt.Println(cmds)
	RunCommand(cmds)
	// sshpass -f password.txt ssh-copy-id -o StrictHostKeyChecking=no  pentium@172.16.155.101
	/*
		child, err := gexpect.Spawn(cmd)
		if err != nil {
			panic(err)
		}
		fmt.Println("into expect")
		child.Expect("password")
		fmt.Println("end expect")
		child.SendLine(sshpassword)
		child.Interact()
		child.Close()
	*/
}

func SshCopy(deployUser string, sshpassword string) {
	//command := ""
	// cannot su user here
	//cmd := fmt.Sprintf("/usr/bin/ssh-copy-id %s", sshAddress)
	//cmd := fmt.Sprintf("sudo -iu %s ssh-copy-id %s@%s", deployUser, deployUser, sshAddress)
	cmd := fmt.Sprintf("sudo -iu %s ssh-copy-id -o StrictHostKeyChecking=no %s@%s", deployUser, deployUser, sshAddress)
	child, err := gexpect.Spawn(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println("into expect")
	child.Expect("password")
	fmt.Println("end expect")
	child.SendLine(sshpassword)
	child.Interact()
	child.Close()
}
func RemoteTaskPipesNoWait(addr string, port string, cmds string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	fmt.Println(addrPort)
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}
	err = session.Shell()
	cmds = fmt.Sprintf("%s;%s", cmds, "exit")
	fmt.Println(cmds)
	//cmds := "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q; exit"
	cmdlist := strings.Split(cmds, ";")
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
		fmt.Println(c)

	}
	time.Sleep(10 * time.Second)

	return
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
}

func RemoteTaskPipes(addr string, port string, cmds string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	fmt.Println(addrPort)
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}
	err = session.Shell()
	cmds = fmt.Sprintf("%s;%s", cmds, "exit")
	fmt.Println(cmds)
	//cmds := "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q; exit"
	cmdlist := strings.Split(cmds, ";")
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
		fmt.Println(c)

	}

	err = session.Wait()
	fmt.Println("session out")
	//if err != nil {
	//		log.Fatal(err)
	//	}
	fmt.Println((outbt.String() + errbt.String()))
	return
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
}

func CreateSshKey(addr string, port string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	fmt.Println(addrPort)
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}
	err = session.Shell()
	cmds := "ssh-keygen -t rsa -C \"comment\" -P \"\" -f \".ssh/id_rsa\" -q; exit"
	cmdlist := strings.Split(cmds, ";")
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
		fmt.Println(c)

	}

	session.Wait()
	fmt.Println((outbt.String() + errbt.String()))
	return
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
}
func GetPasswd() string {

	//filename := os.Args[1]
	filename := "password.txt"
	passwd, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("password:", string(passwd))
	s := strings.Trim(string(passwd), "\n")
	s = strings.TrimSpace(s)
	return s
}

func GetRkeUser() (string, error) {

	var config RkeConfig
	filename := "cluster.yml"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("raw file:", string(source))
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println("shown in rke user reading", config.Nodes[0].User)
	return config.Nodes[0].User, nil

}

func GetAddress() Configs {

	var config Configs

	//filename := os.Args[1]
	filename := "precluster.yml"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//  source := []byte(data)
	fmt.Println("raw file:", string(source))
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("--- config:\n%v\n\n", config)
	fmt.Println("shown in reading", config.Nodes)
	fmt.Println("shown in reading", config.Password)
	//fmt.Println("len of cfg", len(config.Cfgs))
	//fmt.Println("len of value", len(config.Cfgs[0].Info))
	//fmt.Println("first info value", config.Cfgs[0].Info[0])
	//fmt.Println(config.Cfgs[0].Address)
	return config

}

/*
func RunCommand(commandStr string) string {
	cmdstr := commandStr
	out, _ := exec.Command("sh", "-c", cmdstr).Output()
	strout := string(out)
	return strout
}
*/
func CheckSSHKeyExisted() {
	if _, err := os.Stat("/root/.ssh"); os.IsNotExist(err) {
		fmt.Println("---------------------------------------------------------")
		fmt.Println("need sshkey existed for building password free for rke up")
		fmt.Println("---------------------------------------------------------")
		os.Exit(3)
	}
}
