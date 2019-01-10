package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonParse interface{
	Load(filename string) error
}

//定义配置文件解析后的结构
type Config struct {
	Static  string
	Template string
}

type Employee struct {
	Name string
	Icon string
}

type EmployeeList []Employee


func (jst *Config) Load(filename string) error{
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, jst)
	if err != nil {
		return err
	}
	return nil
}

func (jst *EmployeeList) Load(filename string) error{
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, jst)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}