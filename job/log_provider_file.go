package job

import (
	"fmt"
	"os"
	"sync"
)

type (
	FileLogProvider struct {
		config      fileProviderConfig
		openFiles   map[string]*os.File
		openFilesMu *sync.Mutex
	}

	fileProviderConfig struct {
		LogProviderSharedConfig
		Directory string `yaml:"directory"`
	}
)

func init() {
	LogProviders["file"] = newFileLogProvider
}

func newFileLogProvider(config any) LogProvider {
	fConf := config.(fileProviderConfig)

	// TODO: parse directory and make parent dirs if they need to exist

	return &FileLogProvider{
		config: fConf,
	}
}

func (f *FileLogProvider) WriteLine(line string, meta LogMeta) error {
	// TODO: check if we should even write this log level, if not return nil

	fileDirectory := fmt.Sprintf("%s/worker=%s", f.config.Directory, meta.WorkerID)
	filePath := fmt.Sprintf("%s/task=%s.txt", fileDirectory, meta.TaskID)

	// Make the directory if we haven't
	err := os.MkdirAll(fileDirectory, 0777)
	if err != nil {
		return fmt.Errorf("error in os.MkdirAll for %s: %w", fileDirectory, err)
	}

	// Open the file if we haven't
	// There is definitely an optimization to lock per file when writing, but it's fast enough
	f.openFilesMu.Lock()
	defer f.openFilesMu.Unlock()

	openFile, exists := f.openFiles[filePath]
	if !exists {
		openFile, err = os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error in os.Create for %s: %w", filePath, err)
		}
	}

	_, err = openFile.WriteString(line)
	if err != nil {
		return fmt.Errorf("error in WriteString for %s: %w", filePath, err)
	}

	return nil
}

func (f *FileLogProvider) Close() error {
	f.openFilesMu.Lock()
	defer f.openFilesMu.Unlock()
	for _, filePath := range f.openFiles {
		err := filePath.Close()
		if err != nil {
			return fmt.Errorf("error in file.Close for %s: %w", filePath, err)
		}
	}

	return nil
}
