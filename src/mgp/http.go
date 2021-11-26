package mgp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func HttpPostForm(transferUrl string, msgStruct ProduceData, uuid string) string{
	t1 := time.Now()

	args := url.Values{}
	args.Add("ct", strconv.Itoa(msgStruct.Ct))
	args.Add("uuid", msgStruct.Uuid)
	args.Add("biz_json", msgStruct.BizJson)

	resp, err := http.PostForm(transferUrl, args)
	if err != nil {
		fmt.Println("error")
		return ""
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error")
		return ""
	}

	elapsed := time.Since(t1)
	log.Printf("%s | ET: %s | HttpRes: %s", uuid, elapsed, bs)

	return string(bs)
}


