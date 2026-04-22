package console

import "testing"

func TestBoxErrorItem(t *testing.T) {
	got := BoxErrorItem("desc/api/user.api")
	want := "│  " + CrossMark() + " desc/api/user.api"
	if got != want {
		t.Fatalf("BoxErrorItem() = %q, want %q", got, want)
	}
}

func TestBoxInfoItem(t *testing.T) {
	got := BoxInfoItem("Executing echo test")
	want := "│  " + Cyan(">") + " Executing echo test"
	if got != want {
		t.Fatalf("BoxInfoItem() = %q, want %q", got, want)
	}
}

func TestBoxFooters(t *testing.T) {
	if got, want := BoxSuccessFooter(), "└─ "+Cyan("✓")+" "+Cyan("Complete"); got != want {
		t.Fatalf("BoxSuccessFooter() = %q, want %q", got, want)
	}

	if got, want := BoxErrorFooter(), "└─ "+CrossMark()+" "+Red("Abort"); got != want {
		t.Fatalf("BoxErrorFooter() = %q, want %q", got, want)
	}
}
