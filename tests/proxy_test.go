package tests

import (
	"ProxyBuild/proxy"
	"testing"
)

func TestShouldExecuteHook_NoConditions(t *testing.T) {
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"test"},
		When:    "before",
	}

	if !proxy.ShouldExecuteHook(hook, []string{"up"}, false, "win") {
		t.Error("Hook without conditions should always execute")
	}
}

func TestShouldExecuteHook_OnErrorTrue(t *testing.T) {
	trueVal := true
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"error"},
		When:    "after",
		Conditions: proxy.Conditions{
			OnError: &trueVal,
		},
	}

	// Should execute when there was an error
	if !proxy.ShouldExecuteHook(hook, []string{"up"}, true, "win") {
		t.Error("Hook with on_error:true should execute when there was an error")
	}

	// Should NOT execute when there was no error
	if proxy.ShouldExecuteHook(hook, []string{"up"}, false, "win") {
		t.Error("Hook with on_error:true should NOT execute when there was no error")
	}
}

func TestShouldExecuteHook_OnErrorFalse(t *testing.T) {
	falseVal := false
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"success"},
		When:    "after",
		Conditions: proxy.Conditions{
			OnError: &falseVal,
		},
	}

	// Should execute when there was no error
	if !proxy.ShouldExecuteHook(hook, []string{"up"}, false, "win") {
		t.Error("Hook with on_error:false should execute when there was no error")
	}

	// Should NOT execute when there was an error
	if proxy.ShouldExecuteHook(hook, []string{"up"}, true, "win") {
		t.Error("Hook with on_error:false should NOT execute when there was an error")
	}
}

func TestShouldExecuteHook_ArgsContain(t *testing.T) {
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"detached"},
		When:    "after",
		Conditions: proxy.Conditions{
			ArgsContain: []string{"-d"},
		},
	}

	// Should execute when args contain -d
	if !proxy.ShouldExecuteHook(hook, []string{"up", "-d"}, false, "win") {
		t.Error("Hook should execute when args contain required string")
	}

	// Should NOT execute when args don't contain -d
	if proxy.ShouldExecuteHook(hook, []string{"up"}, false, "win") {
		t.Error("Hook should NOT execute when args don't contain required string")
	}
}

func TestShouldExecuteHook_ArgsContainMultiple(t *testing.T) {
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"cleanup"},
		When:    "after",
		Conditions: proxy.Conditions{
			ArgsContain: []string{"--volumes", "down"},
		},
	}

	// Should execute when args contain both strings
	if !proxy.ShouldExecuteHook(hook, []string{"down", "--volumes"}, false, "win") {
		t.Error("Hook should execute when args contain all required strings")
	}

	// Should NOT execute when args contain only one string
	if proxy.ShouldExecuteHook(hook, []string{"down"}, false, "win") {
		t.Error("Hook should NOT execute when args don't contain all required strings")
	}
}

func TestShouldExecuteHook_ArgsMatch(t *testing.T) {
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"following logs"},
		When:    "before",
		Conditions: proxy.Conditions{
			ArgsMatch: []string{"-f"},
		},
	}

	// Should execute when args match exactly
	if !proxy.ShouldExecuteHook(hook, []string{"logs", "-f"}, false, "win") {
		t.Error("Hook should execute when args match exactly")
	}

	// Should NOT execute when args don't match
	if proxy.ShouldExecuteHook(hook, []string{"logs", "--follow"}, false, "win") {
		t.Error("Hook should NOT execute when args don't match exactly")
	}
}

func TestShouldExecuteHook_CombinedConditions(t *testing.T) {
	falseVal := false
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"success with volumes"},
		When:    "after",
		Conditions: proxy.Conditions{
			OnError:     &falseVal,
			ArgsContain: []string{"--volumes"},
		},
	}

	// Should execute when all conditions are met
	if !proxy.ShouldExecuteHook(hook, []string{"down", "--volumes"}, false, "win") {
		t.Error("Hook should execute when all conditions are met")
	}

	// Should NOT execute when error occurred (even with correct args)
	if proxy.ShouldExecuteHook(hook, []string{"down", "--volumes"}, true, "win") {
		t.Error("Hook should NOT execute when error occurred")
	}

	// Should NOT execute when args don't match (even without error)
	if proxy.ShouldExecuteHook(hook, []string{"down"}, false, "win") {
		t.Error("Hook should NOT execute when args don't match")
	}
}

func TestShouldExecuteHook_ArgsContainSubstring(t *testing.T) {
	hook := proxy.Hook{
		Command: "echo",
		Args:    []string{"port mapping"},
		When:    "after",
		Conditions: proxy.Conditions{
			ArgsContain: []string{"8080"},
		},
	}

	// Should execute when substring is found
	if !proxy.ShouldExecuteHook(hook, []string{"up", "-p", "8080:80"}, false, "win") {
		t.Error("Hook should execute when substring is found in args")
	}

	// Should NOT execute when substring is not found
	if proxy.ShouldExecuteHook(hook, []string{"up", "-p", "9090:90"}, false, "win") {
		t.Error("Hook should NOT execute when substring is not found")
	}
}

func TestConfig_JSONStructure(t *testing.T) {
	config := proxy.Config{
		BaseCommand: "docker-compose",
		Hooks: map[string][]proxy.Hook{
			"up": {
				{
					Command: "echo",
					Args:    []string{"starting"},
					When:    "before",
				},
			},
		},
	}

	if config.BaseCommand != "docker-compose" {
		t.Error("BaseCommand not set correctly")
	}

	if len(config.Hooks) != 1 {
		t.Error("Hooks not set correctly")
	}

	if len(config.Hooks["up"]) != 1 {
		t.Error("Hook count not correct")
	}
}
