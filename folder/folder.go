package folder

import (
	"fmt"
	"github.com/gofrs/uuid"
)

type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []*Folder
	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) []*Folder

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]*Folder, error)
}

type driver struct {
	// define attributes here
	// data structure to store folders
	// or preprocessed data

	// example: feel free to change the data structure, if slice is not what you want
	folders []*Folder
}

func NewDriver(folders []*Folder) *driver {
	return &driver{
		// initialize attributes here
		folders: folders,
	}
}

// PrintFolders recursively prints the folder tree structure
func PrintFolders(folders []*Folder) {
    fmt.Printf("\n\n\n*********************************************************************************************\n")
    for _, folder := range folders {
        printFolder(folder, 0)
        fmt.Printf("\n\n")
    }
    fmt.Printf("*********************************************************************************************\n\n\n")
}

// Helper function to print a folder and its children
func printFolder(folder *Folder, level int) {
    indent := ""
    for i := 0; i < level; i++ {
        indent += "  "
    }

    fmt.Printf("%sFolder Name: %s, Path: %s, OrgId: %s, Address: %p\n", indent, folder.Name, folder.Paths, folder.OrgId, folder)

    // Print parent information if available
    if folder.Parent != nil {
        fmt.Printf("%s  Parent: %s\n", indent, folder.Parent.Name)
    } else {
        fmt.Printf("%s  Parent: None (root folder)\n", indent)
    }

    // Recursively print child folders
    if len(folder.Children) > 0 {
        fmt.Printf("%s  Children:\n", indent)
        for _, child := range folder.Children {
            printFolder(child, level+1)
        }
    } else {
        fmt.Printf("%s  Children: None\n", indent)
    }
}
