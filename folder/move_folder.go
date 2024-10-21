package folder

import (
	"fmt"
	"log"
)

// isDescendant checks if dest is a descendant of source
func isDescendant(source, dest *Folder) bool {
	current := dest
	for current != nil {
		if current == source {
			return true
		}
		current = current.Parent
	}
	return false
}

// removeChild removes a child folder from the parent's children slice
func removeChild(parent, childToRemove *Folder) {
	if parent == nil {
		return
	}
	for i, child := range parent.Children {
		if child == childToRemove {
			// Remove the child at index i
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			break
		}
	}
}

// updatePaths updates the Paths of the folder and its descendants
func updatePaths(folder *Folder, parentPath string) {
	folder.Paths = parentPath + "." + folder.Name
	for _, child := range folder.Children {
		updatePaths(child, folder.Paths)
	}
}

// MoveFolder moves a folder and its subtree to a new destination folder
func (f *driver) MoveFolder(name string, dst string) ([]*Folder, error) {
	// Build a map of folder names to folder pointers for quick lookup
	folderMap := make(map[string]*Folder)
	for _, folder := range f.folders {
		folderMap[folder.Name] = folder
	}

	// Get the source and destination folders
	sourceFolder, sourceExists := folderMap[name]
	destFolder, destExists := folderMap[dst]

	// Error handling
	if !sourceExists {
		log.Printf("Error: Source folder '%s' does not exist", name)
		return nil, fmt.Errorf("source folder '%s' does not exist", name)
	}
	if !destExists {
		log.Printf("Error: Destination folder '%s' does not exist", dst)
		return nil, fmt.Errorf("destination folder '%s' does not exist", dst)
	}

	// Error handling for moving to itself
	if sourceFolder == destFolder {
		log.Printf("Error: Cannot move folder '%s' to itself", name)
		return nil, fmt.Errorf("cannot move folder '%s' to itself", name)
	}

	// Error handling for moving to a child of itself
	if isDescendant(sourceFolder, destFolder) {
		log.Printf("Error: Cannot move folder '%s' to a child of itself", name)
		return nil, fmt.Errorf("cannot move folder '%s' to a child of itself", name)
	}

	// Error handling for moving to a different organization
	if sourceFolder.OrgId != destFolder.OrgId {
		log.Printf("Error: Cannot move folder '%s' to a different organization", name)
		return nil, fmt.Errorf("cannot move folder '%s' to a different organization", name)
	}

	// Remove the source folder from its current parent's children
	if sourceFolder.Parent != nil {
		removeChild(sourceFolder.Parent, sourceFolder)
	}
	// fmt.Printf("Remove the source folder from its current parent's children\n")
	// PrintFolders(f.folders)

	// Update the parent of the source folder to be the destination folder
	sourceFolder.Parent = destFolder
	destFolder.Children = append(destFolder.Children, sourceFolder)

	// fmt.Printf("Update the parent of the source folder to be the destination folder\n")
	// PrintFolders(f.folders)

	// Update the paths of the source folder and its descendants
	updatePaths(sourceFolder, destFolder.Paths)
	// fmt.Printf("Update the paths of the source folder and its descendants\n")
	// PrintFolders(f.folders)

	// Return the updated folder structure
	return f.folders, nil
}
