package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	TOML "github.com/pelletier/go-toml"
)

type runitSpec struct {
	RunitCfg runitCfg       `toml:"runit"`
	Services []runitService `toml:"service"`
}

type runitCfg struct {
	Directory string `toml:"directory"`
}

type runitService struct {
	Name        string `toml:"name"`
	Directory   string `toml:"directory"`
	BootCmds    string `toml:"boot"`
	RunCmds     string `toml:"run"`
	Conditional string `toml:"conditional"`
}

func abspath(path string) string {
	if len(path) > 2 {
		if path[:2] == "~/" {
			usr, err := user.Current()
			if err != nil {
				panic(err.Error())
			}
			path = filepath.Join(usr.HomeDir, path[2:])
		}
	}
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return path
}

func writeScript(basePath, srvType, srvName, dir, cmds string) {
	srvPath := fmt.Sprintf("%s/%s/%s", basePath, srvType, srvName)
	err := os.MkdirAll(srvPath, 0744)
	if err != nil {
		log.Fatal(err)
	}
	fn := fmt.Sprintf("%s/run", srvPath)
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := "#!/bin/sh\n\n"
	s += fmt.Sprintf("cd %s\n", abspath(dir))
	s += fmt.Sprintf(". %s/envvars\n\n", basePath) // POSIX standard "source"
	if srvType == "service" {
		//
		// TODO: service can *only() be 1 line of commands, due to "exec" !
		//
		s += fmt.Sprintf("exec %s\n", cmds)
	} else {
		//
		// TODO: can be multiline
		//
		s += cmds
	}
	f.WriteString(s)
	err = os.Chmod(fn, 0744)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// usage: runit-conf config.toml

	// check command line
	if len(os.Args) != 2 {
		fmt.Println("Error: No TOML file specified")
		os.Exit(1)
	}

	// read input file
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// parse toml from input
	v := runitSpec{}
	err = TOML.Unmarshal(b, &v)
	if err != nil {
		log.Fatal(err)
	}

	// (re)create runit path
	basePath := abspath(v.RunitCfg.Directory)
	err = os.RemoveAll(basePath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(basePath, 0744)
	if err != nil {
		log.Fatal(err)
	}

	// write all scripts (both stage 1: boot, and stage 2: run)
	var bootAllList []string
	numServices := 0
	for _, srv := range v.Services {
		// skip if service depends on an envvar that is not set
		if len(srv.Conditional) > 0 {
			if os.Getenv(srv.Conditional) != "1" {
				continue
			}
		}

		// TODO: ensure srv.Name is unqiue and correct format type,
		//       else ignore
		// TODO: any srv.Directory can be:
		//        - absolute
		//        - user expanded
		//        - relative to the main running location

		// write shell scripts
		if len(srv.BootCmds) > 0 {
			writeScript(basePath, "boot", srv.Name, srv.Directory, srv.BootCmds)
			bootAllList = append(bootAllList, srv.Name)
		}
		if len(srv.RunCmds) > 0 {
			writeScript(basePath, "service", srv.Name, srv.Directory, srv.RunCmds)
			numServices++
		}
	}

	// write dummy program if no services
	if numServices < 1 {
		writeScript(basePath, "service", "dummy", "", "sleep 2147483647")
	}

	// write "boot_all" main script
	s := "#!/bin/sh\n\n"
	for _, srvName := range bootAllList {
		s += fmt.Sprintf("%s/boot/%s/run\n", basePath, srvName)
	}
	fn := fmt.Sprintf("%s/boot_all", basePath)
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(s)
	err = os.Chmod(fn, 0744)
	if err != nil {
		log.Fatal(err)
	}

	// return the absolute base path on success
	// (to be consumed by calling runit-bootstrap)
	fmt.Print(basePath)
}
