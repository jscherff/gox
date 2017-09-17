package main
import "github.com/jscherff/gox/log"
func main() {
	//ml := log.NewMLogger("test", log.LstdFlags, true, false, "test1.log", "test2.log")
	ml := log.NewMLogger("test", log.LstdFlags, true, false, "test1.log", "test2.log")
	ml.Println("This is a test")
	ml.Write([]byte("This is a second test"))
	ml.AddFile("test3.log")
	ml.Println("This is a third test")
	ml.SetPrefix("NEWPREFIX")
	ml.Write([]byte("This is a fourth test"))
	ml.SetStderr(true)
	ml.Println("This is a final test")
	ml.Close()
}
