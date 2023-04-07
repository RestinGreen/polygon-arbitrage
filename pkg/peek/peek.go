package peek

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/RestinGreen/polygon-arbitrage/pkg/memory"
)

type Peek struct {
	memory *memory.Memory
}

func NewPeek() *Peek {

	return &Peek{}
}

func (p *Peek) SetDexMemory(memory *memory.Memory) {

	p.memory = memory
}

func (p *Peek) StartPeek() {

	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	go func() {

		for {
			select {
			case stdin, _ := <-ch:
				switch stdin {
				// case "l":
				// 	fmt.Println("-----------------------------------------------------")
				// 	for router, v := range p.dexMemory.DexMemory {
				// 		fmt.Println(router)
				// 		fmt.Println("\t", v.Factory)
				// 		fmt.Println("\t pairs")
				// 		for _, pair := range v.Pairs {
				// 			fmt.Println("\t\t------------------------------------------")
				// 			fmt.Println("\t\t ", pair.PairAddress)
				// 			fmt.Println("\t\t ", pair.Token0Address)
				// 			fmt.Println("\t\t ", pair.Token1Address)
				// 			fmt.Println("\t\t ", pair.Reserve0)
				// 			fmt.Println("\t\t ", pair.Reserve1)
				// 			fmt.Println("\t\t ", pair.LastUpdated)
				// 			fmt.Println("\t\t------------------------------------------")
				// 		}
				// 	}
				// 	fmt.Println("-----------------------------------------------------")
				case "c":
					fmt.Println("Factories: ", len(p.memory.DexMemory.Dexs))
					for _, f := range p.memory.DexMemory.Dexs {
						fmt.Println("Factory", f.Factory, "has", *f.NumPairs, "pairs.")
					}
				case "q":
					exec.Command("stty", "-F", "/dev/tty", "echo").Run()
					os.Exit(0)
				}
			default:
				// fmt.Println("Working..")
			}
			time.Sleep(time.Millisecond * 100)
		}

	}()
}
