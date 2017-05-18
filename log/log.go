package log

import (
	"net/http"

	"github.com/NeowayLabs/logger"
	"github.com/julienschmidt/httprouter"
)

type Logger struct {
	log *logger.Logger
}

func (self Logger) Debug(msg string) {
	self.log.Debug(msg)
}

func (self Logger) Debugf(msg string, args ...interface{}) {
	self.log.Debug(msg, args...)
}

func (self *Logger) Info(msg string) {
	self.log.Info(msg)
}

func (self Logger) Infof(msg string, args ...interface{}) {
	self.log.Info(msg, args...)
}

func (self Logger) Error(msg string) {
	self.log.Fatal(msg)
}

func (self Logger) Errorf(msg string, args ...interface{}) {
	self.log.Fatal(msg, args)
}

func (self Logger) Warn(msg string) {
	self.log.Warn(msg)
}

func (self Logger) Warnf(msg string, args ...interface{}) {
	self.log.Warn(msg, args...)
}

func (self Logger) ChangeLogLevel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.HTTPFunc(w, r)
}

func NewLogger(name string) Logger {

	return Logger{
		log: logger.Namespace(name),
	}
}
