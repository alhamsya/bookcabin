package public

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func RunRESTApp(ctx context.Context) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	/* === GENERAL === */
	/* === HANDLER === */
	/* === DATABASE === */
	/* === REDIS === */

	return nil
}
