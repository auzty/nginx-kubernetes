package main

import(
  "fmt"
  "io/ioutil"
  "text/template"
  "os"
)

func (org *Organization) WriteNginx() error {
  template_path := "/opt/nginx-k8sapi/conf.example"
  buf, err := ioutil.ReadFile(template_path)
  failOnError(err,"template file not found")

  text := string(buf)
  conf_file := fmt.Sprintf("%s.conf", org.InternalDomain)
  fd, err := os.Create(fmt.Sprintf("/etc/nginx/conf.d/%s",conf_file))
  failOnError(err,"File Create Error")

  defer func(){
    fd.Close()
  }()

  tmpl, err := template.New("conf").Parse(text)
  failOnError(err,"template parsing error")

  tmpl.Execute(fd,org)
  return nil
}
