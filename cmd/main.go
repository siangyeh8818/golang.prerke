package main

import (
	"fmt"
	"time"

	utility "github.com/siangyeh8818/golang.prerke/internal/utility"
)

/*
type Configs struct {
	Cfgs []Config `nodes`
}
*/

//sshuser="root"
//sshpassword="promise"

func main() {

	//sshuser := "root"
	sshpassword := "solarvm12345"
	sshport := "22"
	sshAddress := "172.16.155.170"
	deployUser := "Jia-Siang"

	fmt.Println("---------------------------")
	fmt.Println("--need to install sshpass--")
	fmt.Println("---------------------------")
	utility.CheckSSHPASS()
	utility.CheckSSHKeyExisted()

	//sshpassword = "pentiumvm"

	//cmds := ""

	var cmds string

	/*
		if user, err := utility.GetRkeUser(); err != nil {
			fmt.Println(err)
		} else {
			deployUser = user
		}
	*/
	config := utility.GetAddress()
	deployUser = config.User
	sshpassword1 := config.Password
	//sshpassword := getPasswd()
	sshpassword = sshpassword1
	fmt.Println("password:", sshpassword, sshpassword1, len(sshpassword), len(sshpassword1))
	//fmt.Println("show first address:", config.Cfgs[0].Address)
	//fmt.Println("lens of address:", len(config.Cfgs))
	fmt.Println("----- check for 5 secs--------")
	fmt.Println("deploy user:", deployUser)
	fmt.Println("deploy password:", sshpassword)
	fmt.Println("deploy node:", config.Nodes)
	fmt.Println("rke url:", config.RkeUrl)
	fmt.Println("------------------------------")
	time.Sleep(5 * time.Second)

	fmt.Println("----- Wget Rke --------")
	cmds = " yum install -y wget"
	utility.RunCommand(cmds)
	cmds = "wget " + config.RkeUrl
	utility.RunCommand(cmds)

	for ia := 0; ia < len(config.Nodes); ia++ {
		sshAddress = config.Nodes[ia].Address
		installdocker := config.Nodes[ia].Dockerversion
		fmt.Println("installed of address:", sshAddress)
		if installdocker == "" {
			installdocker = "https://releases.rancher.com/install-docker/18.06.sh"
		}
		fmt.Println("installed of docker location:", installdocker)

		//cmds = "systemctl stop firewalld"
		cmds = "systemctl stop firewalld"
		utility.RemoteTaskPipesNoWait(sshAddress, sshport, cmds)
		fmt.Println("remove firewalld done")
		cmds = "systemctl disable firewalld"
		utility.RemoteTaskPipesNoWait(sshAddress, sshport, cmds)
		fmt.Println("disable firewalld done")

		//os.Exit(3)

		/*
			cmd := "curl https://releases.rancher.com/install-docker/17.03.sh | sh"
			ret := RemoteSSHRun(sshAddress, sshport, cmd)
			fmt.Println(ret)
			userTask(sshAddress, sshport)
		*/
		/* install docker */
		cmds = fmt.Sprintf("curl %s | sh", installdocker)
		//cmds = "curl https://releases.rancher.com/install-docker/17.03.sh | sh"
		utility.RemoteTaskPipes(sshAddress, sshport, cmds)

		/* create user */
		// for ubuntu
		//cmds = fmt.Sprintf("adduser pentium --gecos \"First Last,RoomNumber,WorkPhone,HomePhone\" --disabled-password;echo \"pentium:%s\" | sudo chpasswd;gpasswd -a pentium docker", sshpassword)
		//cmds = fmt.Sprintf("adduser pentium ;echo \"pentium:%s\" | sudo chpasswd;usermod -aG wheel pentium;gpasswd -a pentium docker", sshpassword)
		if deployUser != "root" {
			cmds = fmt.Sprintf("adduser %s ;echo \"%s:%s\" | sudo chpasswd;usermod -aG wheel %s;gpasswd -a %s docker", deployUser, deployUser, sshpassword, deployUser, deployUser)
			utility.RemoteTaskPipes(sshAddress, sshport, cmds)
		}

		/* generate ssh key */
		//for ubuntu
		//cmds = "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q"
		//for centos

		//generate sshkey in root space
		//cmds = "ssh-keygen -t rsa -C \"comment\" -P \"\" -f \"/root/.ssh/id_rsa\" -q"
		//remoteTaskPipes(sshAddress, sshport, cmds)

		//remoteExpect()

		/* ssh-keygen in user pentium */
		//for ubuntu
		//cmds = fmt.Sprintf("sudo -iu %s ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q", deployUser)
		//for centos

		fmt.Println("create remote pentum sshkey")
		if deployUser != "root" {
			deployPath := fmt.Sprintf("/home/%s", deployUser)

			//cmds = fmt.Sprintf("sudo -iu %s ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \"/home/%s/.ssh/id_rsa\" -q", deployUser, deployUser)
			cmds = fmt.Sprintf("sudo -iu %s  sh -c 'echo \"y\"|ssh-keygen -t rsa -C \"comment\" -P \"\" -f \"%s/.ssh/id_rsa\" -q '", deployUser, deployPath)
			fmt.Println("remote command:", cmds)
			utility.RemoteTaskPipes(sshAddress, sshport, cmds)
		}

		/* ssh-copy-id something */
		fmt.Println("into sshcopy")
		utility.SshCopy(deployUser, sshpassword)
		utility.SshCopyRoot(deployUser, sshpassword)
		fmt.Println("trying launch the following command for testing")
		fmt.Printf("sudo -u pentium ssh 'pentium@%s' docker ps \n", sshAddress)
		fmt.Printf("sudo -u pentium ssh 'pentium@%s' systemctl status firewalld\n", sshAddress)
	}

}
