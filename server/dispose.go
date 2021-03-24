/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:42
 */

package server

import (
	"fmt"

	"github.com/nichtsen/go-stress-testing/server/statistics"

	"github.com/nichtsen/go-stress-testing/server/client"
	"github.com/nichtsen/go-stress-testing/server/golink"

	"go-stress-testing/utils"

	"github.com/nichtsen/go-stress-testing/model"

	"sync"
	"time"

	"go-stress-testing/server/verify"

	"golang.org/x/net/context"
)

const (
	connectionMode = 1 // 1:顺序建立长链接 2:并发建立长链接
)

// 注册验证器
func init() {

	// http
	model.RegisterVerifyHttp("statusCode", verify.HttpStatusCode)
	model.RegisterVerifyHttp("json", verify.HttpJson)

	// webSocket
	model.RegisterVerifyWebSocket("json", verify.WebSocketJson)
}

// 处理函数
func Dispose(concurrency, totalNumber, extraJsonLength uint64, request *model.Request) {

	// 设置接收数据缓存
	ch := make(chan *model.RequestResults, 1000)
	ctx, cancel := context.WithCancel(context.Background())
	var (
		wg          sync.WaitGroup // 发送数据完成
		wgReceiving sync.WaitGroup // 数据处理完成
		//isStop      chan struct{}
	)

	wgReceiving.Add(1)
	go statistics.ReceivingResults(concurrency, ch, &wgReceiving)

	for i := uint64(0); i < concurrency; i++ {
		wg.Add(1)
		switch request.Form {
		case model.FormTypeHttp:

			go golink.Http(i, ch, totalNumber, &wg, request)

		case model.FormTypeWebSocket:
			// isStop = make(chan struct{}, 0)
			var str string
			if extraJsonLength > 0 {
				str = utils.RandString(extraJsonLength)
			}

			switch connectionMode {
			case 1:
				// 连接以后再启动协程
				ws := client.NewWebSocket(request.Url)
				err := ws.GetConn()
				if err != nil {
					fmt.Println("连接失败:", i, err)

					continue
				}
				worker := golink.NewWorker(ch, i, totalNumber, str, &wg, request, *ws)

				go worker.Run(ctx)
			case 2:
				// 并发建立长链接
				go func(i uint64) {
					// 连接以后再启动协程
					ws := client.NewWebSocket(request.Url)
					err := ws.GetConn()
					if err != nil {
						fmt.Println("连接失败:", i, err)

						return
					}

					worker := golink.NewWorker(ch, i, totalNumber, str, &wg, request, *ws)

					go worker.Run(ctx)
				}(i)

				// 注意:时间间隔太短会出现连接失败的报错 默认连接时长:20毫秒(公网连接)
				time.Sleep(5 * time.Millisecond)
			default:

				data := fmt.Sprintf("不支持的类型:%d", connectionMode)
				panic(data)
			}

		case model.FormTypeGRPC:
			// 连接以后再启动协程
			ws := client.NewGrpcSocket(request.Url)
			err := ws.Link()
			if err != nil {
				fmt.Println("连接失败:", i, err)
				continue
			}
			go golink.Grpc(i, ch, totalNumber, &wg, request, ws)
		default:
			// 类型不支持
			wg.Done()
		}
	}

	// 等待所有的数据都发送完成
	wg.Wait()

	// 延时1毫秒 确保数据都处理完成了
	time.Sleep(1 * time.Millisecond)
	close(ch)
	cancel()
	// wait a little while for all goroutines to quit
	time.Sleep(20 * time.Second)

	// 数据全部处理完成了
	wgReceiving.Wait()

	return
}
