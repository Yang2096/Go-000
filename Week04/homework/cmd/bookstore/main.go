package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"week4/api/bookstore"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type bookstoreServer struct {
	bookstore.UnimplementedBookstoreServer
}

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

		lis, err := net.Listen("tcp", "20000")
		s := grpc.NewServer()

		bookstore.RegisterBookstoreServer(s, &bookstoreServer{})
		if err = s.Serve(lis); err != nil {
			close(httpWork)
		}
		select {
		case <-httpWork:
			close(done)
			return fmt.Errorf("stoped by http")
		case <-done:
			// 其他协程退出了，http server 也关闭
			s.Stop()
			return fmt.Errorf("stoped by another go routine")
		}
	})

	if err := group.Wait(); err != nil {
		log.Printf("error : %v", err)
	}
}
