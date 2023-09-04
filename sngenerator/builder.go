package sngenerator

type builder struct{}

var B = builder{}

func (b builder) Const(s string) TConst {
	return TConst(s)
}

func (b builder) Env(name ExpVal) TEnv {
	return TEnv{Name: name}
}

func (b builder) Code(exp ExpVal) *TCode {
	return &TCode{Exp: exp, OnlyCode: false}
}
func (b builder) CodeString(exp string) *TCode {
	return &TCode{Exp: TConst(exp), OnlyCode: false}
}
func (b builder) CodeEnv(name string) *TCode {
	return &TCode{Exp: TEnv{Name: TConst(name)}, OnlyCode: false}
}

// CodeOnly 仅业务代号
func (b builder) CodeOnly(exp ExpVal) *TCode {
	return &TCode{Exp: exp, OnlyCode: true}
}
func (b builder) CodeOnlyString(exp string) *TCode {
	return &TCode{Exp: TConst(exp), OnlyCode: true}
}

// Incr 计数器
func (b builder) Incr(min, step int64) *TIncr {
	if step == 0 {
		step = 1
	}
	return &TIncr{Min: min, Step: step}
}

func (b builder) IncrRightZero(min int64, len int) *TIncr {
	return &TIncr{Min: min, Step: 1, Len: len, Align: AlignRight, Pad: "0"}
}

// Rand 随机数
func (b builder) Rand(len int, Chr string) TRand {
	return TRand{Len: len, Chr: Chr}
}

func (b builder) RandDigit(len int) TRand {
	return TRand{Len: len, Chr: Digit}
}

func (b builder) RandAlpha(len int) TRand {
	return TRand{Len: len, Chr: Alpha}
}

func (b builder) RandHex(len int) TRand {
	return TRand{Len: len, Chr: Hex}
}

func (b builder) RandAlphaNum(len int) TRand {
	return TRand{Len: len, Chr: AlphaNum}
}

// RandFill 随机填充在右侧
func (b builder) RandFill(exp ExpVal, len int, chr string) TRandFill {
	return TRandFill{Exp: exp, Len: len, Chr: chr}
}

// RandFillDigit 随机填充在右侧
func (b builder) RandFillDigit(exp ExpVal, len int) TRandFill {
	return TRandFill{Exp: exp, Len: len, Chr: Digit}
}

// RandFillAlpha 随机填充在右侧
func (b builder) RandFillAlpha(exp ExpVal, len int) TRandFill {
	return TRandFill{Exp: exp, Len: len, Chr: Alpha}
}

// RandFillHex 随机填充在右侧
func (b builder) RandFillHex(exp ExpVal, len int) TRandFill {
	return TRandFill{Exp: exp, Len: len, Chr: Hex}
}

// RandFillAlphaNum 随机填充在右侧
func (b builder) RandFillAlphaNum(exp ExpVal, len int) TRandFill {
	return TRandFill{Exp: exp, Len: len, Chr: AlphaNum}
}

// RandFillLeft 随机填充在左侧
func (b builder) RandFillLeft(exp ExpVal, len int, chr string) TRandFill {
	return TRandFill{Exp: exp, Len: len, Align: AlignRight, Chr: chr}
}

// RandFillLeftDigit 随机填充在左侧
func (b builder) RandFillLeftDigit(exp ExpVal, len int) TRandFill {
	return TRandFill{Exp: exp, Len: len, Align: AlignRight, Chr: Digit}
}

// RandFillLeftAlpha 随机填充在左侧
func (b builder) RandFillLeftAlpha(exp ExpVal, len int) TRandFill {
	return TRandFill{Exp: exp, Len: len, Align: AlignRight, Chr: Alpha}
}

func (b builder) Time(format string) TTime {
	return TTime{Format: format}
}

func (b builder) TimeYear() TTime {
	return TTime{Format: DateFmtYear}
}

func (b builder) TimeMonth() TTime {
	return TTime{Format: DateFmtMonth}
}

func (b builder) TimeDay() TTime {
	return TTime{Format: DateFmtDay}
}

func (b builder) TimeWithCode(format, formatCode string) TTime {
	return TTime{Format: format, FormatCode: formatCode, IsCode: true}
}

func (b builder) Join(sep string, args ...ExpVal) *TJoin {
	return &TJoin{Sep: sep, Args: args}
}

func (b builder) Session(exp ExpVal) TSession {
	return TSession{Exp: exp}
}
