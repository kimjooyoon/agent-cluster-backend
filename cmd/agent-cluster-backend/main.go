// agent-cluster-backend: placeholder entrypoint.
//
// The real server (GraphQL + SSE, consumer of generated contract clients)
// arrives with the first vertical-slice decision after 002-dumb-agent-role.
package main

import "fmt"

func main() {
	fmt.Println("agent-cluster-backend: bootstrap-minimum build")
	fmt.Println("contracts SSOT: https://github.com/kimjooyoon/agent-cluster-contracts")
	fmt.Println("Real server arrives with the first vertical-slice decision.")
}
