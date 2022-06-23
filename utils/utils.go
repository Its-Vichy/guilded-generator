package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"os"
)

var (
	proxy_index     int
	usernames_index int
	bio_index       int
	pfp_index       int
	status_index    int

	Proxies, _   = ReadLines("data/proxies.txt")
	usernames, _ = ReadLines("data/usernames.txt")
	bio, _       = ReadLines("data/bio.txt")
	pfp, _       = ReadLines("data/pfp.txt")
	status, _    = ReadLines("data/status.txt")
)

func RemoveIProxy(proxy string, lst []string) []string {
	new := []string{}

	for _, item := range lst {
		if item != proxy {
			new = append(new, item)
		}
	}

	return new
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func RandHexString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func AppendLine(path string, line string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		return
	}
}

func GetNexBio() string {
	bio_index++
	if bio_index == len(bio) || bio_index > len(bio) {
		bio_index = 0
	}
	return bio[bio_index]
}

func GetNexPfP() string {
	pfp_index++
	if pfp_index == len(pfp) || pfp_index > len(pfp) {
		pfp_index = 0
	}
	return pfp[usernames_index]
}

func GetNexStatus() string {
	status_index++
	if status_index == len(status) || status_index > len(status) {
		status_index = 0
	}
	return status[status_index]
}

func GetNexProxie() string {
	proxy_index++

	// Giga chad proxies reloading
	if len(Proxies) <= 5 {
		Proxies, _ = ReadLines("data/proxies.txt")
	}
	
	if proxy_index == len(Proxies) || proxy_index > len(Proxies) {
		proxy_index = 0
	}

	return Proxies[proxy_index]
}

func GetNexUsername() string {
	usernames_index++
	if usernames_index == len(usernames) || usernames_index > len(usernames) {
		usernames_index = 0
	}
	return usernames[usernames_index]
}
