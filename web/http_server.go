package web

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/browser"
)

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	runLocalhost()

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func RunHttpServer() {
	log.Printf("main: starting HTTP server")
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startHttpServer(httpServerExitDone)

	// program opens the browser to the following address
	go openLocalhost("http://localhost:8080/home")
	time.Sleep(1 * time.Second)

	// Program waits for input in the console before shutting down http server
	waitForInput()

	log.Printf("main: stopping HTTP server")

	// now close the server gracefully ("shutdown")
	// timeout could be given with a proper context
	// (in real world you shouldn't use TODO()).
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	// wait for goroutine started in startHttpServer() to stop
	httpServerExitDone.Wait()

	log.Printf("main: done. exiting")
}

func openLocalhost(url string) {
	browser.OpenURL(url)
}

func waitForInput() {
	var stop bool = false

	for !stop {
		log.Printf("main: running, type to stop")
		reader := bufio.NewReader(os.Stdin)
		// ReadString will block until the delimiter is entered
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("An error occured while reading input. Please try again", err)
			continue
		}
		if input != "" {
			stop = true
		}

	}
}
