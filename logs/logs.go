package logs

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func CtxInfo(ctx context.Context, message string, args ...interface{}) {
	log.WithContext(ctx).Info(fmt.Sprintf("%s %s", message, args))
}

func CtxError(ctx context.Context, message string, args ...interface{}) {
	log.WithContext(ctx).Error(fmt.Sprintf("%s %s", message, args))
}

func CtxDebug(ctx context.Context, message string, args ...interface{}) {
	log.WithContext(ctx).Debug(fmt.Sprintf("%s %s", message, args))
}

func CtxWarn(ctx context.Context, message string, args ...interface{}) {
	log.WithContext(ctx).Warn(fmt.Sprintf("%s %s", message, args))
}
