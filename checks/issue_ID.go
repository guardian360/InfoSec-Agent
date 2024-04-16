package checks

// This is a list of all the Result IDs for the checks that are performed. It starts at 1 and then iterates up.
const (
	BluetoothID int = iota + 1
	ExternalDevicesID
	GuestAccountID
	PasswordManagerID
	LocationID
	MicrophoneID
	WebcamID
	AppointmentsID
	ContactsID
	PortsID
	RemoteDesktopID
	SmbID
	UacID
	WindowsDefenderID
	LastPasswordChangeID
	LoginMethodID
	WindowsOutdatedID
	SecureBootID
	StartupID
	ExtensionChromiumID
	ExtensionEdgeID
	HistoryChromiumID
	HistoryEdgeID
	SearchChromiumID
	SearchEdgeID
	CookiesFirefoxID
	ExtensionFirefoxID
	AdblockFirefoxID
	SearchFirefoxID
	HistoryFirefoxID
	NetworkProfileTypeID
)
