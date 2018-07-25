package deploy

import (
	"github.com/lvyong1985/go-jarvis/models"
	"path/filepath"
	"strconv"
	"github.com/sirupsen/logrus"
	"github.com/lvyong1985/go-jarvis/config"
	"github.com/lvyong1985/go-jarvis/funcs"
	"strings"
	"fmt"
)

var jarvis = "jarvis"

type Deploy struct {
	Deployment    *models.Deployment
	ProjectDetail *models.ProjectDetail
	logConsole    *funcs.LogConsole
	logPath       string
	msg           chan []byte
	quit          chan int
}

func NewDeploy(d *models.Deployment, p *models.ProjectDetail, logPath string) *Deploy {
	deploy := &Deploy{
		Deployment:    d,
		ProjectDetail: p,
		msg:           make(chan []byte, 1),
		quit:          make(chan int),
	}
	deploy.logPath = logPath
	deploy.logConsole = funcs.NewLogConsole(deploy.msg, logPath)
	return deploy
}

func (d *Deploy) Exec() {
	var releaseFile string
	var err error
	if d.Deployment.ReleasePath == "" {
		d.logf("# 项目 【%s】【%s】", d.ProjectDetail.Code, "===========打包步骤========")
		if releaseFile, err = d.buildFilePathFromSource(); err != nil {
			d.Done(models.FAIL)
			return
		}
		d.Deployment.ReleasePath = releaseFile
		d.Deployment.PreDeploy()
		d.logf("# 项目 【%s】【%s】", d.ProjectDetail.Code, "===========打包结束========")
	} else {
		releaseFile = d.Deployment.ReleasePath
	}

	d.logf("# 项目 【%s】【%s】", d.ProjectDetail.Code, "===========发版步骤========")
	d.log("## 1 文件同步到部署服务器")

	tag := d.Deployment.Tag
	finaleDeploymentPath := filepath.Join(config.Instance().Deploy.Path, jarvis, d.ProjectDetail.Code, tag)
	linkPath := d.Deployment.DeploymentPath
	project := d.ProjectDetail

	var cmd string
	var prefix string
	for idx, server := range d.ProjectDetail.ServerList {
		serverDetail := models.GetServerDetail(server)
		prefix = "1." + strconv.Itoa(idx)

		sshClient, err := funcs.NewSSH(serverDetail.Ip, serverDetail.Username, serverDetail.Password, serverDetail.PrivateKey, serverDetail.Port, d.logConsole)
		if err != nil {
			logrus.Infof("[%s]create ssh client error :%s", serverDetail.Ip, err)
			continue
		}
		logrus.Infof("开始发布 Project:%s, Idx:%s Server:%s ", d.ProjectDetail.Code, strconv.Itoa(idx), server)

		d.log("### " + prefix + "同步文件到服务器(" + server + ")")
		// mkdir
		d.log("#### " + prefix + ".1 创建文件存储目录")
		d.log("在部署服务器(" + server + ")上创建目录：" + finaleDeploymentPath)
		d.log("```shell\n" + "$ mkdir " + finaleDeploymentPath + "\n" + "```")
		sshClient.ExecCmd("mkdir -p " + finaleDeploymentPath)

		// sync file
		d.log("#### " + prefix + ".2 同步文件")
		d.log(
			"把文件(" + releaseFile + ")同步到部署服务器(" + server + ")上的" + finaleDeploymentPath + "目录")
		d.log("```shell\n" + "sftp> put " + releaseFile + " " + finaleDeploymentPath + "\n" + "```")
		sshClient.Put(releaseFile, finaleDeploymentPath)

		// link
		d.log("#### " + prefix + ".3 创建软链")
		d.log("在部署服务器(" + server + ")上创建软链")
		d.log("```shell\n" + "$ ln -vnsf " + finaleDeploymentPath + " " + linkPath + "\n" + "```")
		common := "ln -vnsf " + finaleDeploymentPath + " " + linkPath
		sshClient.ExecCmd(common)

		// 解压
		prefix = "2." + strconv.Itoa(idx)
		d.log("## 2 在部署服务器上执行命令")
		d.log("### " + prefix + " 在部署服务器(" + server + ")上执行命令")
		cmd = "cd " + linkPath + ";" + decompressFileCommand(releaseFile)
		d.log("```shell\n $ " + cmd + " \n```")
		sshClient.ExecCmd(cmd)

		exportCommand := "export TERM=xterm"
		sourceCommand := "source /etc/profile"
		bashProfile := "source ~/.bash_profile"
		bashRc := "source ~/.bashrc"
		currentFolder := "cd " + linkPath

		// before deployment command

		if project.BeforeDeploymentScript != "" {
			d.log("#### " + prefix + ".1 在部署服务器(" + server + ")上执行执行发布前脚本命令")
			d.log("```shell\n $" + project.BeforeDeploymentScript + "\n" + "```")
			d.log("命令执行结果")
			sshClient.ExecMulti(sourceCommand, exportCommand, bashProfile, bashRc, currentFolder, project.BeforeDeploymentScript, "sleep 1")
		}
		// deployment command
		if project.DeploymentScript != "" {

			d.log("#### " + prefix + ".2 在部署服务器(" + server + ")上执行执行发布脚本命令")
			d.log("```shell\n $" + project.DeploymentScript + "\n" + "```")
			sshClient.ExecMulti(sourceCommand, exportCommand, bashProfile, bashRc, currentFolder, project.DeploymentScript, "sleep 1")
		}

		if project.AfterDeploymentScript != "" {
			// after deployment command
			d.log("#### " + prefix + ".3 在部署服务器(" + server + ")上执行执行发布后脚本命令")
			d.log("```shell\n $" + project.AfterDeploymentScript + "\n" + "```")
			sshClient.ExecMulti(sourceCommand, exportCommand, bashProfile, bashRc, currentFolder, project.AfterDeploymentScript, "sleep 1")
		}
	}
	d.logf("# 项目 【%s】【%s】", d.ProjectDetail.Code, "===========发版结束========")
	d.Done(models.SUCCESS)
}

