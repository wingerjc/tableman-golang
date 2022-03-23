package web

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// type utilData struct {
// 	comp *compiler.Compiler
// 	prog *program.Program
// }

// func setup() *utilData {
// 	c, _ := compiler.NewCompiler()
// 	return &utilData{
// 		comp: c,
// 		prog: program.NewProgram(make(program.TableMap)),
// 	}
// }

func testingTimer() func() time.Time {
	t := time.Now()
	return func() time.Time {
		t = t.Add(time.Second)
		return t
	}
}

func TestSessionDropping(t *testing.T) {
	assert := assert.New(t)
	set := NewSessionSet(3, 3*time.Second)
	assert.Empty(set.sessions)
	timeNow = testingTimer()

	// Fill up the session set.
	first := set.NewSession()
	second := set.NewSession()
	third := set.NewSession()
	assert.Len(set.sessions, 3)
	assert.Contains(set.sessions, first)
	assert.Contains(set.sessions, second)
	assert.Contains(set.sessions, third)

	// Add an extra session and make sure the oldest is dropped.
	set.NewSession()
	assert.Len(set.sessions, 3)
	assert.NotContains(set.sessions, first)
	assert.Contains(set.sessions, second)
	assert.Contains(set.sessions, third)

	// Touch the second session to make the third one get deleted for this drop.
	set.sessions[second].Touch()
	set.NewSession()
	assert.Len(set.sessions, 3)
	assert.NotContains(set.sessions, third)
	assert.Contains(set.sessions, second)
}
