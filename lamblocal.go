package lamblocal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/exp/slog" // TODO: Go 1.21 may use "slog".
)

var Logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

// Run runs a lambda hander func detect the environment (lambda or not) and run it.
func Run[T any](ctx context.Context, fn func(context.Context, T) error) {
	if strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda") || os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambda.Start(fn)
	} else {
		if err := RunCLI(ctx, os.Stdin, fn); err != nil {
			Logger.Error(err.Error())
			os.Exit(1)
		}
	}
}

// RunCLI is a helper function for running a lambda hander func on CLI.
func RunCLI[T any](ctx context.Context, src io.Reader, fn func(context.Context, T) error) error {
	payload := new(T)
	if err := json.NewDecoder(src).Decode(payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}
	return fn(ctx, *payload)
}