func (d *Deploy) buildFilePathFromSource() (releasFilePath string, err error) {
	detail := d.ProjectDetail
	sourcePath := detail.SourcePath
	wd := funcs.GetWorkPath()
	cmd := funcs.NewExecCmd(wd, d.logConsole)
	d.log("## 1 从代码服务器下载代码")
	if err = cmd.Exec("git clone --depth 1 " + sourcePath); err != nil {
		logrus.Errorf("[%s] git clone %s error ", detail.Code, sourcePath)
		return
	}
	d.log("## 2 执行打包脚本")
	if len(strings.TrimSpace(detail.PackageScript)) == 0 {
		if err = cmd.Exec("tar czvf ROOT.tar.gz * && mkdir output && mv ROOT.tar.gz ./output"); err != nil {
			logrus.Errorf("[%s] exec package_script error", detail.Code)
			return
		}
		return wd + "/output/ROOT.zip", nil
	}
	packageScript := detail.PackageScript
	if err = cmd.Exec(strings.Join(strings.Split(packageScript, "\n"), ";")); err != nil {
		logrus.Errorf("[%s] exec package_script error", detail.Code)
		return
	}
	return wd + "/output/ROOT.zip", nil
}

func (d *Deploy) log(log string) {
	d.msg <- []byte(log + "\n\n")
}

func (d *Deploy) logf(f string, args ...interface{}) {
	d.msg <- []byte(fmt.Sprintf(f, args...) + "\n\n")
}

func decompressFileCommand(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".gz":
		return "tar xzvf " + filepath.Base(filename)
	case ".tar":
		return "tar xvf " + filepath.Base(filename)
	case ".zip":
		return "unzip -o " + filepath.Base(filename)
	default:
		return ""
	}
}

func (d *Deploy) Done(code models.DeploymentStatusCode) {
	switch code {
	case models.SUCCESS:
		d.Deployment.Success()
		return
	case models.FAIL:
		d.Deployment.Fail()
		return
	default:
	}
}
