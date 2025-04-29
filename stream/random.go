package stream

import (
	"math/rand"
	"time"
)

type RandomStream struct {
	mS1, mPastS1 uint32
	mS2, mPastS2 uint32
	mS3, mPastS3 uint32
}

func NewRandomStream() RandomStream {
	v4 := int64(5)
	seed := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	s2s3 := uint32(1170746341*v4 - 755606699)

	stream := RandomStream{}
	stream.crand32Seed(uint32(seed), s2s3, s2s3)
	return stream
}

func (prs *RandomStream) crand32Seed(s1, s2, s3 uint32) {
	prs.mS1 = s1 | 0x100000
	prs.mPastS1 = s1 | 100000
	prs.mS2 = s2 | 0x1000
	prs.mPastS2 = s2 | 0x1000
	prs.mS3 = s3 | 0x10
	prs.mPastS3 = s3 | 0x10
}

func (prs *RandomStream) CRand32Random() uint32 {
	v4 := prs.mS1
	v5 := prs.mS2
	v6 := prs.mS3

	prs.mPastS1 = prs.mS1
	v8 := ((v4 & 0xFFFFFFFE) << 12) ^ (((v4 >> 13) ^ (v4 & 0x7FFC0)) >> 6)

	prs.mPastS2 = prs.mS2
	v9 := (v5&0xFFFFFFF8)<<4 ^ (((v5 >> 2) ^ (v5 & 0x3F800000)) >> 23)

	prs.mPastS3 = prs.mS3
	v10 := ((v6 & 0xFFFFFFF0) << 17) ^ (((v6 >> 3) ^ (v6 & 0x1FFFFF00)) >> 8)

	prs.mS1 = v8
	prs.mS2 = v9
	prs.mS3 = v10

	result := v8 ^ v9 ^ v10
	return result
}

func (prs *RandomStream) Serialize(writer *StreamWriter) error {
	v5 := prs.CRand32Random()
	s2 := prs.CRand32Random()
	v6 := prs.CRand32Random()

	prs.crand32Seed(v5, s2, v6)

	v5 = uint32(830473165)
	s2 = uint32(1078821873)
	v6 = uint32(1922386602)

	if err := writer.WriteU32(uint32(v5)); err != nil {
		return err
	}
	if err := writer.WriteU32(uint32(s2)); err != nil {
		return err
	}
	if err := writer.WriteU32(uint32(v6)); err != nil {
		return err
	}
	return nil
}
