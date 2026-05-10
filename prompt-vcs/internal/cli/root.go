package cli

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"prompt-vcs/internal/api"
	"prompt-vcs/internal/models"
	storage "prompt-vcs/internal/storage/db"
	"prompt-vcs/internal/version"

	"github.com/spf13/cobra"
)

var (
	dbPath      string
	promptTags  []string
	collection  string
	showVersion bool
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt-vcs",
		Short: "A version control system for LLM prompts",
		Long: `Prompt VCS is a command-line tool and web application for managing,
versioning, and organizing LLM prompts with Git-like semantics.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if dbPath == "" {
				dbPath = os.Getenv("HOME") + "/.prompt-vcs/prompts.db"
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Println("prompt-vcs version", version.Version)
				return
			}
			cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(&dbPath, "db", "d", "", "Database path (default: ~/.prompt-vcs/prompts.db)")
	cmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Print the version number")

	cmd.AddCommand(initCmd())
	cmd.AddCommand(addCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(showCmd())
	cmd.AddCommand(diffCmd())
	cmd.AddCommand(commitCmd())
	cmd.AddCommand(rollbackCmd())
	cmd.AddCommand(searchCmd())
	cmd.AddCommand(serveCmd())

	return cmd
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new prompt repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			fmt.Println("Initialized empty prompt repository in", dbPath)
			return nil
		},
	}
}

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name] [content]",
		Short: "Add a new prompt to the repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()

			prompt := &models.Prompt{
				ID:         generateID(),
				Name:       args[0],
				Content:    args[1],
				Tags:       promptTags,
				Collection: collection,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := db.CreatePrompt(prompt); err != nil {
				return err
			}

			fmt.Printf("Created prompt '%s' (%s)\n", prompt.Name, prompt.ID)
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&promptTags, "tags", "t", []string{}, "Comma-separated tags")
	cmd.Flags().StringVarP(&collection, "collection", "c", "default", "Collection name")

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [collection]",
		Short: "List all prompts in a collection",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()

			coll := "default"
			if len(args) > 0 {
				coll = args[0]
			}

			prompts, err := db.ListPrompts(coll)
			if err != nil {
				return err
			}

			if len(prompts) == 0 {
				fmt.Printf("No prompts found in collection '%s'\n", coll)
				return nil
			}

			fmt.Printf("Prompts in collection '%s':\n", coll)
			fmt.Println("ID\t\t\t\t\tName\t\tVersion\tUpdated")
			for _, p := range prompts {
				fmt.Printf("%s\t%s\t%d\t%s\n", p.ID, p.Name, p.Version, p.UpdatedAt.Format("2006-01-02"))
			}
			return nil
		},
	}
}

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [id]",
		Short: "Show prompt details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()

			prompt, err := db.GetPrompt(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Prompt: %s\n", prompt.Name)
			fmt.Printf("ID: %s\n", prompt.ID)
			fmt.Printf("Version: %d\n", prompt.Version)
			fmt.Printf("Collection: %s\n", prompt.Collection)
			fmt.Printf("Tags: %v\n", prompt.Tags)
			fmt.Printf("Content:\n%s\n", prompt.Content)
			return nil
		},
	}
}

func diffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff [id] [version1] [version2]",
		Short: "Show differences between prompt versions",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Diff feature coming soon!")
			return nil
		},
	}
}

func commitCmd() *cobra.Command {
	var commitContent string

	cmd := &cobra.Command{
		Use:   "commit [id] [message]",
		Short: "Commit a new version of a prompt",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if commitContent == "" {
				return fmt.Errorf("content is required (use --content flag)")
			}

			db, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()

			if err := db.UpdatePromptVersion(args[0], commitContent); err != nil {
				return err
			}

			hash := generateID()
			if err := db.AddCommit(args[0], hash, args[1], "user"); err != nil {
				return err
			}

			fmt.Printf("Committed prompt %s: %s (hash: %s)\n", args[0], args[1], hash)
			return nil
		},
	}

	cmd.Flags().StringVarP(&commitContent, "content", "m", "", "New content for the prompt")

	return cmd
}

func rollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback [id] [version]",
		Short: "Rollback a prompt to a specific version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Rollback feature coming soon!")
			return nil
		},
	}
}

func searchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Search prompts by content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Searching for: %s\n", args[0])
			fmt.Println("Search feature coming soon!")
			return nil
		},
	}
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the web API server",
		Long:  "Start the web API server with RESTful endpoints for managing prompts",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := api.NewServer(dbPath)
			if err != nil {
				return err
			}
			fmt.Println("Web API server starting on :8080")
			return server.Start(":8080")
		},
	}
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func Execute() error {
	return rootCmd().Execute()
}
