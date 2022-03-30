package prog

// func LoadProgram(pathFile string) (*Program, error) {
// 	fd, err := os.Open(pathFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer fd.Close()

// 	sizeInst := int64(binary.Size(inst.Inst{}))

// 	fi, err := fd.Stat()
// 	if err != nil {
// 		return nil, err
// 	}
// 	sizeFile := fi.Size()

// 	p := Program(make([]inst.Inst, sizeFile/sizeInst, sizeFile/sizeInst))
// 	err = binary.Read(fd, binary.BigEndian, &p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &p, nil
// }
