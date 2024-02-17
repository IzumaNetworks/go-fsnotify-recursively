package fsnotifyr

type Folders []*folder

func (folders Folders) Len() int {
	return len(folders)
}

func (folders Folders) Swap(i, j int) {
	folders[i], folders[j] = folders[j], folders[i]
}

func (folders Folders) Less(i, j int) bool {
	return folders[i].path < folders[j].path
}
