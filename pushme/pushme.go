package main

import (
	"flag"
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"gopkg.in/gcfg.v1"
	"log"
	"os/user"
	"path"
)

type Config struct {
	Pushme struct {
		Token string
	}
}

func main() {
	var cfg Config
	usr, _ := user.Current()
	dir := usr.HomeDir
	cfg_file := path.Join(dir, ".config", "golang-utils", "user.cfg")
	log.Println("Read configuration file: ", cfg_file)
	err := gcfg.ReadFileInto(&cfg, cfg_file)
	if err != nil {
		log.Fatalf("Failed to parse configuration file: %s", err)
	}
	token := cfg.Pushme.Token
	log.Println("Token: ", token)
	flag.Parse()
	// create a client for pushbullet
	pb := pushbullet.New(token)

	// create a 'note' push
	n := requests.NewNote()
	n.Title = flag.Args()[0]
	n.Body = flag.Args()[1]

	// send the note
	if _, err := pb.PostPushesNote(n); err != nil {
		log.Fatalln("Error: %s\n", err)
	}
}
