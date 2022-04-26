package parser

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"groupsCrawl/cmd"
	"groupsCrawl/config"
	"groupsCrawl/watcher/models"
	"groupsCrawl/watcher/services/qrgen"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	redis = config.Redis()
	ctx   = context.Background()
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
	Url      string `json:"url"`
}

func (parser *VkParser) createRedisModel(hash string, link string, socketId string) RedisModel {
	data, _ := redis.Connection.Get(ctx, "WATCHER:"+hash).Result()
	model := RedisModel{}
	if len(data) > 0 {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			fmt.Println(cmd.BlackBg + ">>> VK PARSER <<<" + cmd.Reset + "\t" + err.Error() + "")
		}
		model.Counter += 1
		model.SocketId = socketId
	} else {
		model.Counter += 1
		model.SocketId = socketId
		model.Url = link
		parser.Urls.Push(&models.Item{
			Value:    link,
			Priority: 100,
		})
	}
	return model
}

func (parser *VkParser) saveDataRedis(model *RedisModel, hash string) {
	payload, err := json.Marshal(model)
	err = redis.Connection.Set(ctx, "WATCHER:"+hash, payload, 0).Err()
	if err != nil {
		fmt.Println(cmd.BlackBg + ">>> VK PARSER <<<" + cmd.Reset + "\t" + err.Error() + "")
	}
}

func (parser *VkParser) AddLink(link string) {
	hashedLink := md5.Sum([]byte(link))
	hash := hex.EncodeToString(hashedLink[:])

	urlData := parser.createRedisModel(hash, link, "ASD")
	parser.saveDataRedis(&urlData, hash)

	fmt.Printf(cmd.BlackBg+">>> VK PARSER <<<"+cmd.Reset+"\taccept_link: %s \thash: %s\n", link, hex.EncodeToString(hashedLink[:]))
	qrgen.QrFromUrl(link, hash)
}

func (parser *VkParser) StartParser() {
	parser.loadUrls()
	for {
		if parser.Urls != nil {
			parser.makeRequests()
			time.Sleep(time.Second * 1)
		}
	}
}

func (parser *VkParser) loadUrls() {
	keys, err := redis.Connection.Keys(ctx, "WATCHER:*").Result()
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, key := range keys {
		fmt.Println(key)
		keyData, err := redis.Connection.Get(ctx, key).Result()
		model := RedisModel{}
		err = json.Unmarshal([]byte(keyData), &model)
		if err != nil {
			log.Fatal(err.Error())
		}
		parser.Urls.Push(&models.Item{Value: model.Url})
	}
}

func (parser *VkParser) getQueueUrls() []string {
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
}

func (parser *VkParser) parseUrl(url string, wg *sync.WaitGroup) {
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

	parser.Urls.Push(&models.Item{
		Value: url,
	})
}

func (parser *VkParser) makeRequests() {
	links := parser.getQueueUrls()
	wg := sync.WaitGroup{}
	for _, item := range links {
		wg.Add(1)
		go parser.parseUrl(item, &wg)
	}
	wg.Wait()
}
