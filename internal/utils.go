package internal

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	. "fmt"
	"github.com/dlclark/regexp2"
	"github.com/go-redis/redis/v7"
	"net/http"
	"strings"
	"time"
)

func FindFirstOrDefaultMatchUseRegex2(content string, exp string) string {
	regexp2 := regexp2.MustCompile(exp, 0)
	match, err := regexp2.FindStringMatch(content)
	if err != nil || match == nil {
		return ""
	}

	groups := match.Groups()
	if len(groups) > 1 {
		return groups[1].Capture.String()
	}

	return match.String()
}

func FindAllMatchString(content string,exp string) []string {
	regexp2 := regexp2.MustCompile(exp, 0)
	match, err := regexp2.FindStringMatch(content)
	if err != nil || match == nil {
		return nil
	}

	result := make([]string,0)
	groups := match.Groups()
	for i := 0; i < len(groups); i++ {
		result = append(result, groups[i].Capture.String())
	}

	return result
}

func TrimSpace(content string) string {
	content = strings.Replace(content, "\n", "", -1)
	content = strings.TrimSpace(content)
	return content
}

func OK(response http.ResponseWriter,data interface{}){
	result := ApiResponse{
		Code: 200,
		Data: data,
		Date: time.Now(),
	}

	jsonValue, _ := json.Marshal(result)
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	Fprint(response, string(jsonValue))
}

func BadRequest(response http.ResponseWriter) {
	result := ApiResponse{
		Code: 400,
		Data: nil,
		Date: time.Now(),
	}

	jsonValue, _ := json.Marshal(result)
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	Fprint(response, string(jsonValue))
}

func RedisNewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{Addr:RedisServerAddr})
	_,err := client.Ping().Result()
	if err == nil {
		return client
	}

	return nil
}

func GetValueFromCache(key string,v interface{}) interface{} {
	client := RedisNewClient()
	if client == nil {
		return nil
	}

	defer client.Close()

	cacheValue,err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	if err := json.Unmarshal([]byte(cacheValue),v); err == nil {
		return v
	}

	return nil
}

func SetValueToCache(key string,v interface{}) {
	client := RedisNewClient()
	if client == nil {
		return
	}

	defer client.Close()
	client.SetNX(key,v,CacheExpirationTime).Result()
}

func Md5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}