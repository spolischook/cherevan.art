package artWork

func (sl ArtWorks) ForEach(f func(ArtWork)) {
	for i := range sl {
		f(sl[i])
	}
}
