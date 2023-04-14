package server

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/jithin-kg/webpa-common/logging"
)

// SignalWait blocks until any of a set of signals is encountered.  The signal which caused this function
// to exit is returned.  A nil return indicates that the signals channel was closed.
//
// If no waitOn signals are supplied, this function will never return until the signals channel is closed.
//
// In all cases, the supplied logger is used to log information about signals that are ignored.
/**
@author jithin-kg
function takes a logger of type log.Logger
 channel of OS signals called signals, and an arbitrary number of os.Signal types to wait for called waitOn.
 it returns os.Signal type
 so signals channel will be passing value of type os.Signal
 and we have a list of signals that we are interested in waitOn slice
in first loop we create a map of signals we are interested in and we assign true to it
in the second for loop;->The function then enters an infinite loop and waits for signals to be received on the 'signals'channel.
whenever an os signal receive and if that signal is present in our interested signal we just return the signal else we ignore and log
**/
func SignalWait(logger log.Logger, signals <-chan os.Signal, waitOn ...os.Signal) os.Signal {
	// here we create a map of type os.signl, boolean
	filter := make(map[os.Signal]bool)
	for _, s := range waitOn {
		filter[s] = true
	}

	for s := range signals {
		if filter[s] {
			return s
		}

		logger.Log(logging.MessageKey(), "ignoring signal", "signal", s.String())
	}

	return nil
}
