/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nichtsen/go-stress-testing/heper"
	"github.com/nichtsen/go-stress-testing/model"
	"github.com/nichtsen/go-stress-testing/server/client"
)

const (
	firstTime    = 1 * time.Second // 连接以后首次请求数据的时间
	intervalTime = 5 * time.Second // 发送数据的时间间隔
)

var (
	// 请求完成以后是否保持连接
	keepAlive bool
	// extra string for larger bandwith
	extraStr string
)

func init() {
	keepAlive = true
}

type WSWorker struct {
	client.WebSocket
	chanId          uint64
	ch              chan<- *model.RequestResults
	totalNumber     uint64
	extraJsonLength string
	wg              *sync.WaitGroup
	request         *model.Request
}

func NewWorker(ch chan<- *model.RequestResults, chanId uint64, totalNumber uint64, extraJsonLength string, wg *sync.WaitGroup, request *model.Request, ws client.WebSocket) *WSWorker {
	return &WSWorker{
		WebSocket:       ws,
		chanId:          chanId,
		ch:              ch,
		totalNumber:     totalNumber,
		extraJsonLength: extraJsonLength,
		wg:              wg,
		request:         request,
	}
}

// web socket go link
func (w *WSWorker) Run(ctx context.Context) {
	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	defer func() {
		w.Close()
	}()

	var (
		i uint64
	)

	t := time.NewTimer(firstTime)
	for {
		select {
		case <-t.C:
			t.Reset(intervalTime)

			// 请求
			if ok := w.doRequest(i); !ok {
				keepAlive = false
				goto end
			}

			// 结束条件
			i = i + 1
			if i >= w.totalNumber {
				goto end
			}
		}
	}

end:
	w.wg.Done()
	t.Stop()

	if keepAlive {
		// 保持连接
		// chWaitFor := make(chan int, 0)
		// <-chWaitFor
		select {
		case <-ctx.Done():
			if w.chanId%100 == 0 {
				fmt.Printf("断开协程(%d)\n", w.chanId)
			}
			break
		}
	} else {
		fmt.Printf("异常断开协程(%d)\n", w.chanId)
	}

	return
}

// 请求
func (w *WSWorker) doRequest(i uint64) bool {

	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = model.HttpOk
	)

	// 需要发送的数据
	seq := fmt.Sprintf("%d_%d", w.chanId, i)
	// err := ws.Write([]byte(`{"seq":"` + seq + `","cmd":"ping","data":{}}`))
	// add extra data for larger bandwidth
	data := w.request.Body
	if len(w.extraJsonLength) > 0 {
		var builder strings.Builder
		idx := strings.Index(data, "}}")
		left := data[:idx+1]
		right := data[idx+1:]
		builder.WriteString(left)
		builder.WriteString(",\"extra\":\"")
		builder.WriteString(w.extraJsonLength)
		builder.WriteString("\"")
		builder.WriteString(right)
		data = builder.String()
	}

	//data = append(data, extra...)
	err := w.Write([]byte(data))
	if err != nil {
		errCode = model.RequestErr // 请求错误
		fmt.Printf("%v", err)
		return false
	} else {

		// time.Sleep(1 * time.Second)
		msg, err := w.Read()
		if err != nil {
			errCode = model.ParseError
			fmt.Println("读取数据 失败~")
		} else {
			// fmt.Println(msg)
			errCode, isSucceed = w.request.GetVerifyWebSocket()(w.request, seq, msg)
		}
	}

	requestTime := uint64(heper.DiffNano(startTime))

	requestResults := &model.RequestResults{
		Time:      requestTime,
		IsSucceed: isSucceed,
		ErrCode:   errCode,
	}

	requestResults.SetId(w.chanId, i)

	w.ch <- requestResults
	return true
}
