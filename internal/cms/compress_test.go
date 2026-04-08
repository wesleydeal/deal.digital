package cms

import "testing"

func TestCompressionPlanFor(t *testing.T) {
	tests := []struct {
		path   string
		want   compressionPlan
		wantOK bool
	}{
		{path: "index.html", want: compressionPlan{gzip: true, brotli: true, zstd: true}, wantOK: true},
		{path: "style.css", want: compressionPlan{gzip: true, brotli: true, zstd: true}, wantOK: true},
		{path: "script.mjs", want: compressionPlan{gzip: true, brotli: true, zstd: true}, wantOK: true},
		{path: "data.json", want: compressionPlan{gzip: true, brotli: true, zstd: true}, wantOK: true},
		{path: "image.svg", want: compressionPlan{gzip: true, brotli: true, zstd: true}, wantOK: true},
		{path: "font.woff", want: compressionPlan{gzip: true, brotli: false, zstd: true}, wantOK: true},
		{path: "archive.swf", want: compressionPlan{}, wantOK: false},
		{path: "archive.zip", want: compressionPlan{}, wantOK: false},
		{path: "image.png", want: compressionPlan{}, wantOK: false},
		{path: "audio.ogg", want: compressionPlan{}, wantOK: false},
		{path: "bundle.WEBMANIFEST", want: compressionPlan{gzip: true, brotli: true, zstd: true}, wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, ok := compressionPlanFor(tt.path)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("plan = %#v, want %#v", got, tt.want)
			}
		})
	}
}
