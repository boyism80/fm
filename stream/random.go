package stream

type RandomStream struct {
	mS1, mPastS1 int64
	mS2, mPastS2 int64
	mS3, mPastS3 int64
}

func NewRandomStream() RandomStream {
	v4 := int64(5)

	s2s3 := int64(uint32(1170746341*v4) - 755606699)

	stream := RandomStream{}
	stream.crand32Seed(19930809, s2s3, s2s3)
	return stream
}

func (prs *RandomStream) crand32Seed(s1, s2, s3 int64) {
	prs.mS1 = int64(s1 | 0x100000)
	prs.mPastS1 = int64(s1 | 100000)
	prs.mS2 = int64(s2 | 0x1000)
	prs.mPastS2 = int64(s2 | 0x1000)
	prs.mS3 = int64(s3 | 0x10)
	prs.mPastS3 = int64(s3 | 0x10)
}

func (prs *RandomStream) CRand32Random() int64 {
	v4 := prs.mS1
	v5 := prs.mS2
	v6 := prs.mS3
	v7 := prs.mS1

	prs.mPastS1 = v4
	v8 := ((uint64(v4) & 0xFFFFFFFFFFFFFFFE) << 12) ^ uint64((v7&0x7FFC0^(v4>>13))>>6)

	prs.mPastS2 = v5
	v9 := ((uint64(v5) & 0xFFFFFFFFFFFFFFF8) << 4) ^ uint64(((v5>>2)^v5&0x3F800000)>>23)

	prs.mPastS3 = v6
	v10 := ((uint64(v6) & 0xFFFFFFFFFFFFFFF0) << 17) ^ uint64(((v6>>3)^v6&0x1FFFFF00)>>8)

	prs.mS1 = int64(v8)
	prs.mS2 = int64(v9)
	prs.mS3 = int64(v10)

	return int64(uint32(int64(v8) ^ int64(v9) ^ int64(v10)))
}

func (prs *RandomStream) Serialize(writer *StreamWriter) error {
	v5 := prs.CRand32Random()
	s2 := prs.CRand32Random()
	v6 := prs.CRand32Random()

	prs.crand32Seed(v5, s2, v6)

	writer.Write32(int32(v5))
	writer.Write32(int32(s2))
	writer.Write32(int32(v6))
	return nil
}
