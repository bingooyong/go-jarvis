package deploy

import (
	"fmt"
	"github.com/lvyong1985/go-jarvis/config"
	"github.com/lvyong1985/go-jarvis/funcs"
	"github.com/lvyong1985/go-jarvis/models"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
	"strings"
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

	var prefix string
	for idx, server := range d.ProjectDetail.ServerList {
		serverDetail := models.GetServerDetail(server)
		prefix = "1." + strconv.Itoa(idx)

		sshClient, err := funcs.NewSSH(serverDetail.Ip, serverDetail.Username, serverDetail.Password, serverDetail.PrivateKey, serverDetail.Port, d.logConsole)
		if err != nil {
			logrus.Infof("[%s]create ssh client error :%s", serverDetail.Ip, err)
			continue
		}
		defer sshClient.Close()
		logrus.Infof("开始发布 Project:%s, Idx:%s Server:%s ", d.ProjectDetail.Code, strconv.Itoa(idx), server)
		commands := []Command{
			&Mkdir{path: finaleDeploymentPath, name: "创建目录"},
			&Sync{sourcePath: releaseFile, targetPath: finaleDeploymentPath, name: "上传发布包"},
			&Link{sourcePath: finaleDeploymentPath, targetPath: linkPath, name: "创建软链接"},
			&Unzip{path: linkPath, file: releaseFile, name: "解压发布包"},
			&Script{workPath: linkPath, execScript: project.BeforeDeploymentScript, name: "发布前脚本"},
			&Script{workPath: linkPath, execScript: project.DeploymentScript, name: "发布脚本"},
			&Script{workPath: linkPath, execScript: project.AfterDeploymentScript, name: "发布后脚本"},
		}
		d.log("### " + prefix + " 同步文件到服务器(" + server + ")")
		for i, cmd := range commands {
			d.log(funcs.Sprintf("#### %s.%d 【%s】 \n", prefix, i, cmd.getName()), funcs.Sprintf("%v", cmd))
			if err := cmd.exec(sshClient); err != nil {
				d.logf("# 项目 【%s】【%s】", d.ProjectDetail.Code, "===========发版失败========")
				d.Done(models.FAIL)
				return
			}
		}
	}
	d.logf("# 项目 【%s】【%s】", d.ProjectDetail.Code, "===========发版成功========")
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

func (d *Deploy) log(logs ...string) {
	for _, log := range logs {
		d.msg <- []byte(log + "\n\n")
	}
}

func (d *Deploy) logf(f string, args ...interface{}) {
	d.msg <- []byte(fmt.Sprintf(f, args...) + "\n\n")
}

func (d *Deploy) Done(code models.DeploymentStatusCode) {
	defer d.logConsole.Close()
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
