package or

import (
	"fmt"
	"time"
)

// ExampleOr - пример использования
func ExampleOr() {

	start := time.Now()

	<-Or(
		Sig(2*time.Hour),
		Sig(5*time.Minute),
		Sig(1*time.Second),
		Sig(1*time.Hour),
		Sig(1*time.Minute),
	)

	fmt.Printf("done after %v\n", time.Since(start).Round(time.Second))

	// Output:
	// done after 1s

}
