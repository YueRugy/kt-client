package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
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
	con := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold:  20,
		SleepWindow:            5,
	}
	hystrix.ConfigureCommand("my_command", con)
	//w := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		//go func() {
		//w.Add(1)
		//defer w.Done()
		prodChan := make(chan Prod, 1)
		errs := hystrix.Go("my_command", func() error {
			prod, _ := getProd()
			prodChan <- prod
			return nil
		}, func(err error) error {
			prod, _ := recProd()
			//prodChan <- prod
			e := errors.New("my time out")
			fmt.Print(prod)
			return e
			//return errors.New("my time out")
			//return nil
		})
		select {
		case p := <-prodChan:
			fmt.Println(p)
		case e := <-errs:
			fmt.Println(e)
		}
		time.Sleep(time.Second * 1)
		//}()
	}
	//w.Wait()
	//	err := hystrix.Do("my_command", func() error {

	/*if err != nil {
		fmt.Println(err.Error())
	}*/
	//time.Sleep(time.Second * 1)

}

type Prod struct {
	ID    int
	Title string
	Price int
}

func getProd() (Prod, error) {
	rn := rand.Intn(10)
	if rn < 5 {
		time.Sleep(time.Second * 3)
	}
	return Prod{
		ID:    101,
		Title: "GoLang 教程",
		Price: 100,
	}, nil
}

func recProd() (Prod, error) {
	return Prod{
		ID:    999,
		Title: "go-kit教程",
		Price: 114,
	}, nil
}
