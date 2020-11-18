package main

import (
    "fmt"
    "strings"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "os"
)


type Step struct {
  Cmd string `yaml:"cmd"`
  Location string `yaml:"location"`
  Value string `yaml:"value,omitempty"`
  Name string `yaml:"name,omitempty"`
  Desc string `yaml:"desc,omitempty"`
}

type Steps struct {
  Version string `yaml:"version"`
  GroupName string `yaml:"groupName"`
  Entrance string `yaml:"entrance,omitempty"`
  Steps []Step `yaml:"steps"`
}

func (c *Steps) getStep(f string) *Steps {

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

func generateSteps(f string) Steps{
  var c Steps
  c.getStep(f)

  extends_step_path_dash := c.Entrance
  extends_step_path :=  strings.Replace(extends_step_path_dash, "__", "/", 3)
  if extends_step_path == "" {
    return c
  }
  p := "./input/" + extends_step_path + ".yaml"
  c1  := generateSteps(p)
  if c1.Entrance == extends_step_path_dash  {
    err := fmt.Sprintf("Recursion Entrance Error: You should never set entrance to yourself at [ %s ] [%s]", p, c1.Entrance)
    panic(err)
  }
  all_steps := append(c1.Steps, c.Steps...)
  c.Steps = all_steps
  return c

}

func extendsSteps(f string) {
  c := generateSteps(f)
  cs, err := yaml.Marshal(&c)
  if err != nil {
    panic("Marshal failure")
  }

  f2 :=  strings.Replace(f, "input", "outputs", 3)
  err = ioutil.WriteFile(f2, cs, 0644)
  if err != nil {
    panic("Marshal failure")
  }
}

func allSteps(d string) {
  // filter files by configure the case here
  // case about your product, which you want to test.
  files, err := ioutil.ReadDir(d)
  if err != nil {
      log.Fatal(err)
  }

  for _, f := range files {
          fi, err := os.Stat("./input/" +  f.Name())
          if err != nil {
              fmt.Println(err)
              continue // skip dir

          }
          switch mode := fi.Mode(); {
          case mode.IsDir():
              // do directory stuff
              fmt.Println("directory")
          case mode.IsRegular():
              // do file stuff
              fmt.Println("file")
              extendsSteps("./input/" + f.Name())
          }


  }
}



func main() {
  allSteps("./input")
}
