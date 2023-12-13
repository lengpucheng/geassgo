/*
----------------------------------------
@Create 2023/11/17-14:35
@Author lpc<lengpucheng@qq.com>
@Program Geassgo
@Describe main
----------------------------------------
@Version 1.0 2023/11/17
@Memo create this file
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lengpucheng/Geassgo/pkg/coderender"
	"github.com/lengpucheng/Geassgo/pkg/profess/helper"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"time"
)

var valuePath string
var targetPath string
var outputPath string
var mod string

func init() {
	flag.StringVar(&mod, "m", "task", "mod  default task\n1. package : package chart use -t dir -o chart.tar.gz\n"+
		"2. task : execute task use -t task.yaml -v values.yaml\n"+
		"3. chart : execute chart use -t chart.tar.gz -v values.yaml")
	flag.StringVar(&targetPath, "t", "", "target path")
	flag.StringVar(&valuePath, "v", "", "values file (可选)")
	flag.StringVar(&outputPath, "o", ".", "output path default .")
}

func main() {
	flag.Parse()
	if targetPath == "" {
		flag.Usage()
		return
	}
	var err error
	now := time.Now()
	switch mod {
	case "package":
		err = coderender.Archive(targetPath, outputPath)
	case "chart":
		var defValue = map[string]any{}
		if valuePath != "" {
			file, err := os.ReadFile(valuePath)
			if err != nil {
				slog.Error("load values fail!", "valuesPath", valuePath)
			} else {
				if err := yaml.Unmarshal(file, &defValue); err != nil {
					slog.Error("load values fail!", "valuesPath", valuePath)
				}
			}
		}
		_, err = helper.RunChart(context.Background(), targetPath, defValue)
	case "task":
		_, err = helper.RunClaim(context.Background(), targetPath, valuePath)
	default:
		slog.Error("the mod is not support", "mode", mod)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("用时： ", time.Now().UnixMilli()-now.UnixMilli(), "ms")
}
