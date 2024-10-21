package folder

import (
	"log" // Import log package to log errors
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []*Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []*Folder {
	folders := f.folders

	res := []*Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res
}

// func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) []Folder {
// 	// Your code here...

// 	return []Folder{}
// }

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) []*Folder {
	// Retrieve all folders belonging to the given organization ID.
	folders := f.GetFoldersByOrgID(orgID)

	// Check if no folders exist for the orgID.
	if len(folders) == 0 {
		log.Printf("Error: No folders found for orgID '%s'", orgID)
		return []*Folder{}
	}

	// Find the base folder by the provided name.
	var baseFolder *Folder
	for _, folder := range folders {
		if folder.Name == name {
			baseFolder = folder
			break
		}
	}

	// If the base folder doesn't exist, log the error and return an empty list.
	if baseFolder == nil {
		log.Printf("Error: Folder '%s' does not exist in orgID '%s'", name, orgID)
		return []*Folder{}
	}

	// Prepare the list to store all child folders.
	var childFolders []*Folder

	// Base folder path to check for children.
	basePath := baseFolder.Paths + "."

	// Find all folders whose paths start with the base folder's path.
	for _, folder := range folders {
		// Ensure folder is a descendant (starts with basePath) but is not the base folder itself.
		if strings.HasPrefix(folder.Paths, basePath) {
			childFolders = append(childFolders, folder)
		}
	}

	// If no child folders are found, log that the folder has no children.
	if len(childFolders) == 0 {
		log.Printf("Info: Folder '%s' has no child folders in orgID '%s'", name, orgID)
		return []*Folder{}
	}

	return childFolders
}
