package farmhash

import "fmt"

func ExampleHash32() {
	str := "hello world"
	bytes := []byte(str)
	hash := Hash32(bytes)
	fmt.Printf("Hash32(%s) is %x\n", str, hash)
}
