package artWork

import "github.com/cherevan.art/src/tool"

func (sl ArtWorks) MaxId() int {
	maxId := 0
	sl.ForEach(func(w ArtWork) { maxId = tool.Max(maxId, w.ID) })
	return maxId
}
