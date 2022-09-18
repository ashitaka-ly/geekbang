package vlog

import "github.com/golang/glog"

func LogFatal(format string, args ...interface{}) {
	if args == nil {
		glog.V(1).Infof("FATAL-"+format, args)
	}
	glog.V(1).Infof("FATAL-"+format, args)
}

func LogError(format string, args ...interface{}) {
	glog.V(2).Info("ERROR-"+format, args)
}

func LogErrorf(format string, args ...interface{}) {
	glog.V(2).Infof("ERROR-"+format, args)
}

func LogWarning(format string, args ...interface{}) {
	glog.V(3).Infof("WARNING-"+format, args)
}

func LogInfo(format string, args ...interface{}) {
	glog.V(4).Infof("INFO-"+format, args)
}
