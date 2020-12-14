package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kitHttp "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"io"
	service "kt-client/services"
	"net/url"
	"os"
	"time"
)

func main1() {
	target, _ := url.Parse("http://localhost:8080")
	client := kitHttp.NewClient("GET", target, service.GetUserInfoRequest, service.GetUserInfoResponse)
	getUserInfo := client.Endpoint()
	ctx := context.Background()
	response, _ := getUserInfo(ctx, service.UserRequest{
		Uid: 101,
	})
	us := response.(service.UserResponse)
	fmt.Println(us.Name)
}

func main() {
	{
		//创建一个consulClient 指定地址和端口
		config := consulapi.DefaultConfig()
		config.Address = "localhost:8500"
		consul_client, _ := consulapi.NewClient(config)
		// 创建一个client通过consul
		client := consul.NewClient(consul_client)
		//go-kit 的logger
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
		}
		//创建一个实例
		{
			tags := []string{"primary"}
			instancer := consul.NewInstancer(client, logger, "userservice", tags, true)
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) {
					target, _ := url.Parse("http://" + service_url)
					client := kitHttp.NewClient("GET", target, service.GetUserInfoRequest, service.GetUserInfoResponse)
					return client.Endpoint(), nil, nil
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				//endpoints, _ := endpointer.Endpoints()
				//负载均衡器
				//mylb := lb.NewRoundRobin(endpointer)
				mylb := lb.NewRandom(endpointer, time.Now().UnixNano())
				//轮循执行
				for {
					e, _ := mylb.Endpoint()
					ctx := context.Background()
					response, _ := e(ctx, service.UserRequest{
						Uid: 101,
					})
					//us := response.(service.UserResponse)
					//fmt.Println(us.Name)
					fmt.Println(response)
					time.Sleep(time.Millisecond * 60)
				}

			}
		}
	}

}
