package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// Test_folder_GetFoldersByOrgID tests the GetFoldersByOrgID method of the folder package.
func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	// Sample UUIDs for the tests
	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())
	orgIDInvalid := uuid.Must(uuid.NewV4())   // UUID with no folders in the list
	orgIDMixedCase := uuid.Must(uuid.NewV4()) // UUID for testing mixed-case folder names
	orgIDEmpty := uuid.Nil                    // Empty/invalid UUID

	// Sample folders for the tests
	folders := []*folder.Folder{
		{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
		{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
		{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
		{Name: "delta", OrgId: orgID2, Paths: "delta"},
		{Name: "echo", OrgId: orgID2, Paths: "delta.echo"},
		{Name: "Alpha", OrgId: orgIDMixedCase, Paths: "Alpha"}, // Case sensitivity test
		{Name: "beta", OrgId: orgIDMixedCase, Paths: "Alpha.beta"},
		{Name: "Gamma", OrgId: orgIDMixedCase, Paths: "Alpha.Gamma"},
	}

	// Define the test cases
	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []*folder.Folder
		want    []*folder.Folder
	}{
		{
			name:    "OrgID with multiple folders", // Test case for orgID with multiple folders
			orgID:   orgID1,
			folders: folders,
			want: []*folder.Folder{
				{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
			},
		},
		{
			name:    "OrgID with one folder", // Test case for orgID with one folder
			orgID:   orgID2,
			folders: folders,
			want: []*folder.Folder{
				{Name: "delta", OrgId: orgID2, Paths: "delta"},
				{Name: "echo", OrgId: orgID2, Paths: "delta.echo"},
			},
		},
		{
			name:    "OrgID with no folders", // Test case for orgID with no associated folders
			orgID:   orgIDInvalid,
			folders: folders,
			want:    []*folder.Folder{}, // Expecting empty list because no folders exist for this orgID
		},
		{
			name:    "Empty folder list", // Test case for empty folder list
			orgID:   orgID1,
			folders: []*folder.Folder{}, // No folders provided
			want:    []*folder.Folder{}, // Expecting empty list
		},
		{
			name:    "Case sensitivity in folder names", // Case sensitivity test for folder names
			orgID:   orgIDMixedCase,
			folders: folders,
			want: []*folder.Folder{
				{Name: "Alpha", OrgId: orgIDMixedCase, Paths: "Alpha"},
				{Name: "beta", OrgId: orgIDMixedCase, Paths: "Alpha.beta"},
				{Name: "Gamma", OrgId: orgIDMixedCase, Paths: "Alpha.Gamma"},
			},
		},
		{
			name:  "Multiple organizations with same folder names", // Test with same folder names under different orgs
			orgID: orgID1,
			folders: []*folder.Folder{
				{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
				{Name: "alpha", OrgId: orgID2, Paths: "alpha"},
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
			},
			want: []*folder.Folder{
				{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
			},
		},
		{
			name:    "Nil OrgID", // Test case for nil/empty orgID
			orgID:   orgIDEmpty,
			folders: folders,
			want:    []*folder.Folder{}, // Expecting empty list due to invalid orgID
		},
	}

	// Run through each test case
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Initialize the folder driver with the test folder data
			f := folder.NewDriver(tt.folders)

			// Call GetFoldersByOrgID with the provided orgID
			got := f.GetFoldersByOrgID(tt.orgID)

			// Assert that the result matches the expected output
			assert.Equal(t, tt.want, got)
		})
	}
}

// Test_folder_GetAllChildFolders tests the GetAllChildFolders method.
func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	// Sample UUIDs for the tests
	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())

	// Sample folders for the tests
	folders := []*folder.Folder{
		{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
		{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
		{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
		{Name: "delta", OrgId: orgID1, Paths: "alpha.delta"},
		{Name: "echo", OrgId: orgID1, Paths: "echo"},
		{Name: "foxtrot", OrgId: orgID2, Paths: "foxtrot"},
	}

	// Define the test cases
	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []*folder.Folder
		base    string
		want    []*folder.Folder
	}{
		{
			name:    "Multiple child folders", // Test with a folder that has multiple child folders
			orgID:   orgID1,
			folders: folders,
			base:    "alpha",
			want: []*folder.Folder{
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
				{Name: "delta", OrgId: orgID1, Paths: "alpha.delta"},
			},
		},
		{
			name:    "Single child folder", // Test with a folder that has one child folder
			orgID:   orgID1,
			folders: folders,
			base:    "bravo",
			want: []*folder.Folder{
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
			},
		},
		{
			name:    "No child folders", // Test with a folder that has no child folders
			orgID:   orgID1,
			folders: folders,
			base:    "charlie",
			want:    []*folder.Folder{}, // Expecting empty list since 'charlie' has no children
		},
		{
			name:    "Non-existent folder", // Test for a folder that doesn't exist
			orgID:   orgID1,
			folders: folders,
			base:    "nonexistent",
			want:    []*folder.Folder{}, // Expecting empty list since the folder doesn't exist
		},
		{
			name:    "No folders for orgID", // Test for an orgID with no associated folders
			orgID:   orgID2,
			folders: folders,
			base:    "alpha",
			want:    []*folder.Folder{}, // Expecting empty list since 'alpha' doesn't exist in orgID2
		},
		{
			name:    "Folder from a different orgID", // Test with a folder that belongs to a different orgID
			orgID:   orgID1,
			folders: folders,
			base:    "foxtrot",
			want:    []*folder.Folder{}, // Expecting empty list since 'foxtrot' is in orgID2, not orgID1
		},
		{
			name:    "Folder with no children in orgID2", // Test with a folder in a different orgID but no children
			orgID:   orgID2,
			folders: folders,
			base:    "foxtrot",
			want:    []*folder.Folder{}, // Expecting empty list since 'foxtrot' has no children
		},
		{
			name:    "Empty base folder name", // Test with an empty base folder name
			orgID:   orgID1,
			folders: folders,
			base:    "",
			want:    []*folder.Folder{}, // Expecting empty list since base folder name is empty
		},
		{
			name:    "Base folder is the root", // Test where the base folder is the root (no dot in the path)
			orgID:   orgID1,
			folders: []*folder.Folder{
				{Name: "root", OrgId: orgID1, Paths: "root"},
				{Name: "child1", OrgId: orgID1, Paths: "root.child1"},
				{Name: "child2", OrgId: orgID1, Paths: "root.child2"},
			},
			base: "root",
			want: []*folder.Folder{
				{Name: "child1", OrgId: orgID1, Paths: "root.child1"},
				{Name: "child2", OrgId: orgID1, Paths: "root.child2"},
			},
		},
		{
			name:    "Invalid orgID format", // Test with an invalid orgID (uuid.Nil)
			orgID:   uuid.Nil, // Invalid UUID for testing
			folders: folders,
			base:    "alpha",
			want:    []*folder.Folder{}, // Expecting empty list since the orgID is invalid
		},
		{
			name:    "Case-sensitive folder names", // Case sensitivity test for folder names
			orgID:   orgID1,
			folders: []*folder.Folder{
				{Name: "Alpha", OrgId: orgID1, Paths: "Alpha"},
				{Name: "Bravo", OrgId: orgID1, Paths: "Alpha.Bravo"},
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"}, // Lowercase base
			},
			base: "Alpha", // Test for case sensitivity
			want: []*folder.Folder{
				{Name: "Bravo", OrgId: orgID1, Paths: "Alpha.Bravo"},
			},
		},
		{
			name:    "Folders with special characters", // Test folder names with special characters
			orgID:   orgID1,
			folders: []*folder.Folder{
				{Name: "alpha@", OrgId: orgID1, Paths: "alpha@"},
				{Name: "bravo#", OrgId: orgID1, Paths: "alpha@.bravo#"},
				{Name: "charlie$", OrgId: orgID1, Paths: "alpha@.bravo#.charlie$"},
			},
			base: "alpha@",
			want: []*folder.Folder{
				{Name: "bravo#", OrgId: orgID1, Paths: "alpha@.bravo#"},
				{Name: "charlie$", OrgId: orgID1, Paths: "alpha@.bravo#.charlie$"},
			},
		},
		{
			name:    "Cyclic folder structure", // Test cyclic structure to ensure no infinite loops
			orgID:   orgID1,
			folders: []*folder.Folder{
				{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
				{Name: "alpha", OrgId: orgID1, Paths: "alpha.bravo.alpha"}, // Cyclic reference
			},
			base: "alpha",
			want: []*folder.Folder{
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
				{Name: "alpha", OrgId: orgID1, Paths: "alpha.bravo.alpha"}, // Should handle cyclic references
			},
		},
	}

	// Run through each test case
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Initialize the folder driver with the test folder data
			f := folder.NewDriver(tt.folders)

			// Call GetAllChildFolders with the provided orgID and base folder name
			got := f.GetAllChildFolders(tt.orgID, tt.base)

			// Assert that the result matches the expected output
			assert.Equal(t, tt.want, got)
		})
	}
}
