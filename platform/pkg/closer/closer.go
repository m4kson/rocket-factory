package closer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"
)

const shutdownTimeout = 5 * time.Second

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func(context.Context) error
	log   *slog.Logger
}

var globalCloser = NewWithLogger(slog.Default())

func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

func SetLogger(log *slog.Logger) {
	globalCloser.SetLogger(log)
}

func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

func New(signals ...os.Signal) *Closer {
	return NewWithLogger(slog.Default(), signals...)
}

func NewWithLogger(log *slog.Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done: make(chan struct{}),
		log:  log,
	}
	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}
	return c
}

func (c *Closer) SetLogger(log *slog.Logger) {
	c.log = log
}

func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()

		c.log.InfoContext(ctx, "closing resource",
			slog.String("resource", name),
		)

		err := f(ctx)
		duration := time.Since(start)

		if err != nil {
			c.log.ErrorContext(ctx, "failed to close resource",
				slog.String("resource", name),
				slog.String("err", err.Error()),
				slog.Duration("duration", duration),
			)
		} else {
			c.log.InfoContext(ctx, "resource closed",
				slog.String("resource", name),
				slog.Duration("duration", duration),
			)
		}

		return err
	})
}

func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.log.InfoContext(ctx, "no resources to close")
			return
		}

		c.log.InfoContext(ctx, "starting graceful shutdown",
			slog.Int("resources", len(funcs)),
		)

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						err := fmt.Errorf("panic in closer func: %v", r)
						errCh <- err
						c.log.ErrorContext(ctx, "panic recovered in closer func",
							slog.Any("panic", r),
						)
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				c.log.WarnContext(ctx, "shutdown context cancelled",
					slog.String("err", ctx.Err().Error()),
				)
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.log.InfoContext(ctx, "graceful shutdown complete")
					return
				}
				c.log.ErrorContext(ctx, "error closing resource",
					slog.String("err", err.Error()),
				)
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}

func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case sig := <-ch:
		c.log.Info("received signal, starting graceful shutdown",
			slog.String("signal", sig.String()),
		)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.log.Error("graceful shutdown failed",
				slog.String("err", err.Error()),
			)
		}

	case <-c.done:

	}
}
