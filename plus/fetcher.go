package plus

import (
	"os/exec"
	"strconv"
)

func Fetch(position int) error {
	_, err := exec.Command("/bin/sh", strconv.Itoa(position)).Output()
	if err != nil {
		return err
	}

	// toReturn := string(rand.Intn(100)) + "," + string(rand.Intn(100))
	// fmt.Printf(toReturn)
	// out, err := exec.Command("/usr/sbin/python", filename, string(start)).Output()
	// if err != nil {
	// 	return "", err
	// }

	return nil
}
