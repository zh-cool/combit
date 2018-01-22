package main

const (
	BIT_SIZE    = (1 << 32)
	MATRIX_SIZE = 16
)

type Matric struct {
	matric [MATRIX_SZIE * MATRIX_SIZE >> 3]byte
	flag   uint8
}

func CreateMatric(data []byte) []Matric {
	bits := len(data) << 3
	CMatric := len(data) >> 5

	matric := make([]Matric, CMatric)

	for i := 0; i < CMatric; i++ {
		for j := 0; i < 32; j++ {
			matric[i].matric[j] |= data[j+i<<5]
		}
	}
}

func (matric *Matric) prepareData() []byte {
	size := len(Matric) >> 3
	data := make([]byte, size)
	for i, v := range Matric {
		if v.flag > 0 {
			data[i/8] |= 1 << (i % 8)
		}
	}

	return data
}
