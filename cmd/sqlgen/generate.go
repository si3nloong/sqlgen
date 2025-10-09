package main

import (
	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen"
	"github.com/spf13/cobra"
)

var (
	genOpts struct {
		watch bool
		force bool
	}

	genCmd = &cobra.Command{
		Use:   "generate [source]",
		Short: "Generate boilerplate code based on go struct",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				cfg = codegen.DefaultConfig()
			)

			// // Create new watcher.
			// watcher, err := fsnotify.NewWatcher()
			// if err != nil {
			// 	return err
			// }
			// defer watcher.Close()

			// // Start listening for events.
			// go func() {
			// 	for {
			// 		select {
			// 		case event, ok := <-watcher.Events:
			// 			if !ok {
			// 				return
			// 			}
			// 			log.Println("event:", event)
			// 			if event.Has(fsnotify.Write) {
			// 				log.Println("modified file:", event.Name)
			// 			}
			// 		case err, ok := <-watcher.Errors:
			// 			if !ok {
			// 				return
			// 			}
			// 			log.Println("error:", err)
			// 		}
			// 	}
			// }()

			// // Add a path.
			// if err := watcher.Add(filepath.Join(fileutil.Getpwd(), args[0])); err != nil {
			// 	return err
			// }

			// <-make(chan struct{})

			// If user pass the source, then we refer to it.
			cfg.Source = []string{args[0]}

			return codegen.Generate(cfg)
		},
	}
)

func init() {
	genCmd.Flags().BoolVarP(&genOpts.watch, "watch", "w", false, "watch the file changes and re-generate.")
	genCmd.Flags().BoolVarP(&genOpts.force, "force", "f", false, "force to execute")
}
