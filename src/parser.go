package main

import (
    "fmt"
    "strings"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)


type step struct {
  Cmd string `yaml:"cmd"`
  Location string `yaml:"location"`
  Value string `yaml:"value,omitempty"`
  Name string `yaml:"name,omitempty"`
}

type steps struct {
  Version string `yaml:"version"`
  GroupName string `yaml:"groupName"`
  Entrance string `yaml:"entrance"`
  Steps []step `yaml:"steps"`
}

func (c *steps) getStep(f string) *steps {

    yamlFile, err := ioutil.ReadFile(f)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
    return c
}

func generateSteps(){
  var c steps
  f := "step.yaml"
  c.getStep(f)

  extends_step_path_dash := c.Entrance
  extends_step_path := strings.Replace(extends_step_path_dash, "__", "/", 3)

  fmt.Println(c)
  fmt.Println(extends_step_path)
  p := extends_step_path + ".yaml"
  // open the extends_path_file
  c.getStep(p)
  if c.Entrance == extends_step_path_dash {
    err := fmt.Sprintf("Recursion Entrance Error: You should never set entrance to yourself at [ %s ]", p)
    panic(err)
  }
  fmt.Println(c)

}

func main() {
  generateSteps()
}
