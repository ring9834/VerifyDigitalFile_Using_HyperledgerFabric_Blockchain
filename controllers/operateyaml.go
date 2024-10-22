package controllers

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

////////////////YAML EXAMPLE 1//////////////////
type Stu struct {
	Name  string `yaml:"Name"`
	Age   string `yaml:"Age"`
	Sex   string `yaml:"Sex"`
	Class string `yaml:"class"`
}

func ReadYamlFile(src string) {
	var st Stu //定义一个结构体变量

	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Print(err)
	}
	//yaml文件内容影射到结构体中
	err1 := yaml.Unmarshal(config, &st)
	if err1 != nil {
		fmt.Println("error")
	}
}

func WriteYamlFile(src string) {
	stu := &Stu{
		Name:  "hzx",
		Age:   "12",
		Sex:   "male",
		Class: "six four",
	}
	data, err := yaml.Marshal(stu)
	if err != nil {
		fmt.Print(err)
	}
	err = ioutil.WriteFile(src, data, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

//////////////////////YAML EXAMPLE 2//////////////////////////

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func operateYaml_Example_2() {
	t := T{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))

}

///////////////////////////////////////////
