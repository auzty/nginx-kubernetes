package main

import (
	"fmt"
	"time"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
  "k8s.io/client-go/tools/cache"
  "k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
  "log"
  "io/ioutil"
  "strings"
  "strconv"
  "syscall"
  "os"
)

type Organization struct{
  Name string
  FullInternalName string
  Domain string 
  InternalDomain string 
  Port int32
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "services", v1.NamespaceDefault,
		fields.Everything())
	    _, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second * 0,
		cache.ResourceEventHandlerFuncs{
		    AddFunc: func(obj interface{}) {
	        log.Println("Watching Kubernetes Service")

                        mapping := obj.(*v1.Service)

			svc := mapping.Spec

      clustername := os.Getenv("CLUSTERNAME")
      if (clustername == "") {
        clustername = "cluster.local"
      }

      name := mapping.ObjectMeta.Name
      domain := mapping.ObjectMeta.Labels["domain"]
      internal := mapping.ObjectMeta.Labels["internaldomain"]
      port := svc.Ports[0].Port

      fullinternalname := mapping.ObjectMeta.Name + "." + mapping.ObjectMeta.Namespace + ".svc." + clustername

      details := Organization{Name:name,FullInternalName:fullinternalname,Domain:domain,InternalDomain:internal,Port:port}
      updateConf(details)


		    },
		},
	    )
	stop := make(chan struct{})
	    go controller.Run(stop)
	    for{
		time.Sleep(time.Second)
	    }	
}

func updateConf(org Organization){
 if org.InternalDomain != "" {
  fileloc:= "/etc/nginx/conf.d/"+org.InternalDomain+".conf"
  _, err := ioutil.ReadFile(fileloc)

  if err != nil {
// create new nginx configuration
    org.WriteNginx()
    reloadNginx()
    log.Println("Configuration "+org.InternalDomain+".conf Created..")
  } else {
    log.Println("File is exist...(please remove the existing file(s) to write new nginx conf")
  }
 }
}

func reloadNginx(){
  pidloc := "/var/run/nginx.pid"

  nginxpid, nginxerr := ioutil.ReadFile(pidloc)
  failOnError(nginxerr,"PID not found / nginx is not running")

  // convert []bytes become INT
  pid_str := strings.Replace(string(nginxpid),"\n","",-1)
  pid, err := strconv.Atoi(pid_str)

  fmt.Println("Reloading Nginx Configuration.....")
  time.Sleep(2 * time.Second)

  if nginxerr == nil {
    err =  syscall.Kill(pid,syscall.SIGHUP)
    failOnError(err,"Failed send SIGHUP to nginx")
  }

}
