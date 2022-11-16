package lservercom

//LicenseType license type
type LicenseType int

const (
	//Formal license type, include ip and not online
	Formal LicenseType = iota
	//SubLocal license type, sub company local, include ip and not online
	SubLocal
	//SubOnline license type, sub company online, include ip and online
	SubOnline
	//Test license type, not include ip and not online, 3 months
	Test
	//Trial license type, not include ip adn online
	Trial
	//LicenseTypeMaxLimit max limit of license type
	LicenseTypeMaxLimit
)

const (
	//FormalString license type string
	FormalString = "FORMAL"
	//SubLocalString license type string
	SubLocalString = "SUB_LOCAL"
	//SubOnlineString license type string
	SubOnlineString = "SUB_ONLINE"
	//TestString license type string
	TestString = "TEST"
	//TrialString license type string
	TrialString = "TRIAL"
)

func (ls LicenseType) String() string {
	switch ls {
	case Formal:
		return FormalString
	case SubLocal:
		return SubLocalString
	case SubOnline:
		return SubOnlineString
	case Test:
		return TestString
	case Trial:
		return TrialString
	default:
		return "Unknown license type"
	}
}
