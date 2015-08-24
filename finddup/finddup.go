package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}	
	defer file.Close()

	var prefix string
	if len(flag.Args()) > 1 {
		prefix = flag.Args()[1]
		log.Println(prefix)
	} else {
		prefix = ""
	}
	var out *os.File = nil
	if len(flag.Args()) > 2 {
		out, err = os.Create(flag.Args()[2])
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}
	var index map[string]string
	index = make(map[string]string)
	var duplicated map[string][]string
	duplicated = make(map[string][]string)
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dir, file := filepath.Split(scanner.Text())
		if _, ok := index[file]; ok {
			//log.Println(file, " is already present")
			if _, ok := duplicated[file]; !ok {
				duplicated[file] = make([]string, 0)
				duplicated[file] = append(duplicated[file], filepath.Join(prefix, index[file], file))
			}
			duplicated[file] = append(duplicated[file], filepath.Join(prefix, dir, file))
		} else {
			index[file] = dir
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for k := range duplicated {
		sha := make(map[string]string)
		for _, v := range duplicated[k] {
			hasher := sha256.New()
			ff, err := os.Open(v)
			if err != nil {
				log.Fatal(err)
			}
			defer ff.Close()
			if _, err := io.Copy(hasher, ff); err != nil {
				log.Fatal(err)
			}
			val := hex.EncodeToString(hasher.Sum(nil))
			if _, ok := sha[val]; !ok {
				sha[val] = v
			} else {
				if out == nil {
					log.Println("Duplicated sha:", sha[val], " ", v)
				} else {
					out.WriteString(fmt.Sprintf("\"%s\" \"%s\"\n", sha[val], v))
				}
			}
		}
	}
}

