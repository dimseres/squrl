package parser

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"groupsCrawl/config"
	"groupsCrawl/watcher/cmd"
	"groupsCrawl/watcher/services/models"
	"net/http"
	"sync"
	"time"
)

var (
	redis = config.Redis()
	ctx   = context.Background()
	wg    = sync.WaitGroup{}
)

type VkParser struct {
	Urls        *models.PriorityQueue
	Bus         chan models.ChanBus
	activeUrls  []string
	waitingUrls []string
	processing  [50]string
}

type RedisModel struct {
	Counter  int    `json:"counter"`
	SocketId string `json:"socketId"`
}

func (parser *VkParser) AddLink(link string) {
	hashedLink := md5.Sum([]byte(link))
	hash := hex.EncodeToString(hashedLink[:])

	data, err := redis.Connection.Get(ctx, "WATCHER:"+hash).Result()

	_payload := RedisModel{}
	if len(data) > 0 {
		err := json.Unmarshal([]byte(data), &_payload)
		if err != nil {
			fmt.Println(cmd.BlackBg + ">>> VK PARSER <<<" + cmd.Reset + "\t" + err.Error() + "")
		}
		_payload.Counter += 1
		_payload.SocketId = "ASD"

	} else {
		_payload.Counter += 1
		_payload.SocketId = "KOKO"
		parser.activeUrls = append(parser.activeUrls, hash)

	}

	parser.Urls.Push(&models.Item{
		Value: link,
	})

	fmt.Println(data)
	payload, err := json.Marshal(_payload)
	err = redis.Connection.Set(ctx, "WATCHER:"+hash, payload, 0).Err()
	if err != nil {
		fmt.Println(cmd.BlackBg + ">>> VK PARSER <<<" + cmd.Reset + "\t" + err.Error() + "")
	}
	fmt.Printf(cmd.BlackBg+">>> VK PARSER <<<"+cmd.Reset+"\taccept_link: %s \thash: %s\n", link, hex.EncodeToString(hashedLink[:]))
}

func (parser *VkParser) StartParser() {
	for {
		if parser.Urls != nil {
			parser.makeRequests()
			time.Sleep(time.Second * 1)
		}
	}
}

func (parser *VkParser) makeRequests() {
	links := func() []string {
		maxLen := 25
		out := make([]string, 0, maxLen)
		var selectLen int
		if parser.Urls.Len() > maxLen {
			selectLen = 25
		} else {
			selectLen = parser.Urls.Len()
		}
		for i := 0; i < selectLen; i++ {
			item := parser.Urls.Pop().(*models.Item)
			out = append(out, item.Value)
		}
		return out
	}()
	for _, item := range links {
		wg.Add(1)
		go func(url string) {
			var start = time.Now().UnixMilli()
			defer wg.Done()
			response, err := http.Get(url)
			if err != nil {
				fmt.Println(cmd.PurpleBg + ">>> VK PARSER <<<\t" + cmd.Reset + "error: " + err.Error())
				return
			}
			end := time.Now().UnixMilli() - start
			fmt.Printf(">>> [VK PARSER]: Request Done\t time: %d ms, link: %s\n", end, url)

			parser.Bus <- models.ChanBus{
				Service: "VK Parser",
				Message: "Request Done",
				Payload: response.Request.Host,
			}
			time.Sleep(time.Second * 2)
		}(item)
		parser.Urls.Push(&models.Item{
			Value: item,
		})
	}
	wg.Wait()
}
