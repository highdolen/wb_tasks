package main

// WaitForQuorum ждёт, пока N/2+1 серверов подтвердят успешную обработку
func WaitForQuorum(nodes int, results chan string) {
	quorum := nodes/2 + 1
	success := 0

	for res := range results {
		if res == "ok" {
			success++
		}
		if success >= quorum {
			break
		}
	}
}
