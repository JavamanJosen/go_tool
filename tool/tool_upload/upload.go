package tool_upload

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"os/exec"
)

var (
	log = miner.Log()
)

func UploadToHadoopPath(filePathName, hadoopPath string) (string, error){
	//上传文件到Hadoop集群路径
	//uploadCommand := `cd %s && /app/hadoop/hadoop/bin/hdfs dfs -put %s hdfs://nswx%s`
	uploadCommand := `/app/hadoop/hadoop/bin/hdfs dfs -put %s hdfs://nswx%s`
	//把文件中的日期给补充上
	uploadCommand = fmt.Sprintf(uploadCommand, filePathName, hadoopPath)

	cmd := exec.Command("/bin/bash", "-c", uploadCommand)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("Execute Shell:%s failed with error:%s", uploadCommand, err.Error())
	}
	return string(output), err
}
