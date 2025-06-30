// internal/logger/logger.go
package logger

import (
    "os"

    "github.com/sirupsen/logrus"
)

var Debug = logrus.New()

func init() {
    // Formato texto con timestamp
    Debug.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    Debug.SetLevel(logrus.DebugLevel)

    file, err := os.OpenFile("debug.log",
        os.O_CREATE|os.O_WRONLY|os.O_APPEND,
        0666,
    )
    if err == nil {
        Debug.Out = file
    }
}
