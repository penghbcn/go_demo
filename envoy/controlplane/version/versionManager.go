package version

import "strconv"

var Vsersion uint32 = 0

func GetNewVersion() string {
	Vsersion++
	return strconv.Itoa(int(Vsersion))
}
