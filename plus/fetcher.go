package plus

import (
	"math/rand"
)

func Fetch(iface string) (string, error) {
	toReturn := string(rand.Intn(100)) + "," + string(rand.Intn(100))
	// fmt.Printf(toReturn)
	// out, err := exec.Command("/usr/sbin/python", filename, string(start)).Output()
	// if err != nil {
	// 	return "", err
	// }

	return toReturn, nil
}
