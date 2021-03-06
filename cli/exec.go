package cli

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func ExecFile(pathToFile string, opts Options) {
	query, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	Exec(string(query), opts)
}

func Exec(query string, opts Options) {
	ferret := compiler.New()

	prog, err := ferret.Compile(query)

	if err != nil {
		fmt.Println("Failed to compile the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	l := NewLogger()

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			cancel()
			l.Close()
		}
	}()

	out, err := prog.Run(
		ctx,
		runtime.WithBrowser(opts.Cdp),
		runtime.WithLog(l),
		runtime.WithLogLevel(logging.DebugLevel),
		runtime.WithParams(opts.Params),
		runtime.WithProxy(opts.Proxy),
		runtime.WithUserAgent(opts.UserAgent),
	)

	if err != nil {
		fmt.Println("Failed to execute the query")
		fmt.Println(err)
		os.Exit(1)
		return
	}

	fmt.Println(string(out))
}
