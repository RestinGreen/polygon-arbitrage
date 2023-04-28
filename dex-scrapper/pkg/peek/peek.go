package peek

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/memory"
)

type Peek struct {
	memory *memory.Memory
}

func NewPeek(memory *memory.Memory) *Peek {

	return &Peek{
		memory: memory,
	}
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
				case "l":
					fmt.Println("There are", len(p.memory.PairMemory.Pairs), "pairs loaded in the memory.")
				case "c":
					fmt.Println("Factories: ", len(p.memory.DexMemory.Dexs))
					for _, f := range p.memory.DexMemory.Dexs {
						fmt.Println("Factory", f.Factory, "has", *f.NumPairs, "pairs.")
					}
				case "q":
					exec.Command("stty", "-F", "/dev/tty", "echo").Run()

					// p.memory.CreationMutex.Lock()
					// for _, pair := range p.memory.PairMemory.Pairs {
					// 	p.db.UpdatePair(pair.PairAddress.Hex(), pair.Reserve0, pair.Reserve1, pair.LastUpdated)
					// }
					// p.memory.CreationMutex.Unlock()
					os.Exit(0)
				}
			default:
				// fmt.Println("Working..")
			}
			time.Sleep(time.Millisecond * 100)
		}

	}()
}
