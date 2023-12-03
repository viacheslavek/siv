package main

import "github.com/VyacheslavIsWorkingNow/siv/ssh_testing/server/internal"

// ssh user@127.0.0.1 -p 2222

// TODO: для сервера в internal нужно рассортировать все по папкам

func main() {
	internal.RunServer()
}
