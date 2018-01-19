package config

import (
	"bufio"
	"io"
	"os"
	"strings"
	"fmt"
)

const middle = "."

// TODO 写死的 不应该这么做
// TODO 分离到多个文件中 用文件名 + key 获取
const configPath = "./config/app.conf"

type configObj struct {
	Config map[string]string
	Model  string
}

var C configObj

func init() {
	C.InitConfig(configPath)
}

func (c *configObj) InitConfig(path string) {
	c.Config = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// 去掉两端的空白
		s := strings.TrimSpace(string(b))
		// 跳过注释
		if strings.Index(s, "#") == 0 {
			continue
		}

		// 解析配置模块标识
		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.Model = strings.TrimSpace(s[n1+1: n2])
			continue
		}

		if len(c.Model) == 0 {
			fmt.Println("no Model")
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		pos := strings.Index(value, "\t#")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, " #")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, "\t//")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, " //")
		if pos > -1 {
			value = value[0:pos]
		}

		if len(value) == 0 {
			value = ""
		}

		model := c.Model + middle + key
		c.Config[model] = strings.TrimSpace(value)
	}
}

func Read(key string) string {
	v, found := C.Config[key]
	if !found {
		return ""
	}
	return v
}
