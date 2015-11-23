// Copyright 2015 reposkeeper@reposkeeper.net. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

/*
Package shell is a shell engine.
可以对shell脚本进行一系列的执行，其实，此引擎并不只是对shell，它只是提供了一个shell环境，可以执行
Agent上面任何脚本，如 perl python ruby 等等

Author: reposkeeper
Time:   2015-11-22
Version: 0.1
Last update:

主要具备如下的特点：
	1. 可对脚本使用特定的用户执行
	2. 可对脚本传任何参数
	3. 可实时或异步查询执行返回
*/
package shell

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	parallelMax  = 10
	tmpDirectory = "./shell_tmp"
)

type shellRuntime struct {
	fileName    string            // shell的文件名
	tmpFileName string            // 临时的的shell名
	role        string            // 执行的用户
	param       map[string]string // shell的参数
	output      string            // 执行输出存放的文件名
	startTime   time.Time         // shell执行的开始时间
	timeout     int32             // shell 执行的超时时间，单位是秒
	isStop      bool              // 是否执行完毕
	ret         bool              // 执行结果成功还是失败
	feedBack    bool              // 执行结束之后是否取结果
}


type shellEngineManager struct {
	parallelCount int                        // 同时执行的shell数
	waitQueue     list.List                  // 等待执行的队列
	excutingQueue *[parallelMax]shellRuntime // 正在执行的队列

}

func (sr *shellRuntime) prepareEnv() (e error) {

	if sr.fileName == "" {
		return errors.New("shell文件名为空")
	}

	// 创建临时shell文件
	fp, err := ioutil.TempFile(tmpDirectory, sr.fileName)
	if err != nil {
		return errors.New("临时文件创建失败: " + err.Error())
	}
	sr.tmpFileName = fp.Name()

	// 创建临时shell的输出文件
	ofp, e := ioutil.TempFile(tmpDirectory, sr.fileName+"_output")
	if err != nil {
		return errors.New("临时输出文件创建失败: " + err.Error())
	}
	sr.output = ofp.Name()
	ofp.Close()

	// 设置bash 设置输出文件
	fp.WriteString("#!/bin/bash  \nexec 2>1 1>" + sr.output + "\n")

	// 设置传入的参数
	params := ""
	for k, v := range sr.param {
		params += fmt.Sprintln("export %s=%s", k, v)
	}
	fp.WriteString(params)
	
	
	return nil

}
