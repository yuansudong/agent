package agent

// IsWindows shorthand function to check if OS == Windows
func (ua *UserAgent) IsWindows() bool {
	return ua.OS == _Windows
}

// IsAndroid shorthand function to check if OS == Android
func (ua *UserAgent) IsAndroid() bool {
	return ua.OS == _Android
}

// IsMacOS shorthand function to check if OS == MacOS
func (ua *UserAgent) IsMacOS() bool {
	return ua.OS == _MacOS
}

// IsIOS shorthand function to check if OS == IOS
func (ua *UserAgent) IsIOS() bool {
	return ua.OS == _IOS
}

// IsLinux shorthand function to check if OS == Linux
func (ua *UserAgent) IsLinux() bool {
	return ua.OS == _Linux
}

// IsOpera shorthand function to check if Name == Opera
func (ua *UserAgent) IsOpera() bool {
	return ua.Name == _Opera
}

// IsOperaMini shorthand function to check if Name == Opera Mini
func (ua *UserAgent) IsOperaMini() bool {
	return ua.Name == _OperaMini
}

// IsChrome shorthand function to check if Name == Chrome
func (ua *UserAgent) IsChrome() bool {
	return ua.Name == _Chrome
}

// IsFirefox shorthand function to check if Name == Firefox
func (ua *UserAgent) IsFirefox() bool {
	return ua.Name == _Firefox
}

// IsInternetExplorer shorthand function to check if Name == Internet Explorer
func (ua *UserAgent) IsInternetExplorer() bool {
	return ua.Name == _InternetExplorer
}

// IsSafari shorthand function to check if Name == Safari
func (ua *UserAgent) IsSafari() bool {
	return ua.Name == _Safari
}

// IsEdge shorthand function to check if Name == Edge
func (ua *UserAgent) IsEdge() bool {
	return ua.Name == _Edge
}

// IsGooglebot shorthand function to check if Name == Googlebot
func (ua *UserAgent) IsGooglebot() bool {
	return ua.Name == _Googlebot
}

// IsTwitterbot shorthand function to check if Name == Twitterbot
func (ua *UserAgent) IsTwitterbot() bool {
	return ua.Name == _Twitterbot
}

// IsFacebookbot shorthand function to check if Name == FacebookExternalHit
func (ua *UserAgent) IsFacebookbot() bool {
	return ua.Name == _FacebookExternalHit
}
