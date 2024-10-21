package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// Helper function to initialize folders and create a folder map
func initializeFolders(orgID1, orgID2 uuid.UUID) ([]*folder.Folder, map[string]*folder.Folder) {
	alpha := &folder.Folder{Name: "alpha", OrgId: orgID1, Paths: "alpha"}
	bravo := &folder.Folder{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo", Parent: alpha}
	charlie := &folder.Folder{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie", Parent: bravo}
	delta := &folder.Folder{Name: "delta", OrgId: orgID1, Paths: "alpha.delta", Parent: alpha}
	echo := &folder.Folder{Name: "echo", OrgId: orgID1, Paths: "alpha.delta.echo", Parent: delta}
	foxtrot := &folder.Folder{Name: "foxtrot", OrgId: orgID2, Paths: "foxtrot"}
	golf := &folder.Folder{Name: "golf", OrgId: orgID1, Paths: "golf"}

	// Attach children to their respective parents
	alpha.Children = []*folder.Folder{bravo, delta}
	bravo.Children = []*folder.Folder{charlie}
	delta.Children = []*folder.Folder{echo}

	// Create a slice of folder pointers
	folders := []*folder.Folder{alpha, bravo, charlie, delta, echo, foxtrot, golf}

	// Create a map of folder names to folder pointers for easy access
	folderMap := map[string]*folder.Folder{
		"alpha":   alpha,
		"bravo":   bravo,
		"charlie": charlie,
		"delta":   delta,
		"echo":    echo,
		"foxtrot": foxtrot,
		"golf":    golf,
	}

	return folders, folderMap
}

// test functions with single operation, including positive and negative cases
func Test_folder_MoveFolder(t *testing.T) {
	// Sample UUIDs for testing
	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())

	// Define the test cases
	tests := []struct {
		name          string
		source        string
		dest          string
		expectError   bool
		expectedError string
		validateFunc  func(*testing.T, map[string]*folder.Folder)
	}{
		// Invalid move - moving to a child of itself
		{
			name:          "Invalid move - bravo to charlie (moving to a child of itself)",
			source:        "bravo",
			dest:          "charlie",
			expectError:   true,
			expectedError: "cannot move folder 'bravo' to a child of itself",
		},
		// Invalid move - moving to itself
		{
			name:          "Invalid move - bravo to itself",
			source:        "bravo",
			dest:          "bravo",
			expectError:   true,
			expectedError: "cannot move folder 'bravo' to itself",
		},
		// Invalid move - source folder doesn't exist
		{
			name:          "Invalid move - source folder doesn't exist",
			source:        "invalid_folder",
			dest:          "delta",
			expectError:   true,
			expectedError: "source folder 'invalid_folder' does not exist",
		},
		// Invalid move - destination folder doesn't exist
		{
			name:          "Invalid move - destination folder doesn't exist",
			source:        "bravo",
			dest:          "invalid_folder",
			expectError:   true,
			expectedError: "destination folder 'invalid_folder' does not exist",
		},
		// Invalid move - moving to a different organization
		{
			name:          "Invalid move - bravo to foxtrot (different org)",
			source:        "bravo",
			dest:          "foxtrot",
			expectError:   true,
			expectedError: "cannot move folder 'bravo' to a different organization",
		},
		// Valid move - bravo to delta
		{
			name:        "Valid move - bravo to delta",
			source:      "bravo",
			dest:        "delta",
			expectError: false,
			validateFunc: func(t *testing.T, folderMap map[string]*folder.Folder) {
				bravo := folderMap["bravo"]
				delta := folderMap["delta"]
				alpha := folderMap["alpha"]

				// Check bravo's new path and parent
				assert.Equal(t, "alpha.delta.bravo", bravo.Paths, "bravo's path should be updated correctly")
				assert.Equal(t, delta, bravo.Parent, "bravo's parent should be delta")

				// Check that bravo is now a child of delta
				found := false
				for _, child := range delta.Children {
					if child == bravo {
						found = true
						break
					}
				}
				assert.True(t, found, "bravo should be in delta's children")

				// Check that bravo is removed from alpha's children
				for _, child := range alpha.Children {
					assert.NotEqual(t, bravo, child, "bravo should be removed from alpha's children")
				}

				// Check charlie's path is updated correctly
				if len(bravo.Children) > 0 {
					charlie := bravo.Children[0] // Assuming charlie is bravo's only child
					assert.Equal(t, "alpha.delta.bravo.charlie", charlie.Paths, "charlie's path should be updated correctly")
					assert.Equal(t, bravo, charlie.Parent, "charlie's parent should be bravo")
				}
			},
		},
		// Valid move - echo to bravo
		{
			name:        "Valid move - echo to bravo",
			source:      "echo",
			dest:        "bravo",
			expectError: false,
			validateFunc: func(t *testing.T, folderMap map[string]*folder.Folder) {
				echo := folderMap["echo"]
				bravo := folderMap["bravo"]
				delta := folderMap["delta"]

				// Check echo's new path and parent
				assert.Equal(t, "alpha.bravo.echo", echo.Paths, "echo's path should be updated correctly")
				assert.Equal(t, bravo, echo.Parent, "echo's parent should be bravo")

				// Check that echo is now a child of bravo
				found := false
				for _, child := range bravo.Children {
					if child == echo {
						found = true
						break
					}
				}
				assert.True(t, found, "echo should be in bravo's children")

				// Check that echo is removed from delta's children
				for _, child := range delta.Children {
					assert.NotEqual(t, echo, child, "echo should be removed from delta's children")
				}
			},
		},
		// Valid move - golf to alpha
		{
			name:        "Valid move - golf to alpha",
			source:      "golf",
			dest:        "alpha",
			expectError: false,
			validateFunc: func(t *testing.T, folderMap map[string]*folder.Folder) {
				golf := folderMap["golf"]
				alpha := folderMap["alpha"]

				// Check golf's new path and parent
				assert.Equal(t, "alpha.golf", golf.Paths, "golf's path should be updated correctly")
				assert.Equal(t, alpha, golf.Parent, "golf's parent should be alpha")

				// Check that golf is now a child of alpha
				found := false
				for _, child := range alpha.Children {
					if child == golf {
						found = true
						break
					}
				}
				assert.True(t, found, "golf should be in alpha's children")
			},
		},
		// Valid move - delta to golf
		{
			name:        "Valid move - delta to golf",
			source:      "delta",
			dest:        "golf",
			expectError: false,
			validateFunc: func(t *testing.T, folderMap map[string]*folder.Folder) {
				delta := folderMap["delta"]
				golf := folderMap["golf"]
				alpha := folderMap["alpha"]
				// echo := folderMap["echo"]

				// Check delta's new path and parent
				assert.Equal(t, "golf.delta", delta.Paths, "delta's path should be updated correctly")
				assert.Equal(t, golf, delta.Parent, "delta's parent should be golf")

				// Check that delta is now a child of golf
				found := false
				for _, child := range golf.Children {
					if child == delta {
						found = true
						break
					}
				}
				assert.True(t, found, "delta should be in golf's children")

				// Check that delta is removed from alpha's children
				for _, child := range alpha.Children {
					assert.NotEqual(t, delta, child, "delta should be removed from alpha's children")
				}

				// Check that echo's path is updated correctly
				if len(delta.Children) > 0 {
					echo := delta.Children[0] // Assuming echo is delta's only child
					assert.Equal(t, "golf.delta.echo", echo.Paths, "echo's path should be updated correctly")
					assert.Equal(t, delta, echo.Parent, "echo's parent should be delta")
				}
			},
		},
		// Valid move - moving a folder with no children
		{
			name:        "Valid move - moving golf under delta",
			source:      "golf",
			dest:        "delta",
			expectError: false,
			validateFunc: func(t *testing.T, folderMap map[string]*folder.Folder) {
				golf := folderMap["golf"]
				delta := folderMap["delta"]

				// Check golf's new path and parent
				assert.Equal(t, "alpha.delta.golf", golf.Paths, "golf's path should be updated correctly")
				assert.Equal(t, delta, golf.Parent, "golf's parent should be delta")

				// Check that golf is now a child of delta
				found := false
				for _, child := range delta.Children {
					if child == golf {
						found = true
						break
					}
				}
				assert.True(t, found, "golf should be in delta's children")
			},
		},
		// Valid move - moving a folder to become a sibling of its parent
		{
			name:        "Valid move - moving bravo under golf",
			source:      "bravo",
			dest:        "golf",
			expectError: false,
			validateFunc: func(t *testing.T, folderMap map[string]*folder.Folder) {
				bravo := folderMap["bravo"]
				golf := folderMap["golf"]
				alpha := folderMap["alpha"]

				// Check bravo's new path and parent
				assert.Equal(t, "golf.bravo", bravo.Paths, "bravo's path should be updated correctly")
				assert.Equal(t, golf, bravo.Parent, "bravo's parent should be golf")

				// Check that bravo is now a child of golf
				found := false
				for _, child := range golf.Children {
					if child == bravo {
						found = true
						break
					}
				}
				assert.True(t, found, "bravo should be in golf's children")

				// Check that bravo is removed from alpha's children
				for _, child := range alpha.Children {
					assert.NotEqual(t, bravo, child, "bravo should be removed from alpha's children")
				}

				// Check charlie's path is updated correctly
				if len(bravo.Children) > 0 {
					charlie := bravo.Children[0] // Assuming charlie is bravo's only child
					assert.Equal(t, "golf.bravo.charlie", charlie.Paths, "charlie's path should be updated correctly")
					assert.Equal(t, bravo, charlie.Parent, "charlie's parent should be bravo")
				}
			},
		},
		// Invalid move - moving a root folder under itself
		{
			name:          "Invalid move - moving alpha under itself",
			source:        "alpha",
			dest:          "alpha",
			expectError:   true,
			expectedError: "cannot move folder 'alpha' to itself",
		},
		// Invalid move - moving a folder to a non-existent destination
		{
			name:          "Invalid move - moving echo to a non-existent destination",
			source:        "echo",
			dest:          "nonexistent",
			expectError:   true,
			expectedError: "destination folder 'nonexistent' does not exist",
		},
		// Invalid move - moving a folder with cyclic reference
		{
			name:          "Invalid move - moving alpha under charlie (creates cyclic reference)",
			source:        "alpha",
			dest:          "charlie",
			expectError:   true,
			expectedError: "cannot move folder 'alpha' to a child of itself",
		},
		// Edge case: moving a folder to the root level (assuming not supported)
		{
			name:          "Invalid move - moving bravo to root level (empty dest)",
			source:        "bravo",
			dest:          "",
			expectError:   true,
			expectedError: "destination folder '' does not exist",
		},
	}

	// Run through each test case
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Initialize folders and driver
			folders, folderMap := initializeFolders(orgID1, orgID2)
			driver := folder.NewDriver(folders)

			_, err := driver.MoveFolder(tt.source, tt.dest)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)

				// Run validation function for checking paths, parent, and children
				if tt.validateFunc != nil {
					tt.validateFunc(t, folderMap)
				}
			}
		})
	}
}

