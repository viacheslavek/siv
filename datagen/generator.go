package datagen

import (
	"fmt"
	"math/rand"
	"time"
)

// Batch представляет сущность batch.
type Batch struct {
	Time   time.Time
	Length int
	Data   []int
}

// GeneratorMode определяет режим генерации данных.
type GeneratorMode int

const (
	RandomMode GeneratorMode = iota
	AmplitudeK1Mode
	AmplitudeK2Mode
)

// DataGenerator генерирует данные на основе заданной схемы и режима.
func DataGenerator(t time.Duration, l, N, K1, K2 int, mode GeneratorMode) <-chan Batch {
	dataStream := make(chan Batch)

	go func() {
		defer close(dataStream)
		rand.Seed(time.Now().UnixNano())

		for {
			time.Sleep(t)
			batch := Batch{
				Time:   time.Now(),
				Length: l,
			}

			for i := 0; i < l; i++ {
				var value int

				switch mode {
				case RandomMode:
					value = rand.Intn(N + 1)
				case AmplitudeK1Mode:
					value = rand.Intn(K1 + 1)
				case AmplitudeK2Mode:
					value = rand.Intn(K2 + 1)
				}

				batch.Data = append(batch.Data, value)
			}

			dataStream <- batch
		}
	}()

	return dataStream
}

func example() {
	// Пример использования с переключением режима
	generator := DataGenerator(1*time.Second, 5, 10, 5, 15, RandomMode)

	for i := 0; i < 3; i++ {
		batch := <-generator
		fmt.Printf("Time: %v, Data: %v\n", batch.Time, batch.Data)
	}

	// Изменение режима генерации
	generator = DataGenerator(1*time.Second, 5, 10, 5, 15, AmplitudeK1Mode)

	for i := 0; i < 3; i++ {
		batch := <-generator
		fmt.Printf("Time: %v, Data: %v\n", batch.Time, batch.Data)
	}
}
