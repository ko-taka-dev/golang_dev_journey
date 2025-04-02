package logger

import (
    "log"
    "os"
    "sync"
)

// ログレベルを定義
const (
    DEBUG = iota
    INFO
    WARN
    ERROR
    FATAL
)

var (
    once   sync.Once
    logger *Logger
)

// Logger はアプリケーションのロギングを担当
type Logger struct {
    debugLogger *log.Logger
    infoLogger  *log.Logger
    warnLogger  *log.Logger
    errorLogger *log.Logger
    fatalLogger *log.Logger
    level       int
}

// GetLogger はシングルトンのLoggerインスタンスを返す
func GetLogger() *Logger {
    once.Do(func() {
        logger = &Logger{
            debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
            infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
            warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime),
            errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
            fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
            level:       INFO, // デフォルトはINFOレベル
        }
    })
    return logger
}

// SetLevel はロガーのログレベルを設定
func (l *Logger) SetLevel(level int) {
    l.level = level
}

// Debug はデバッグレベルのログを出力
func (l *Logger) Debug(v ...interface{}) {
    if l.level <= DEBUG {
        l.debugLogger.Println(v...)
    }
}

// Debugf はフォーマット付きのデバッグレベルのログを出力
func (l *Logger) Debugf(format string, v ...interface{}) {
    if l.level <= DEBUG {
        l.debugLogger.Printf(format, v...)
    }
}

// Info は情報レベルのログを出力
func (l *Logger) Info(v ...interface{}) {
    if l.level <= INFO {
        l.infoLogger.Println(v...)
    }
}

// Infof はフォーマット付きの情報レベルのログを出力
func (l *Logger) Infof(format string, v ...interface{}) {
    if l.level <= INFO {
        l.infoLogger.Printf(format, v...)
    }
}

// Warn は警告レベルのログを出力
func (l *Logger) Warn(v ...interface{}) {
    if l.level <= WARN {
        l.warnLogger.Println(v...)
    }
}

// Warnf はフォーマット付きの警告レベルのログを出力
func (l *Logger) Warnf(format string, v ...interface{}) {
    if l.level <= WARN {
        l.warnLogger.Printf(format, v...)
    }
}

// Error はエラーレベルのログを出力
func (l *Logger) Error(v ...interface{}) {
    if l.level <= ERROR {
        l.errorLogger.Println(v...)
    }
}

// Errorf はフォーマット付きのエラーレベルのログを出力
func (l *Logger) Errorf(format string, v ...interface{}) {
    if l.level <= ERROR {
        l.errorLogger.Printf(format, v...)
    }
}

// Fatal は致命的なエラーレベルのログを出力し、プログラムを終了する
func (l *Logger) Fatal(v ...interface{}) {
    if l.level <= FATAL {
        l.fatalLogger.Println(v...)
        os.Exit(1)
    }
}

// Fatalf はフォーマット付きの致命的なエラーレベルのログを出力し、プログラムを終了する
func (l *Logger) Fatalf(format string, v ...interface{}) {
    if l.level <= FATAL {
        l.fatalLogger.Printf(format, v...)
        os.Exit(1)
    }
}