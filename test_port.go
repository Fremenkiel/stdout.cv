package main
import (
	"fmt"
	"net"
)
func main() {
	l, err := net.Listen("tcp", ":")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Listening on", l.Addr().String())
	l.Close()
}
