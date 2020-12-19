package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	group := errgroup.Group{}
	// done 用于协程间退出的通知
	done := make(chan struct{})

	group.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case s := <-sigs:
			close(done)
			return fmt.Errorf("stoped by signal : %s", s.String())
		case <-done:
			return fmt.Errorf("stoped by another go routine")
		}
	})

	group.Go(func() error {
		// httpWork 仅用于通知当前协程 http.ListenAndServe 协程的退出
		httpWork := make(chan struct{})

		mux := http.NewServeMux()
		mux.HandleFunc("/index", func(writer http.ResponseWriter, request *http.Request) {
			io.WriteString(writer, "Hello, world!\n")
			// 退出通知
			close(httpWork)
			return
		})
		server := http.Server{
			Addr:    "127.0.0.1:24000",
			Handler: mux,
		}

		go func() {
			log.Println(server.ListenAndServe())
		}()

		select {
		case <-httpWork:
			close(done)
			return fmt.Errorf("stoped by http")
		case <-done:
			// 其他协程退出了，http server 也关闭
			server.Shutdown(context.TODO())
			return fmt.Errorf("stoped by another go routine")
		}
	})

	if err := group.Wait(); err != nil {
		log.Printf("error : %v", err)
	}
}
