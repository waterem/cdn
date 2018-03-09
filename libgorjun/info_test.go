package gorjun

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

// Test for info rest
// Two user uploads same 5 templates
// Each user requests by his token
func TestList(t *testing.T) {
	g := NewGorjunServer()
	g.Register(g.Username)

	k := SecondNewGorjunServer()
	k.Register(k.Username)

	repos := [2]string{"template", "raw"}

	for i:= 0; i < len(repos) ;i++ {

		artifactType := repos[i]

		err := k.Uploads(artifactType, "false")
		if err != nil {
			t.Errorf("Failed to uploads templates: %v", err)
		}

		err = g.Uploads(artifactType, "false")

		if err != nil {
			t.Errorf("Failed to uploads templates: %v", err)
		}

		flist, err := g.List(artifactType, "")
		if err != nil {
			t.Errorf("Failed to retrieve user files: %v", err)
		}
		if len(flist) <= 0 {
			t.Errorf("Resulting array is empty")
		}

		flist, err = g.List(artifactType, "?token")
		if err != nil {
			t.Errorf("Failed to retrieve user files: %v", err)
		}
		if len(flist) <= 0 {
			t.Errorf("Resulting array is empty")
		}

		fmt.Println("Token for user " + g.Username + " = " + g.Token)
		fmt.Println("Token for user " + k.Username + " = " + k.Token)

		flist, err = g.List(artifactType, "?name=nginx&owner="+g.Username)
		owner := flist[0].Owner[0]
		assert.Equal(t, owner, g.Username)

		flist, err = g.List(artifactType, "?name=nginx&owner="+k.Username)
		owner = flist[0].Owner[0]
		assert.Equal(t, owner, k.Username)

		flist, err = g.List(artifactType, "?name=nginx&token="+g.Token)
		owner = flist[0].Owner[0]
		assert.Equal(t, owner, g.Username)

		flist, err = g.List(artifactType, "?name=nginx&token="+k.Token)
		owner = flist[0].Owner[0]
		assert.Equal(t, owner, k.Username)

		err = g.Deletes(artifactType, "")
		if err != nil {
			t.Errorf("Failed to delete templates: %v", err)
		}

		err = k.Deletes(artifactType, "")
		if err != nil {
			t.Errorf("Failed to delete templates: %v", err)
		}
	}
}

// Test for info rest
// Two user uploads same 5 private templates
// Each user requests by his token
func TestPrivateTemplates(t *testing.T) {
	g := NewGorjunServer()
	g.Register(g.Username)

	k := SecondNewGorjunServer()
	k.Register(k.Username)

	repos := [2]string{"template", "raw"}

	for i:= 0; i < len(repos) ;i++ {

		artifactType := repos[i]

		err := k.Uploads(artifactType, "true")
		if err != nil {
			t.Errorf("Failed to uploads templates: %v", err)
		}

		err = g.Uploads(artifactType, "true")
		if err != nil {
			t.Errorf("Failed to uploads templates: %v", err)
		}

		fmt.Println("Token for user " + g.Username + " = " + g.Token)
		fmt.Println("Token for user " + k.Username + " = " + k.Token)

		flist, err := g.List(artifactType, "?name=nginx&token="+g.Token)
		owner := flist[0].Owner[0]
		assert.Equal(t, owner, g.Username)

		flist, err = g.List(artifactType, "?name=nginx&token="+k.Token)
		owner = flist[0].Owner[0]
		assert.Equal(t, owner, k.Username)

		err = g.Deletes(artifactType, "?token="+g.Token)
		if err != nil {
			t.Errorf("Failed to delete templates: %v", err)
		}

		err = k.Deletes(artifactType, "?token="+k.Token)
		if err != nil {
			t.Errorf("Failed to delete templates: %v", err)
		}
	}
}

// Test for info rest
// If can't find user templates it should search
// in shared
func TestShareTemplates(t *testing.T) {
	g := NewGorjunServer()
	g.Register(g.Username)

	k := SecondNewGorjunServer()
	k.Register(k.Username)
	k.AuthenticateUser()

	artifactType := "template"

	err := g.Uploads(artifactType, "false")
	if err != nil {
		t.Errorf("Failed to uploads templates: %v", err)
	}

	fmt.Println("Token for user " + g.Username + " = " + g.Token)
	fmt.Println("Token for user " + k.Username + " = " + k.Token)
	flist, err := g.List(artifactType, "")
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	g.Share(g.Token,flist,k.Username,artifactType)

	err = g.Deletes(artifactType, "?token="+g.Token)
	if err != nil {
		t.Errorf("Failed to delete templates: %v", err)
	}
}

//Test for info?name=master
func TestListByName(t *testing.T) {
	v := VerifiedUser()
	v.Register(v.Username)

	g := SecondNewGorjunServer()
	g.Register(g.Username)

	artifactType := "template"

	err := g.Uploads(artifactType, "false")
	if err != nil {
		t.Errorf("Failed to uploads %s: %v", err,artifactType)
	}

	err  = v.Uploads(artifactType, "false")
	if err != nil {
		t.Errorf("Failed to uploads %s: %v", err,artifactType)
	}

	fmt.Println("Token for user " + g.Username + " = " + g.Token)
	fmt.Println("Token for user " + v.Username + " = " + v.Token)

	flist, err := g.List(artifactType, "?name=nginx")
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	owner := flist[0].Owner[0]
	version := flist[0].Version
	assert.Equal(t, owner, v.Username)
	assert.Equal(t, version, "0.1.11")

	flist, err = g.List(artifactType, "?name=nginx&token=" + g.Token)
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	owner = flist[0].Owner[0]
	version = flist[0].Version
	assert.Equal(t, owner, g.Username)
	assert.Equal(t, version, "0.1.11")

	flist, err = g.List(artifactType, "?name=nginx&owner=" + g.Username)
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	owner = flist[0].Owner[0]
	version = flist[0].Version
	assert.Equal(t, owner, g.Username)
	assert.Equal(t, version, "0.1.11")

	flist, err = g.List(artifactType, "?name=nginx&owner=" + v.Username)
	if err != nil {
		t.Errorf("Failed to retrieve user files: %v", err)
	}
	owner = flist[0].Owner[0]
	version = flist[0].Version
	assert.Equal(t, owner, v.Username)
	assert.Equal(t, version, "0.1.11")
	err = g.Deletes(artifactType, "?token="+g.Token)
	if err != nil {
		t.Errorf("Failed to delete templates: %v", err)
	}
	err = v.Deletes(artifactType, "?token="+v.Token)
	if err != nil {
		t.Errorf("Failed to delete templates: %v", err)
	}
}