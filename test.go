package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"math/rand"
	"net/http"
	"time"
)

func ml(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !l.Allow() {
			http.Error(writer, "too many request", http.StatusTooManyRequests)
			return
		}
		//log.Println("1111")
		next.ServeHTTP(writer, request)
	})
}

var (
	l = rate.NewLimiter(1, 5)
)

func main() {
	/*limit := rate.NewLimiter(1, 5)
	for {
		//_ = limit.Wait(context.Background())
		_ = limit.WaitN(context.Background(), 2)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second * 1)
	}*/

	/*mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK!!!"))
		//log.Println("2222")
	})
	_ = http.ListenAndServe(":8090", ml(mux))*/

	rand.Seed(time.Now().UnixNano())
	for {
		prod, _ := getProd()
		fmt.Println(prod)
		time.Sleep(time.Second * 1)
	}
}

type Prod struct {
	ID    int
	Title string
	Price int
}

func getProd() (Prod, error) {
	rn := rand.Intn(10)
	if rn < 6 {
		time.Sleep(time.Second * 3)
	}
	return Prod{
		ID:    101,
		Title: "GoLang 教程",
		Price: 100,
	}, nil
}
