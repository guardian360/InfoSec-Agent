// Package apiconnection provides functions for integrating the application with the Guardian360 Lighthouse platform.
package apiconnection

type IssueResPair struct {
	IssueID  int
	ResultID int
}

// TODO: Update documentation
// IssueMap is a map that maps an issue ID and result ID pair to a boolean value.
// This map is used to determine whether a given issue is considered a problem based on the result of a security or privacy check.
var IssueMap = map[IssueResPair]bool{
	// Bluetooth: No devices found
	{1, 0}: false,
	// Bluetooth: Devices found
	{1, 1}: false,
	// External Devices: No devices found
	{2, 0}: false,
	// External Devices: Devices found
	{2, 1}: false,
	// Guest Account: Guest account not found
	{3, 0}: false,
	// Guest Account: Guest account active
	{3, 1}: true,
	// Guest Account: Guest account not active
	{3, 2}: false,
	// Network Profile Type: No network profiles found
	{4, 0}: false,
	// Network Profile Type: Network profiles found
	{4, 1}: false,
	// Password Manager: Password manager found
	{5, 0}: false,
	// Password Manager: No password manager found
	{5, 1}: true,
	// Location: Apps with location access
	{6, 0}: false,
	// Location: No apps with location access
	{6, 1}: false,
	// Microphone: Apps with microphone access
	{7, 0}: false,
	// Microphone: No apps with microphone access
	{7, 1}: false,
	// Webcam: Apps with webcam access
	{8, 0}: false,
	// Webcam: No apps with webcam access
	{8, 1}: false,
	// Appointments: Apps with calendar access
	{9, 0}: false,
	// Appointments: No apps with calendar access
	{9, 1}: false,
	// Contacts: Apps with contacts access
	{10, 0}: false,
	// Contacts: No apps with contacts access
	{10, 1}: false,
	// Ports: Open ports found
	{11, 0}: false,
	// Remote Desktop: Remote Desktop enabled
	{12, 0}: true,
	// Remote Desktop: Remote Desktop disabled
	{12, 1}: false,
	// SMB: smb1 and smb2 disabled
	{13, 0}: true,
	// SMB: smb1 enabled, smb2 disabled
	{13, 1}: true,
	// SMB: smb1 disabled, smb2 enabled
	{13, 2}: false,
	// SMB: smb1 and smb2 enabled
	{13, 3}: true,
	// UAC: UAC disabled
	{14, 0}: true,
	// UAC: UAC enabled for apps and settings
	{14, 1}: false,
	// UAC: UAC enabled for apps
	{14, 2}: false,
	// UAC: Unknown UAC level
	{14, 3}: false,
	// Defender: Real-time enabled
	{15, 0}: false,
	// Defender: Real-time disabled
	{15, 1}: true,
	// Last Password Change: Password changed more than half a year ago
	{16, 0}: true,
	// Last Password Change: Password changed recently
	{16, 1}: false,
	// Login Method: PIN only
	{17, 1}: true,
	// Login Method: Picture Logon only
	{17, 2}: true,
	// Login Method: PIN and Picture Logon
	{17, 3}: true,
	// Login Method: Password only
	{17, 4}: false,
	// Login Method: PIN and Password
	{17, 5}: true,
	// Login Method: Picture Logon and Password
	{17, 6}: true,
	// Login Method: PIN, Picture Logon, and Password
	{17, 7}: true,
	// Login Method: Fingerprint only
	{17, 8}: false,
	// Login Method: PIN and Fingerprint
	{17, 9}: true,
	// Login Method: Picture Logon and Fingerprint
	{17, 10}: true,
	// Login Method: PIN, Picture Logon, and Fingerprint
	{17, 11}: true,
	// Login Method: Password and Fingerprint
	{17, 12}: false,
	// Login Method: PIN, Password, and Fingerprint
	{17, 13}: true,
	// Login Method: Picture Logon, Password, and Fingerprint
	{17, 14}: true,
	// Login Method: PIN, Picture Logon, Password, and Fingerprint
	{17, 15}: true,
	// Login Method: Facial recognition only
	{17, 16}: true,
	// Login Method: PIN and Facial recognition
	{17, 17}: true,
	// Login Method: Picture Logon and Facial recognition
	{17, 18}: true,
	// Login Method: PIN, Picture Logon, and Facial recognition
	{17, 19}: true,
	// Login Method: Password and Facial recognition
	{17, 20}: true,
	// Login Method: PIN, Password, and Facial recognition
	{17, 21}: true,
	// Login Method: Picture Logon, Password, and Facial recognition
	{17, 22}: true,
	// Login Method: PIN, Picture Logon, Password, and Facial recognition
	{17, 23}: true,
	// Login Method: Fingerprint and Facial recognition
	{17, 24}: true,
	// Login Method: PIN, Fingerprint, and Facial recognition
	{17, 25}: true,
	// Login Method: Picture Logon, Fingerprint, and Facial recognition
	{17, 26}: true,
	// Login Method: PIN, Picture Logon, and Fingerprint, and Facial recognition
	{17, 27}: true,
	// Login Method: Password, Fingerprint, and Facial recognition
	{17, 28}: true,
	// Login Method: PIN, Password, Fingerprint, and Facial recognition
	{17, 29}: true,
	// Login Method: Picture Logon, Password, Fingerprint, and Facial recognition
	{17, 30}: true,
	// Login Method: PIN, Picture Logon, Password, Fingerprint, and Facial recognition
	{17, 31}: true,
	// Login Method: Trust signal only
	{17, 32}: true,
	// Login Method: PIN and Trust signal
	{17, 33}: true,
	// Login Method: Picture Logon and Trust signal
	{17, 34}: true,
	// Login Method: PIN, Picture Logon, and Trust signal
	{17, 35}: true,
	// Login Method: Password and Trust signal
	{17, 36}: true,
	// Login Method: PIN, Password, and Trust signal
	{17, 37}: true,
	// Login Method: Picture Logon, Password, and Trust signal
	{17, 38}: true,
	// Login Method: PIN, Picture Logon, Password, and Trust signal
	{17, 39}: true,
	// Login Method: Fingerprint and Trust signal
	{17, 40}: true,
	// Login Method: PIN, Fingerprint, and Trust signal
	{17, 41}: true,
	// Login Method: Picture Logon, Fingerprint, and Trust signal
	{17, 42}: true,
	// Login Method: PIN, Picture Logon, and Fingerprint, and Trust signal
	{17, 43}: true,
	// Login Method: Password, Fingerprint, and Trust signal
	{17, 44}: true,
	// Login Method: PIN, Password, Fingerprint, and Trust signal
	{17, 45}: true,
	// Login Method: Picture Logon, Password, Fingerprint, and Trust signal
	{17, 46}: true,
	// Login Method: PIN, Picture Logon, Password, Fingerprint, and Trust signal
	{17, 47}: true,
	// Login Method: Facial recognition and Trust signal
	{17, 48}: true,
	// Login Method: PIN, Facial recognition, and Trust signal
	{17, 49}: true,
	// Login Method: Picture Logon, Facial recognition, and Trust signal
	{17, 50}: true,
	// Login Method: PIN, Picture Logon, and Facial recognition, and Trust signal
	{17, 51}: true,
	// Login Method: Password, Facial recognition, and Trust signal
	{17, 52}: true,
	// Login Method: PIN, Password, Facial recognition, and Trust signal
	{17, 53}: true,
	// Login Method: Picture Logon, Password, Facial recognition, and Trust signal
	{17, 54}: true,
	// Login Method: PIN, Picture Logon, Password, Facial recognition, and Trust signal
	{17, 55}: true,
	// Login Method: Fingerprint, Facial recognition, and Trust signal
	{17, 56}: true,
	// Login Method: PIN, Fingerprint, Facial recognition, and Trust signal
	{17, 57}: true,
	// Login Method: Picture Logon, Fingerprint, Facial recognition, and Trust signal
	{17, 58}: true,
	// Login Method: PIN, Picture Logon, Fingerprint, Facial recognition, and Trust signal
	{17, 59}: true,
	// Login Method: Password, Fingerprint, Facial recognition, and Trust signal
	{17, 60}: true,
	// Login Method: PIN, Password, Fingerprint, Facial recognition, and Trust signal
	{17, 61}: true,
	// Login Method: Picture Logon, Password, Fingerprint, Facial recognition, and Trust signal
	{17, 62}: true,
	// Login Method: PIN, Picture Logon, Password, Fingerprint, Facial recognition, and Trust signal
	{17, 63}: true,
	// Outdated: Windows up to date
	{18, 0}: false,
	// Outdated: Windows update available
	{18, 1}: true,
	// Outdated: Windows version not 10 or 11
	{18, 2}: false,
	// Secure Boot: Secure Boot enabled
	{19, 0}: true,
	// Secure Boot: Secure Boot disabled
	{19, 1}: false,
	// Secure Boot: Status unknown
	{19, 2}: false,
	// Startup: No startup programs found
	{20, 0}: false,
	// Startup: Startup programs found
	{20, 1}: false,
	// Chromium Extension: Ad blocker installed
	{21, 0}: false,
	// Chromium Extension: No ad blocker installed
	{21, 1}: true,
	// Edge Extension: Ad blocker installed
	{22, 0}: false,
	// Edge Extension: No ad blocker installed
	{22, 1}: true,
	// Chromium History: Phishing domains found
	{23, 0}: true,
	// Chromium History: No phishing domains found
	{23, 1}: false,
	// Edge History: Phishing domains found
	{24, 0}: true,
	// Edge History: No phishing domains found
	{24, 1}: false,
	// Chromium SearchEngine: Search engine
	{25, 0}: false,
	// Edge SearchEngine: Search engine,
	{26, 0}: false,
	// Firefox Cookies: No tracking cookies found
	{27, 0}: false,
	// Firefox Cookies: Tracking cookies found
	{27, 1}: true,
	// Firefox Extension: List of extensions
	{28, 0}: false,
	// Firefox Adblock: Ad blocker installed
	{29, 0}: false,
	// Firefox Adblock: No ad blocker installed
	{29, 1}: true,
	// Firefox SearchEngine: Search engine
	{30, 0}: false,
	// Firefox History: No phishing domains found
	{31, 0}: false,
	// Firefox History: Phishing domains found
	{31, 1}: true,
	// CIS Registry Settings: Not everything is set correctly
	{32, 0}: true,
	// CIS Registry Settings: Everything is set correctly
	{32, 1}: false,
	// Auto Login: Auto login disabled
	{33, 0}: false,
	// Auto Login: Auto login enabled
	{33, 1}: true,
	// Remote RPC: Remote RPC disabled
	{34, 0}: false,
	// Remote RPC: Remote RPC enabled
	{34, 1}: true,
	// Chromium Cookies: No tracking cookies found
	{35, 0}: false,
	// Chromium Cookies: Tracking cookies found
	{35, 1}: true,
	// Edge Cookies: No tracking cookies found
	{36, 0}: false,
	// Edge Cookies: Tracking cookies found
	{36, 1}: true,
	// Windows Firewall: Firewall enabled for all profiles
	{37, 0}: false,
	// Windows Firewall: Firewall disabled for any/all profile(s)
	{37, 1}: true,
	// Password length: Password length at least 15 characters
	{38, 0}: false,
	// Password length: Password length less than 15 characters
	{38, 1}: true,
	// Credential Guard: Credential Guard running
	{39, 0}: false,
	// Credential Guard: Credential Guard not running
	{39, 1}: true,
	// NetBIOS over TCP/IP: NetBIOS disabled
	{40, 0}: false,
	// NetBIOS over TCP/IP: NetBIOS enabled
	{40, 1}: true,
	// WPAD: WPAD service disabled
	{41, 0}: false,
	// WPAD: WPAD service enabled
	{41, 1}: true,
	// Screen Lock: Screen lock correctly enabled
	{42, 0}: false,
	// Screen Lock: Screen lock not correctly enabled
	{42, 1}: true,
}
