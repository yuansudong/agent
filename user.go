package agent

import (
	"bytes"
	"regexp"
	"strings"
)

// UserAgent 整个user_agent
type UserAgent struct {
	Name      string
	Version   string
	OS        string
	OSVersion string
	Device    string
	Mobile    bool
	Tablet    bool
	Desktop   bool
	Bot       bool
	URL       string
	String    string
}

// GetAgent 根据agent字符串,获得UserAgent结构
func GetAgent(sUserAgent string) *UserAgent {
	ua := &UserAgent{
		String: sUserAgent,
	}

	tokens := parse(sUserAgent)

	// check is there URL
	for k := range tokens {
		if strings.HasPrefix(k, "http://") || strings.HasPrefix(k, "https://") {
			ua.URL = k
			delete(tokens, k)
			break
		}
	}

	// OS lookup
	switch {
	case tokens.exists("Android"):
		ua.OS = _Android
		ua.OSVersion = tokens[_Android]
		for s := range tokens {
			if strings.HasSuffix(s, "Build") {
				ua.Device = strings.TrimSpace(s[:len(s)-5])
				ua.Tablet = strings.Contains(strings.ToLower(ua.Device), "tablet")
			}
		}

	case tokens.exists("iPhone"):
		ua.OS = _IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPhone"
		ua.Mobile = true

	case tokens.exists("iPad"):
		ua.OS = _IOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Device = "iPad"
		ua.Tablet = true

	case tokens.exists("Windows NT"):
		ua.OS = _Windows
		ua.OSVersion = tokens["Windows NT"]
		ua.Desktop = true

	case tokens.exists("Windows Phone OS"):
		ua.OS = _WindowsPhone
		ua.OSVersion = tokens["Windows Phone OS"]
		ua.Mobile = true

	case tokens.exists("Macintosh"):
		ua.OS = _MacOS
		ua.OSVersion = tokens.findMacOSVersion()
		ua.Desktop = true

	case tokens.exists("Linux"):
		ua.OS = _Linux
		ua.OSVersion = tokens[_Linux]
		ua.Desktop = true

	}
	switch {
	case tokens.exists("Googlebot"):
		ua.Name = _Googlebot
		ua.Version = tokens[_Googlebot]
		ua.Bot = true
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Opera Mini"] != "":
		ua.Name = _OperaMini
		ua.Version = tokens[_OperaMini]
		ua.Mobile = true

	case tokens["OPR"] != "":
		ua.Name = _Opera
		ua.Version = tokens["OPR"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["OPT"] != "":
		ua.Name = _OperaTouch
		ua.Version = tokens["OPT"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Opera on iOS
	case tokens["OPiOS"] != "":
		ua.Name = _Opera
		ua.Version = tokens["OPiOS"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Chrome on iOS
	case tokens["CriOS"] != "":
		ua.Name = _Chrome
		ua.Version = tokens["CriOS"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// Firefox on iOS
	case tokens["FxiOS"] != "":
		ua.Name = _Firefox
		ua.Version = tokens["FxiOS"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["Firefox"] != "":
		ua.Name = _Firefox
		ua.Version = tokens[_Firefox]
		_, ua.Mobile = tokens["Mobile"]
		_, ua.Tablet = tokens["Tablet"]

	case tokens["Vivaldi"] != "":
		ua.Name = _Vivaldi
		ua.Version = tokens[_Vivaldi]

	case tokens.exists("MSIE"):
		ua.Name = _InternetExplorer
		ua.Version = tokens["MSIE"]

	case tokens["Edge"] != "":
		ua.Name = _Edge
		ua.Version = tokens["Edge"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["EdgA"] != "":
		ua.Name = _Edge
		ua.Version = tokens["EdgA"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["bingbot"] != "":
		ua.Name = "Bingbot"
		ua.Version = tokens["bingbot"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens["SamsungBrowser"] != "":
		ua.Name = "Samsung Browser"
		ua.Version = tokens["SamsungBrowser"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	// if chrome and Safari defined, find any other tokensent descr
	case tokens.exists(_Chrome) && tokens.exists(_Safari):
		name := tokens.findBestMatch(true)
		if name != "" {
			ua.Name = name
			ua.Version = tokens[name]
			break
		}
		fallthrough

	case tokens.exists("Chrome"):
		ua.Name = _Chrome
		ua.Version = tokens["Chrome"]
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	case tokens.exists("Safari"):
		ua.Name = _Safari
		if v, ok := tokens["Version"]; ok {
			ua.Version = v
		} else {
			ua.Version = tokens["Safari"]
		}
		ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")

	default:
		if ua.OS == "Android" && tokens["Version"] != "" {
			ua.Name = "Android browser"
			ua.Version = tokens["Version"]
			ua.Mobile = true
		} else {
			if name := tokens.findBestMatch(false); name != "" {
				ua.Name = name
				ua.Version = tokens[name]
			} else {
				ua.Name = ua.String
			}
			ua.Bot = strings.Contains(strings.ToLower(ua.Name), "bot")
			ua.Mobile = tokens.existsAny("Mobile", "Mobile Safari")
		}
	}
	if ua.Tablet {
		ua.Mobile = false
	}
	if !ua.Bot {
		ua.Bot = ua.URL != ""
	}

	if !ua.Bot {
		switch ua.Name {
		case _Twitterbot, _FacebookExternalHit:
			ua.Bot = true
		}
	}

	return ua
}

func parse(userAgent string) (clients properties) {
	clients = make(map[string]string, 0)
	slash := false
	isURL := false
	var buff, val bytes.Buffer
	addToken := func() {
		if buff.Len() != 0 {
			s := strings.TrimSpace(buff.String())
			if _, ign := _Ignore[s]; !ign {
				if isURL {
					s = strings.TrimPrefix(s, "+")
				}

				if val.Len() == 0 { // only if value don't exists
					var ver string
					s, ver = checkVer(s) // determin version string and split
					clients[s] = ver
				} else {
					clients[s] = strings.TrimSpace(val.String())
				}
			}
		}
		buff.Reset()
		val.Reset()
		slash = false
		isURL = false
	}

	parOpen := false

	bua := []byte(userAgent)
	for i, c := range bua {

		//fmt.Println(string(c), c)
		switch {
		case c == 41: // )
			addToken()
			parOpen = false

		case parOpen && c == 59: // ;
			addToken()

		case c == 40: // (
			addToken()
			parOpen = true

		case slash && c == 32:
			addToken()

		case slash:
			val.WriteByte(c)

		case c == 47 && !isURL: //   /
			if i != len(bua)-1 && bua[i+1] == 47 && (bytes.HasSuffix(buff.Bytes(), []byte("http:")) || bytes.HasSuffix(buff.Bytes(), []byte("https:"))) {
				buff.WriteByte(c)
				isURL = true
			} else {
				slash = true
			}

		default:
			buff.WriteByte(c)
		}
	}
	addToken()

	return clients
}

func checkVer(s string) (name, v string) {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return s, ""
	}

	//v = s[i+1:]

	switch s[:i] {
	case "Linux", "Windows NT", "Windows Phone OS", "MSIE", "Android":
		return s[:i], s[i+1:]
	default:
		return s, ""
	}
}

type properties map[string]string

func (p properties) exists(key string) bool {
	_, ok := p[key]
	return ok
}

func (p properties) existsAny(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; ok {
			return true
		}
	}
	return false
}

func (p properties) findMacOSVersion() string {
	for k, v := range p {
		if strings.Contains(k, "OS") {
			if ver := findVersion(v); ver != "" {
				return ver
			} else if ver = findVersion(k); ver != "" {
				return ver
			}
		}
	}
	return ""
}

// findBestMatch from the rest of the bunch
// in first cycle only return key vith version value
// if withVerValue is false, do another cycle and return any token
func (p properties) findBestMatch(withVerOnly bool) string {
	n := 2
	if withVerOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		for k, v := range p {
			switch k {
			case _Chrome, _Firefox, _Safari, "Version", "Mobile", "Mobile Safari", "Mozilla", "AppleWebKit", "Windows NT", "Windows Phone OS", _Android, "Macintosh", _Linux, "GSA":
			default:
				if i == 0 {
					if v != "" { // in first check, only return  keys with value
						return k
					}
				} else {
					return k
				}
			}
		}
	}
	return ""
}

var rxMacOSVer = regexp.MustCompile("[_\\d\\.]+")

func findVersion(s string) string {
	if ver := rxMacOSVer.FindString(s); ver != "" {
		return strings.Replace(ver, "_", ".", -1)
	}
	return ""
}
