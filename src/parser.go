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

func generateSteps(f string) steps{
  var c steps
  c.getStep(f)

  extends_step_path_dash := c.Entrance
  extends_step_path := strings.Replace(extends_step_path_dash, "__", "/", 3)
  // fmt.Println("extends_step_path -->", c, extends_step_path, extends_step_path_dash)
  if extends_step_path == "" {
    return c
  }
  p := extends_step_path + ".yaml"
  // fmt.Println("do me")
  // open the extends_path_file
  c1  := generateSteps(p)
  if c1.Entrance == extends_step_path_dash  {
    err := fmt.Sprintf("Recursion Entrance Error: You should never set entrance to yourself at [ %s ] [%s]", p, c1.Entrance)
    panic(err)
  }

  // fmt.Println(c)
  // fmt.Println("c2 steps -->, ", c1.Steps)
  all_steps := append(c1.Steps, c.Steps...)
  // fmt.Println("all ---> ", all_steps)
  c.Steps = all_steps
  // fmt.Println("c.Steps new -->", c.Steps)
  return c

}

func main() {
  f := "step.yaml"
  c := generateSteps(f)

  fmt.Println("c -->", c)
}
