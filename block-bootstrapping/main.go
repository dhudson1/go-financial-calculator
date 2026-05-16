package block_bootstrapping

import (
    "math/rand/v2"
)

type SampleInput[T any] struct {
    SampleData []T
    MinBlockLength uint64
    BlockLengthVarience uint64
    GeneratedLength uint64
    GeneratedSamples uint64
    Seed [2]uint64
} 


func generateBootstrapSample[S ~[]T, T any](input *SampleInput[T], subSeeds [2]uint64, outputChannel chan S) error {
    r := rand.New(rand.NewPCG(subSeeds[0], subSeeds[1]))
    retVal := make(S, input.GeneratedLength)
    leftInBlock := 0
    sampleIdx := 0
    for currentIdx := range input.GeneratedLength {
        if leftInBlock <= 0 {
            if input.BlockLengthVarience > 0 {
                leftInBlock = r.UInt64N(input.BlockLengthVarience) + input.MinBlockLength
            } else {
                leftInBlock = input.MinBlockLength
            }
            sampleIdx = r.UInt64N(len(input.SampleData) - leftInBlock)
        }
        retVal[currentIdx] = input.SampleData[sampleIdx]
        sampleIdx += 1
        leftInBlock -= 1
    }
    outputChannel <- retVal
}

func GenerateBootstrap[S ~[]T, T any](input *SampleInput[T], outputChannel chan S) error {
    if input.MinBlockLength == 0 {
        return errors.New("MinBlockLength must be greater than 0")
    }

    if len(input.SampleData) < (input.MinBlockLength + input.BlockLengthVarience) {
        return errors.New("SampleData's length must be greater than the sum of the MinBlockLength and BlockLengthVarience.")
    }

    r := rand.New(rand.NewPCG(input.Seed[0], input.Seed[1]))
    var wg sync.WaitGroup
    for currentIdx := range input.GeneratedSamples {
        wg.Add(1)
        go generateBootstrapSample[S, T](input, [2]uint64{r.Uint64(), r.Uint64()}, outputChannel)
    }

    wg.Wait()
    return nil
}

