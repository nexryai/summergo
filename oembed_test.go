package summergo

import "testing"

func containsString(arr []string, target string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func TestGetRequiredPermissionsFromIframe(t *testing.T) {
	testIframe := "<iframe width=\"200\" height=\"113\" src=\"https://www.youtube.com/embed/zK-RUYiYLok?feature=oembed\" frameborder=\"0\" allow=\"accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share\" allowfullscreen title=\"【崩壊：スターレイル】EP「制御不能」\"></iframe>"

	permissions := getRequiredPermissionsFromIframe(testIframe)

	if !containsString(permissions, "autoplay") {
		t.Errorf("autoplay should be contained")
	} else if !containsString(permissions, "picture-in-picture") {
		t.Errorf("picture-in-picture should be contained")
	}
}
