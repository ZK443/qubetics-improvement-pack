// SPDX-License-Identifier: MIT
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
)

func main() {
	cmd := flag.String("cmd", "", "send|verify|execute")
	id := flag.String("id", "", "message ID")
	nonce := flag.Uint64("nonce", 0, "nonce (for send)")
	src := flag.String("src", "eth", "source chain")
	dst := flag.String("dst", "qub", "dest chain")
	flag.Parse()

	switch *cmd {
	case "send":
		fmt.Printf("SEND id=%s nonce=%d src=%s dst=%s\n", *id, *nonce, *src, *dst)
	case "verify":
		fmt.Printf("VERIFY id=%s proof=%s\n", *id, hex.EncodeToString([]byte("demo-proof")))
	case "execute":
		fmt.Printf("EXECUTE id=%s\n", *id)
	default:
		fmt.Println("Usage: -cmd send|verify|execute [-id ...] [-nonce ...] [-src ...] [-dst ...]")
	}
}
