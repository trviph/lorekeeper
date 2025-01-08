package lorekeeper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/trviph/collection"
)

type fileInfo struct {
	filePath string
	modtime  time.Time
}

func getArchives(pattern string) (*collection.List[string], error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get archived, caused by %w", err)
	}

	minHeap, err := collection.NewHeap[*fileInfo](func(current, other *fileInfo) bool {
		return current.modtime.Before(other.modtime)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get heap, caused by %w", err)
	}
	for _, match := range matches {
		info, err := getFileInfo(match)
		if err != nil {
			return nil, fmt.Errorf("failed to get file info %s, caused by %w", match, err)
		}
		minHeap.Push(info)
	}

	l := collection.NewList[string]()
	for !minHeap.IsEmpty() {
		min, err := minHeap.Pop()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info, caused by %w", err)
		}
		l.Append(min.filePath)
	}
	return l, nil
}

func getFileInfo(filePath string) (*fileInfo, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file, caused by %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed get file stat, caused by %w", err)
	}
	return &fileInfo{
		filePath: filePath,
		modtime:  stat.ModTime(),
	}, nil
}
