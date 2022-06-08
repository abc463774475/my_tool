package nlog

type Option interface {
	apply(l *loginfo)
}

type OptionFun func(l *loginfo)

func (f OptionFun) apply(l *loginfo) {
	f(l)
}

func WithOneFileMaxLines(maxLines int) Option {
	return OptionFun(func(l *loginfo) {
		l.oneFileMaxLines = maxLines
	})
}

func WithIsWriteLog(isWriteLog bool) Option {
	return OptionFun(func(l *loginfo) {
		l.isWriteLog = isWriteLog
	})
}

func WithCompressType(compressType CompressType) Option {
	return OptionFun(func(l *loginfo) {
		l.comressType = compressType
	})
}

func WithIsAsyn(isAsyn bool) Option {
	return OptionFun(func(l *loginfo) {
		l.isAsyn = isAsyn
	})
}
