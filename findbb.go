package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"flag"
	"os"
	"io"
	"path"
)

// Project Powered by CHAOS - https://chaos.projectdiscovery.io/
// Project Discovery - https://projectdiscovery.io/
// System00Sec - https://system00sec.org/
// Github - https://github.com/system00-security/findbb
// Twitter - https://twitter.com/0xjoyghosh
// copywrite 2023 System00 Security

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
	warn = "\033[33m[WARN]\033[0m No Internet Connection"
	found = "\033[32m[Found Program]\033[0m"
	err = "\033[31m[ERROR]\033[0m"
	info = "\033[34m[INFO]\033[0m"
	
)

var headers = map[string]string{
	"User-Agent": userAgent,
}
func endpoint1(domain string) string {
	req, err := http.Get("https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/master/chaos-bugbounty-list.json")
	if err != nil {
		return warn
	}
	defer req.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return warn
	}

	for _, i := range data["programs"].([]interface{}) {
		for _, j := range i.(map[string]interface{})["domains"].([]interface{}) {
			if j == domain {
				return i.(map[string]interface{})["url"].(string)
			}
		}
	}
	return ""
}

type Program struct {
	Name        string `json:"name"`
	ProgramURL  string `json:"program_url"`
	URL         string `json:"URL"`
	Bounty      bool   `json:"bounty"`
}

func subdomain(program_url string) string {
	resp, err := http.Get("https://chaos-data.projectdiscovery.io/index.json")
	if err != nil {
		fmt.Println(err, "Error retrieving data:", err)
		return ""
	}
	defer resp.Body.Close()
	var programs []Program
	if err := json.NewDecoder(resp.Body).Decode(&programs); err != nil {
		fmt.Println(err, "Error decoding data:", err)
		return ""
	}
	var targetProgram Program
	for _, program := range programs {
		if program.ProgramURL == program_url {
			targetProgram = program
			break
		}
	}
	if (Program{}) == targetProgram {
		fmt.Println(err,"Program not found\n")
		return ""
	}
	return targetProgram.URL
	
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}



func main() {
	logo:=`
    _____         _____  ___ 
   / __(_)__  ___/ / _ )/ _ )
  / _// / _ \/ _  / _  / _  |
 /_/ /_/_//_/\_,_/____/____/  
            system00sec.org 					  
		`
	fmt.Println(logo)

	var domain string
	var enumerate bool
	flag.StringVar(&domain, "d", "", "Domain to search")
	flag.BoolVar(&enumerate, "e", false, "Enable Subdomain Enumeration")
	flag.Parse()

	if domain == "" {
		fmt.Println("Usage: findbb -d example.com\n")
		return
	}
	program_url := endpoint1(domain)
	if program_url == "" {
		fmt.Println(err," Program not found\n")
		return
	}
	fmt.Println(found, program_url)
	if enumerate {
		if program_url != "" {
			zipurl := subdomain(program_url)
			path := path.Base(zipurl)
			fmt.Println(info, "Downloading", path)
			err := downloadFile(path, zipurl)
			if err != nil {
				fmt.Println(err, "Download failed")
			}
			fmt.Println(info, "Downloaded", path,"\n")
		}else {
			fmt.Println(err, " Subdomain enumeration failed \n")
		}
	}



}
