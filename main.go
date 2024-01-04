package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	nomad_api "github.com/hashicorp/nomad/api"
)

const DEFAULT_DIR_PERMS = 0o777

var (
	invalidPathChars = regexp.MustCompile("[^a-zA-Z0-9-_~/]")

	// version of nomad-var-dirsync being run, set with ldflags
	version = "dev"
)

func writeDir(client *nomad_api.Client, root string, sourceDir string) error {
	err := filepath.WalkDir(sourceDir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not walk to %s: %w", path, err)
		}

		if dir.IsDir() {
			return nil
		}

		fileInfo, err := dir.Info()
		if err != nil {
			return fmt.Errorf("failed getting info for %s: %w", path, err)
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed reading file %s: %w", path, err)
		}

		sanitizedPath := invalidPathChars.ReplaceAllString(path, "_")
		sanitizedPath = filepath.Join(root, sanitizedPath)

		newVar := nomad_api.Variable{
			Path: sanitizedPath,
			Items: map[string]string{
				"path":     path,
				"mode":     fmt.Sprintf("%o", fileInfo.Mode()),
				"contents": string(contents),
			},
		}

		if _, _, err := client.Variables().Create(&newVar, nil); err != nil {
			return fmt.Errorf("failed creating var %s for file %s: %w", sanitizedPath, path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Error walking dir %s: %w", sourceDir, err)
	}

	return nil
}

func readDir(client *nomad_api.Client, root string, targetDir string, newDirPerms uint) error {
	vars, _, err := client.Variables().List(&nomad_api.QueryOptions{
		Prefix: root,
	})
	if err != nil {
		return fmt.Errorf("failed reading vars from root %s: %w", root, err)
	}

	for _, varInfo := range vars {
		log.Printf("Reading variable %s", varInfo.Path)

		fileVar, _, err := client.Variables().Read(varInfo.Path, &nomad_api.QueryOptions{})
		if err != nil {
			log.Printf("Failed reading variable %s: %v", varInfo.Path, err)
		}

		filePath := filepath.Join(targetDir, fileVar.Items["path"])
		fileModeString := fileVar.Items["mode"]
		fileContents := fileVar.Items["contents"]

		fileMode, err := strconv.ParseUint(fileModeString, 8, 32)
		if err != nil {
			return fmt.Errorf("Failed parsing file mode for %s. %s: %w", filePath, fileModeString, err)
		}

		parentDir := filepath.Dir(filePath)
		if _, err := os.Stat(parentDir); err != nil {
			if err = os.MkdirAll(parentDir, fs.FileMode(newDirPerms)); err != nil {
				return fmt.Errorf("error creating paretn dir for file at path %s: %w", filePath, err)
			}
		}

		err = os.WriteFile(filePath, []byte(fileContents), os.FileMode(fileMode))
		if err != nil {
			return fmt.Errorf("Failed writing file %s: %w", filePath, err)
		}
	}

	return nil
}

func main() {
	root := flag.String("root-var", "", "root path for nomad variable")
	showVersion := flag.Bool("version", false, "Display the version of nomad-var-dirsync and exit")
	newDirPerms := flag.Uint("dir-perms", DEFAULT_DIR_PERMS, "default permissions for new directories (default: 0o777)")
	flag.Parse()

	// Print version if flag is provided
	if *showVersion {
		fmt.Println("nomad-var-dirsync version:", version)

		return
	}

	action := flag.Arg(0)
	target := flag.Arg(1)

	if *root == "" {
		log.Fatal("Must provide a nomad variable root -root-var")
	}

	targetStat, err := os.Stat(target)
	if err != nil {
		log.Fatalf("Failed reading target file `%s`. %v", target, err)
	}

	if !targetStat.IsDir() {
		log.Fatalf("must provide a path to a directory: %s", target)
	}

	client, err := nomad_api.NewClient(&nomad_api.Config{
		SecretID: os.Getenv("NOMAD_TOKEN"),
	})
	if err != nil {
		log.Fatalf("failed creating nomad client: %v", err)
	}

	switch action {
	case "write":
		if err = writeDir(client, *root, target); err != nil {
			log.Fatalf("Failed writing directory: %v", err)
		}
	case "read":
		if err = readDir(client, *root, target, *newDirPerms); err != nil {
			log.Fatalf("Failed reading to files for path %v", err)
		}
	default:
		log.Fatalf("Expected action read or write, found %s", action)
	}
}
