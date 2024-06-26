package emongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/util/xdebug"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	metricType = "mongo"
)

type Interceptor func(oldProcessFn processFn) (newProcessFn processFn)

func InterceptorChain(interceptors ...Interceptor) func(oldProcess processFn) processFn {
	build := func(interceptor Interceptor, oldProcess processFn) processFn {
		return interceptor(oldProcess)
	}

	return func(oldProcess processFn) processFn {
		chain := oldProcess
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain)
		}
		return chain
	}
}

func debugInterceptor(compName string, c *config) func(processFn) processFn {
	return func(oldProcess processFn) processFn {
		return func(cmd *cmd) error {
			if !eapp.IsDevelopmentMode() {
				return oldProcess(cmd)
			}

			beg := time.Now()
			err := oldProcess(cmd)
			cost := time.Since(beg)
			if err != nil {
				log.Println("emongo.response", xdebug.MakeReqAndResError(fileWithLineNum(), compName,
					fmt.Sprintf("%v", c.keyName), cost, fmt.Sprintf("%s %v", cmd.name, mustJsonMarshal(cmd.req)), err.Error()),
				)
			} else {
				log.Println("emongo.response", xdebug.MakeReqAndResInfo(fileWithLineNum(), compName,
					fmt.Sprintf("%v", c.keyName), cost, fmt.Sprintf("%s %v", cmd.name, mustJsonMarshal(cmd.req)), fmt.Sprintf("%v", cmd.res)),
				)
			}
			return err
		}
	}
}

func metricInterceptor(compName string, c *config, logger *elog.Component) func(processFn) processFn {
	return func(oldProcess processFn) processFn {
		return func(cmd *cmd) error {
			beg := time.Now()
			err := oldProcess(cmd)
			cost := time.Since(beg)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					emetric.ClientHandleCounter.Inc(metricType, compName, cmd.name, c.keyName, "Empty")
				} else {
					emetric.ClientHandleCounter.Inc(metricType, compName, cmd.name, c.keyName, "Error")
				}
			} else {
				emetric.ClientHandleCounter.Inc(metricType, compName, cmd.name, c.keyName, "OK")
			}
			emetric.ClientHandleHistogram.WithLabelValues(metricType, compName, cmd.name, c.keyName).Observe(cost.Seconds())
			return err
		}
	}
}

func accessInterceptor(compName string, c *config, logger *elog.Component) func(processFn) processFn {
	return func(oldProcess processFn) processFn {
		return func(cmd *cmd) error {
			beg := time.Now()
			err := oldProcess(cmd)
			cost := time.Since(beg)

			var fields = make([]elog.Field, 0, 15)
			fields = append(fields,
				elog.FieldMethod(cmd.name),
				elog.FieldCost(cost),
				elog.FieldKey(cmd.dbName),
				elog.String("collName", cmd.collName),
				elog.String("cmdName", cmd.name),
			)
			if c.EnableAccessInterceptorReq {
				fields = append(fields, elog.Any("req", cmd.req))
			}
			if c.EnableAccessInterceptorRes && err == nil {
				fields = append(fields, elog.Any("res", cmd.res))
			}
			event := "normal"
			isSlowLog := false
			if c.SlowLogThreshold > time.Duration(0) && cost > c.SlowLogThreshold {
				event = "slow"
				isSlowLog = true
			}

			if err != nil {
				fields = append(fields, elog.FieldEvent(event), elog.FieldErr(err))
				if errors.Is(err, mongo.ErrNoDocuments) {
					logger.Warn("access", fields...)
					return err
				}
				logger.Error("access", fields...)
				return err
			}

			if c.EnableAccessInterceptor || isSlowLog {
				fields = append(fields, elog.FieldEvent(event))
				if isSlowLog {
					logger.Warn("access", fields...)
				} else {
					logger.Info("access", fields...)
				}
			}
			return nil
		}
	}
}

func mustJsonMarshal(val interface{}) string {
	res, _ := json.Marshal(val)
	return string(res)
}

func fileWithLineNum() string {
	// the second caller usually from internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if (!(strings.Contains(file, "ego-component/emongo") && strings.HasSuffix(file, "wrapped_client.go")) && !(strings.Contains(file, "ego-component/emongo") && strings.Contains(file, "wrapped_collection.go"))) || strings.HasSuffix(file, "_test.go") {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
