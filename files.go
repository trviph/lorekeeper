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
	size     int
	modtime  time.Time
}

func getArchives(pattern string) (*collection.List[*fileInfo], int, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get archived, caused by %w", err)
	}

	minHeap, err := collection.NewHeap(func(current, other *fileInfo) bool {
		return current.modtime.Before(other.modtime)
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get heap, caused by %w", err)
	}
	for _, match := range matches {
		info, err := getFileInfo(match)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get file info %s, caused by %w", match, err)
		}
		minHeap.Push(info)
	}

	l := collection.NewList[*fileInfo]()
	totalSize := 0
	for !minHeap.IsEmpty() {
		min, err := minHeap.Pop()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get file info, caused by %w", err)
		}
		l.Append(min)
		totalSize += min.size
	}
	return l, totalSize, nil
}

func getFileInfo(filePath string) (*fileInfo, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed get file stat, caused by %w", err)
	}
	return &fileInfo{
		filePath: filePath,
		modtime:  stat.ModTime(),
		size:     int(stat.Size()),
	}, nil
}
