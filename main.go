package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	docID := flag.String("id", "1ZXB6pnQEvym9pkDDoSxl49mxiiMAkzx0tS50OuwRt3Q", "Google Drive / Sheet file ID")
	flag.Parse()

	url := fmt.Sprintf("https://drive.google.com/file/d/%s/view", *docID)
	fmt.Printf("[*] Target Document: %s\n", url)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[-] Failed to fetch document: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	// Document Title
	titleRe := regexp.MustCompile(`<title>(.*?) - Google Drive</title>`)
	if m := titleRe.FindStringSubmatch(content); len(m) > 1 {
		fmt.Printf("[+] Title: %s\n", m[1])
	}

	// Email Extraction
	emailRe := regexp.MustCompile(`[a-zA-Z0-9_.+-]+@gmail\.com`)
	emails := emailRe.FindAllString(content, -1)
	if len(emails) == 0 {
		fmt.Println("[-] No owner email found in viewer metadata.")
		return
	}

	email := emails[0]
	username := strings.Split(email, "@")[0]
	fmt.Printf("[+] Recovered Owner Email: %s\n", email)
	fmt.Printf("[+] Extracted Username:   %s\n", username)

	// Look up linked Google identity via public YouTube endpoint
	ytURL := fmt.Sprintf("https://www.youtube.com/@%s", username)
	fmt.Printf("[*] Checking linked profile: %s\n", ytURL)
	ytReq, _ := http.NewRequest("GET", ytURL, nil)
	ytReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")

	ytResp, err := client.Do(ytReq)
	if err == nil {
		defer ytResp.Body.Close()
		ytBody, _ := io.ReadAll(ytResp.Body)
		ytContent := string(ytBody)

		nameRe := regexp.MustCompile(`"mainEntity":\{"type":"Person","name":"([^"]+)"`)
		if m := nameRe.FindStringSubmatch(ytContent); len(m) > 1 {
			fmt.Printf("[!] Discovered Linked Real Name: %s\n", m[1])
		}

		channelRe := regexp.MustCompile(`"channelId":"([^"]+)"`)
		if m := channelRe.FindStringSubmatch(ytContent); len(m) > 1 {
			fmt.Printf("[!] Linked Channel / Gaia ID:   %s\n", m[1])
		}
	}
}
