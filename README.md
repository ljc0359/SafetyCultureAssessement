# Code Explaination
## Set up
* Required dependency is added to `go.mod`.
* Please run `go mod tidy` to download all the required dependency before running the test.

## Testing
### GetFoldersByOrgID
Here’s a summary of all the scenarios being tested in the `Test_folder_GetFoldersByOrgID` function:
* Multiple Folders for a Single OrgId.
* One Folder for a Single OrgId.
* No Folders for an OrgId.
* Empty Folder List.
* Case Sensitivity for Folder Names.
* Handling Folders with the Same Name Across Different Organizations.
* Invalid inputs like `nil` OrgId.

### GetAllChildFolders
Here’s a summary of all the scenarios being tested in the `Test_folder_GetAllChildFolders` function:

* Multiple child folders
* Single child folder
* No child folders
* Non-existent folder
* No folders for orgID
* Folder from a different orgID
* Folder with no children in orgID2
* Empty base folder name
* Base folder is the root
* Invalid orgID format
* Case-sensitive folder names
* Folders with special characters
* Cyclic folder structure

### MoveFolder
Here’s a summary of all the scenarios being tested in the `Test_folder_MoveFolder` function:

#### Negative Cases
* Moving a Folder to Itself
* Moving a Folder Under Its Own Child (Creates Cyclic Reference)
* Source Folder Does Not Exist
* Destination Folder Does Not Exist
* Moving a Folder to a Different Organization
* Moving a Root Folder Under Itself


#### Positive Cases
* Moving a Folder Under a New Parent within the Same Organization
* Moving a Leaf Folder Under Another Folder
* Moving a Folder with Children Under a New Parent
* Moving a Child Folder Under Another Sibling Folder
* Moving a Leaf Folder Under a Non-Root Folder
* Moving a Folder to Become a Sibling of Its Former Parent

#### Multiple operations
Additionally, `Test_folder_MoveFolder_MultipleOperations` valids the folder structure after every operations when multiple MoveFolder operation is used.