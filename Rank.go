package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Rank struct {
	Rank  int    `json:"rank"`
	Nick  string `json:"nick"`
	Score int    `json:"score"`
	Date  string `json:"date"`
}

const (
	RANK_FILE = "top.txt"
	MAX_RANK  = 20
)

func ReadRankList(count int) ([]Rank, error) {

	rankList := make([]Rank, 0, count)

	file, err := os.Open(RANK_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for i := 0; i < count; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				Error.Println(err)
			}
			break
		}

		info := strings.Split(string(line), ",")
		if len(info) != 3 {
			Error.Printf("Invalid line %d: %s", i, string(line))
			break
		}

		s, err := strconv.Atoi(info[1])
		if err != nil {
			Error.Printf("Invalid score %d: %s", i, string(line))
			break
		}

		r := Rank{
			Rank:  i + 1,
			Nick:  info[0],
			Score: s,
			Date:  info[2],
		}

		rankList = append(rankList, r)
	}

	return rankList, nil
}

func SaveTopRank(nick string, score int) (int, error) {

	// 파일을 읽는다.
	file, err := os.OpenFile(RANK_FILE, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	buf := ""
	rank := -1

	reader := bufio.NewReader(file)
	for i := 0; i < MAX_RANK; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				if rank == -1 {
					rank = i + 1
					buf += fmt.Sprintf("%s,%d,%s\n", nick, score, time.Now().Format("2006-01-02 15:04:05"))
				}
			} else {
				Error.Println(err)
			}
			break
		}

		info := strings.Split(string(line), ",")
		if len(info) != 3 {
			Error.Printf("Invalid line %d: %s", i, string(line))
			break
		}

		s, err := strconv.Atoi(info[1])
		if err != nil {
			Error.Printf("Invalid score %d: %s", i, string(line))
			break
		}

		if rank == -1 && score > s {
			i++
			rank = i
			buf += fmt.Sprintf("%s,%d,%s\n", nick, score, time.Now().Format("2006-01-02 15:04:05"))
		}

		if i < 20 {
			buf += string(line) + "\n"
		}
	}

	if rank != -1 {
		// 파일을 다시 쓴다.
		file.Seek(0, 0)
		err := file.Truncate(0) // 기존 내용을 완전히 지움
		if err != nil {
			Error.Printf("Invalid truncate: %s", err)
			return rank, err
		}
		_, err = file.WriteString(buf)
		if err != nil {
			Error.Printf("Invalid write: %s", err)
			return rank, err
		}
	}

	return rank, nil
}
