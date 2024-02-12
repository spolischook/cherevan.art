package artWork

func (sl ArtWorks) Filter(f func(ArtWork) bool) (res ArtWorks) {
	for i := range sl {
		if f(sl[i]) {
			res = append(res, sl[i])
		}
	}
	return
}
