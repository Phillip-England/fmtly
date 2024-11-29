package cli

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyOptionWatch = "--watch"
)

// ##==================================================================
func getOptionList() []string {
	return []string{KeyOptionWatch}
}

// ##==================================================================
type Option interface {
	Print()
	GetType() string
	Inject(ex Executor, process func() error) func() error
}

func NewOption(arg string) (Option, error) {
	match := purse.FindMatchInStrSlice(getOptionList(), arg)
	if match == "" {
		return nil, fmt.Errorf("invalid option selected: %s\nRun 'gtml help' for usage.", arg)
	}
	switch match {
	case KeyOptionWatch:
		opt, err := NewOptionWatch()
		if err != nil {
			return nil, err
		}
		return opt, err
	}
	return nil, fmt.Errorf("invalid option selected: %s\nRun 'gtml help' for usage.", arg)
}

// ##==================================================================
type OptionWatch struct {
	Type string
}

func NewOptionWatch() (*OptionWatch, error) {
	opt := &OptionWatch{
		Type: KeyOptionWatch,
	}
	return opt, nil
}

func (opt *OptionWatch) GetType() string { return opt.Type }
func (opt *OptionWatch) Print()          { fmt.Println(opt.Type) }
func (opt *OptionWatch) Inject(ex Executor, process func() error) func() error {
	return func() error {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		dirToWatch := ex.GetCommand().GetFilteredArgs()[0]
		if err := watcher.Add(dirToWatch); err != nil {
			return err
		}

		err = process() // Initial run
		if err != nil {
			return err
		}

		var debounceTimer *time.Timer
		debounceDuration := 100 * time.Millisecond

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					if event.Op&fsnotify.Write == fsnotify.Write {
						fmt.Printf("File modified: %s\n", event.Name)

						if debounceTimer != nil {
							debounceTimer.Stop()
						}
						debounceTimer = time.AfterFunc(debounceDuration, func() {
							err := process()
							if err != nil {
								fmt.Printf("Error running process: %v\n", err)
							}
						})
					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					fmt.Printf("Watcher error: %v\n", err)
				}
			}
		}()

		fmt.Printf("Watching directory: %s\n", dirToWatch)
		select {} // Block forever.
	}
}
