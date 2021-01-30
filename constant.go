package agent

// _Ignore 忽略的字符
var _Ignore = map[string]struct{}{
	"KHTML, like Gecko": {},
	"U":                 {},
	"compatible":        {},
	"Mozilla":           {},
	"WOW64":             {},
}

// 常量,列举
const (
	_Windows             = "Windows"
	_WindowsPhone        = "Windows Phone"
	_Android             = "Android"
	_MacOS               = "macOS"
	_IOS                 = "iOS"
	_Linux               = "Linux"
	_Opera               = "Opera"
	_OperaMini           = "Opera Mini"
	_OperaTouch          = "Opera Touch"
	_Chrome              = "Chrome"
	_Firefox             = "Firefox"
	_InternetExplorer    = "Internet Explorer"
	_Safari              = "Safari"
	_Edge                = "Edge"
	_Vivaldi             = "Vivaldi"
	_Googlebot           = "Googlebot"
	_Twitterbot          = "Twitterbot"
	_FacebookExternalHit = "facebookexternalhit"
)
