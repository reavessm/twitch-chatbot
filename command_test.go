package chatbot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/twitch-chatbot"
)

func TestIsCommand(t *testing.T) {
	testCases := map[string]bool{
		"":       false,
		" ":      false,
		"find":   false,
		"#find":  false,
		"!":      false,
		"! ":     false,
		"! !":    false,
		"!.":     false,
		"!.find": false,
		"!!":     false,
		"!123":   false,

		"!find":       true,
		" !find":      true,
		" !find ":     true,
		"  !  find  ": true,
		"! find":      true,
		"! find ":     true,
	}

	const prefix = '!'

	for text, expect := range testCases {
		t.Run(text, func(t *testing.T) {
			assert.Equal(t, expect, chatbot.IsCommand(prefix, text))
		})
	}
}

func BenchmarkIsCommand(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		chatbot.IsCommand('!', "!find xantom")
	}
}

func TestNewCommandFromMessage(t *testing.T) {
	t.Run("invalid prefix", func(t *testing.T) {
		cmd, err := chatbot.NewCommandFromMessage('!', "#bar")
		assert.Equal(t, cmd, chatbot.Command{})
		assert.EqualError(t, err, "unable to parse command call")
	})

	t.Run("invalid command", func(t *testing.T) {
		cmd, err := chatbot.NewCommandFromMessage('#', "##")
		assert.Equal(t, cmd, chatbot.Command{})
		assert.EqualError(t, err, "unable to parse command call")
	})

	t.Run("valid", func(t *testing.T) {
		testCases := map[string]chatbot.Command{
			"!find":       {Name: "find", Args: []string{}},
			" !find arg1": {Name: "find", Args: []string{"arg1"}},
			"!find a b c": {Name: "find", Args: []string{"a", "b", "c"}},
		}

		const prefix = '!'

		for text, expect := range testCases {
			t.Run(text, func(t *testing.T) {
				foo, err := chatbot.NewCommandFromMessage(prefix, text)
				assert.Equal(t, expect, foo)
				assert.Nil(t, err)
			})
		}
	})
}

func BenchmarkNewCommandFromMessage(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		chatbot.NewCommandFromMessage('!', "!find xantom")
	}
}

func TestArgsToString(t *testing.T) {
	assert.Equal(t, "foo bar", chatbot.NewCommand("find", "foo", "bar").ArgsToString())
}
