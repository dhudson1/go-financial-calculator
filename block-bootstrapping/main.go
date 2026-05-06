package block_bootstrapping

type SampleInput[T interface{}] interface {
    sampleData T[]
    blockLength int64
    generatedLength int64
} 


type GeneratedOutput[T interface{}] struct {
    generatedData T[]
}


