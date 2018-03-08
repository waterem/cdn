package gorjun

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

// Test for info rest
func TestList(t *testing.T) {
	g := NewGorjunServer()
	g.Register(g.Username)

	k := SecondNewGorjunServer()
	k.Register(k.Username)

	artifactType := "template"

	err := k.UploadTemplates()
	if err != nil {
		t.Errorf("Failed to uploads templates: %v", err)
	}

	err = g.UploadTemplates()
	if err != nil {
		t.Errorf("Failed to uploads templates: %v", err)
	}

	flist, err := g.List(artifactType,"")
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	if len(flist) <= 0 {
		t.Errorf("Resulting array is empty")
	}

	flist, err = g.List(artifactType,"?token")
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	if len(flist) <= 0 {
		t.Errorf("Resulting array is empty")
	}

	fmt.Println("Token for user " + g.Username + " = " + g.Token)
	fmt.Println("Token for user " + k.Username + " = " + k.Token)

	flist, err = g.List(artifactType,"?name=nginx&owner=" + g.Username)
	owner := flist[0].Owner[0]
	assert.Equal(t,owner,g.Username)

	flist, err = g.List(artifactType,"?name=nginx&owner=" + k.Username)
	owner = flist[0].Owner[0]
	assert.Equal(t,owner,k.Username)

	flist, err = g.List(artifactType,"?name=nginx&token=" + g.Token)
	owner = flist[0].Owner[0]
	assert.Equal(t,owner,g.Username)

	flist, err = g.List(artifactType,"?name=nginx&token=" + k.Token)
	owner = flist[0].Owner[0]
	assert.Equal(t,owner,k.Username)

	err = g.DeleteTemplates()
	if err != nil {
		t.Errorf("Failed to delete templates: %v", err)
	}

}