// test function for multiple MoveFolder operations
func Test_folder_MoveFolder_MultipleOperations(t *testing.T) {
	// Sample UUIDs for testing
	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())

	// Initialize folders and driver
	folders, folderMap := initializeFolders(orgID1, orgID2)
	driver := folder.NewDriver(folders)

	// Perform the first MoveFolder operation
	_, err := driver.MoveFolder("bravo", "delta")
	assert.NoError(t, err, "First move (bravo to delta) should succeed")

	// Validate the folder structure after the first operation
	{
		bravo := folderMap["bravo"]
		delta := folderMap["delta"]
		alpha := folderMap["alpha"]

		// Validate bravo
		assert.Equal(t, "alpha.delta.bravo", bravo.Paths)
		assert.Equal(t, delta, bravo.Parent)

		// Check that bravo is now a child of delta
		foundBravo := false
		for _, child := range delta.Children {
			if child == bravo {
				foundBravo = true
				break
			}
		}
		assert.True(t, foundBravo, "After first move, delta should have bravo as a child")

		// Check that bravo is removed from alpha's children
		for _, child := range alpha.Children {
			assert.NotEqual(t, bravo, child, "After first move, bravo should not be in alpha's children")
		}
	}

	// Perform the second MoveFolder operation
	_, err = driver.MoveFolder("delta", "golf")
	assert.NoError(t, err, "Second move (delta to golf) should succeed")

	// Validate the folder structure after the second operation
	{
		delta := folderMap["delta"]
		golf := folderMap["golf"]
		alpha := folderMap["alpha"]

		// Validate delta
		assert.Equal(t, "golf.delta", delta.Paths)
		assert.Equal(t, golf, delta.Parent)

		// Check that delta is now a child of golf
		foundDelta := false
		for _, child := range golf.Children {
			if child == delta {
				foundDelta = true
				break
			}
		}
		assert.True(t, foundDelta, "After second move, golf should have delta as a child")

		// Check that delta is removed from alpha's children
		for _, child := range alpha.Children {
			assert.NotEqual(t, delta, child, "After second move, delta should not be in alpha's children")
		}

		// Since bravo was a child of delta, check bravo's path and parent
		bravo := folderMap["bravo"]
		assert.Equal(t, "golf.delta.bravo", bravo.Paths)
		assert.Equal(t, delta, bravo.Parent)
	}

	// Perform the third MoveFolder operation
	_, err = driver.MoveFolder("echo", "bravo")
	assert.NoError(t, err, "Third move (echo to bravo) should succeed")

	// Validate the folder structure after the third operation
	{
		echo := folderMap["echo"]
		bravo := folderMap["bravo"]
		delta := folderMap["delta"]

		// Validate echo
		assert.Equal(t, "golf.delta.bravo.echo", echo.Paths)
		assert.Equal(t, bravo, echo.Parent)

		// Check that echo is now a child of bravo
		foundEcho := false
		for _, child := range bravo.Children {
			if child == echo {
				foundEcho = true
				break
			}
		}
		assert.True(t, foundEcho, "After third move, bravo should have echo as a child")

		// Check that echo is removed from delta's children
		for _, child := range delta.Children {
			assert.NotEqual(t, echo, child, "After third move, echo should not be in delta's children")
		}
	}
}
