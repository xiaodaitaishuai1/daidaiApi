package console_setting

import (
	"strings"
	"testing"
)

func TestValidateConsoleSettingsCountsUTF16Characters(t *testing.T) {
	content := strings.Repeat("😀", 251)
	settings := `[{"content":"` + content + `","publishDate":"2026-01-01T00:00:00Z"}]`
	if err := ValidateConsoleSettings(settings, "Announcements"); err == nil {
		t.Fatal("expected UTF-16 length validation to reject 251 astral characters")
	}
}

func TestValidateConsoleSettingsAnnouncementExtraLimit(t *testing.T) {
	extra := strings.Repeat("a", 101)
	settings := `[{"content":"ok","publishDate":"2026-01-01T00:00:00Z","extra":"` + extra + `"}]`
	if err := ValidateConsoleSettings(settings, "Announcements"); err == nil {
		t.Fatal("expected announcement extra over 100 characters to be rejected")
	}
}
