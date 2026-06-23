package gotenberg

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChromiumSecurityAndEmbeds(t *testing.T) {
	var receivedParams = make(map[string]string)
	var receivedFiles = make(map[string]string)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forms/chromium/convert/html" {
			t.Errorf("expected path /forms/chromium/convert/html, got %s", r.URL.Path)
		}

		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatal(err)
		}

		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}

			name := part.FormName()
			filename := part.FileName()
			content, err := io.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}

			if filename != "" {
				receivedFiles[name] = filename + ":" + string(content)
			} else {
				receivedParams[name] = string(content)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF_CONTENT"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Chromium().
		ConvertHTML(context.Background(), strings.NewReader("<html></html>")).
		UserPassword("userpass").
		OwnerPassword("ownerpass").
		AllowPrinting(false).
		AllowCopying(false).
		AllowModifying(false).
		AllowAnnotating(false).
		AllowFillingForms(false).
		AllowAssembling(false).
		Embeds("attachment.xml", strings.NewReader("<xml></xml>")).
		Send()

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Assert parameters
	expectedParams := map[string]string{
		"userPassword":      "userpass",
		"ownerPassword":     "ownerpass",
		"allowPrinting":     "false",
		"allowCopying":      "false",
		"allowModifying":    "false",
		"allowAnnotating":   "false",
		"allowFillingForms": "false",
		"allowAssembling":   "false",
	}

	for k, v := range expectedParams {
		if receivedParams[k] != v {
			t.Errorf("expected param %s to be %s, got %s", k, v, receivedParams[k])
		}
	}

	// Assert files
	if receivedFiles["files"] != "index.html:<html></html>" {
		t.Errorf("expected files to contain index.html, got %s", receivedFiles["files"])
	}
	if receivedFiles["embeds"] != "attachment.xml:<xml></xml>" {
		t.Errorf("expected embeds to contain attachment.xml, got %s", receivedFiles["embeds"])
	}
}

func TestLibreOfficeSecurityAndEmbeds(t *testing.T) {
	var receivedParams = make(map[string]string)
	var receivedFiles = make(map[string]string)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forms/libreoffice/convert" {
			t.Errorf("expected path /forms/libreoffice/convert, got %s", r.URL.Path)
		}

		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatal(err)
		}

		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}

			name := part.FormName()
			filename := part.FileName()
			content, err := io.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}

			if filename != "" {
				receivedFiles[name] = filename + ":" + string(content)
			} else {
				receivedParams[name] = string(content)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF_CONTENT"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.LibreOffice().
		Convert(context.Background()).
		File("document.docx", strings.NewReader("docx-data")).
		UserPassword("userpass").
		OwnerPassword("ownerpass").
		AllowPrinting(true).
		Embeds("embed.txt", strings.NewReader("text-data")).
		Send()

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if receivedParams["userPassword"] != "userpass" {
		t.Errorf("expected userPassword userpass, got %s", receivedParams["userPassword"])
	}
	if receivedParams["ownerPassword"] != "ownerpass" {
		t.Errorf("expected ownerPassword ownerpass, got %s", receivedParams["ownerPassword"])
	}
	if receivedParams["allowPrinting"] != "true" {
		t.Errorf("expected allowPrinting true, got %s", receivedParams["allowPrinting"])
	}
	if receivedFiles["files"] != "document.docx:docx-data" {
		t.Errorf("expected files to contain document.docx, got %s", receivedFiles["files"])
	}
	if receivedFiles["embeds"] != "embed.txt:text-data" {
		t.Errorf("expected embeds to contain embed.txt, got %s", receivedFiles["embeds"])
	}
}

func TestPDFEnginesEmbed(t *testing.T) {
	var receivedParams = make(map[string]string)
	var receivedFiles = make(map[string]string)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forms/pdfengines/embed" {
			t.Errorf("expected path /forms/pdfengines/embed, got %s", r.URL.Path)
		}

		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatal(err)
		}

		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}

			name := part.FormName()
			filename := part.FileName()
			content, err := io.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}

			if filename != "" {
				if _, ok := receivedFiles[name]; ok {
					receivedFiles[name] += ";" + filename + ":" + string(content)
				} else {
					receivedFiles[name] = filename + ":" + string(content)
				}
			} else {
				receivedParams[name] = string(content)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF_CONTENT"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.PDFEngines().
		Embed(context.Background()).
		File("document.pdf", strings.NewReader("pdf-data")).
		Embeds("attachment.xml", strings.NewReader("xml-data")).
		Send()

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if receivedFiles["files"] != "document.pdf:pdf-data" {
		t.Errorf("expected files to contain document.pdf, got %s", receivedFiles["files"])
	}
	if receivedFiles["embeds"] != "attachment.xml:xml-data" {
		t.Errorf("expected embeds to contain attachment.xml, got %s", receivedFiles["embeds"])
	}
}

func TestPDFEnginesFacturX(t *testing.T) {
	var receivedParams = make(map[string]string)
	var receivedFiles = make(map[string]string)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forms/pdfengines/factur-x" {
			t.Errorf("expected path /forms/pdfengines/factur-x, got %s", r.URL.Path)
		}

		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatal(err)
		}

		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}

			name := part.FormName()
			filename := part.FileName()
			content, err := io.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}

			if filename != "" {
				receivedFiles[name] = filename + ":" + string(content)
			} else {
				receivedParams[name] = string(content)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF_CONTENT"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.PDFEngines().
		FacturX(context.Background()).
		File("document.pdf", strings.NewReader("pdf-data")).
		FacturXXml("factur-x.xml", strings.NewReader("xml-data")).
		FacturXConformanceLevel("BASIC").
		FacturXDocumentType("INVOICE").
		FacturXVersion("1.0").
		Send()

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if receivedFiles["files"] != "document.pdf:pdf-data" {
		t.Errorf("expected files to contain document.pdf, got %s", receivedFiles["files"])
	}
	if receivedFiles["facturxXml"] != "factur-x.xml:xml-data" {
		t.Errorf("expected facturxXml to contain factur-x.xml, got %s", receivedFiles["facturxXml"])
	}
	if receivedParams["facturxConformanceLevel"] != "BASIC" {
		t.Errorf("expected conformance level BASIC, got %s", receivedParams["facturxConformanceLevel"])
	}
	if receivedParams["facturxDocumentType"] != "INVOICE" {
		t.Errorf("expected document type INVOICE, got %s", receivedParams["facturxDocumentType"])
	}
	if receivedParams["facturxVersion"] != "1.0" {
		t.Errorf("expected version 1.0, got %s", receivedParams["facturxVersion"])
	}
}

func BenchmarkChromiumBuilder(b *testing.B) {
	client, err := NewClient(&http.Client{}, "http://localhost:3000")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.Chromium().
			ConvertURL(context.Background(), "http://localhost").
			UserPassword("userpass").
			OwnerPassword("ownerpass").
			AllowPrinting(false).
			AllowCopying(false)
	}
}

func BenchmarkPDFEnginesBuilder(b *testing.B) {
	client, err := NewClient(&http.Client{}, "http://localhost:3000")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.PDFEngines().
			Merge(context.Background()).
			FacturXConformanceLevel("BASIC").
			FacturXDocumentType("INVOICE").
			FacturXVersion("1.0")
	}
}


