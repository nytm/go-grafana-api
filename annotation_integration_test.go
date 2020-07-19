// +build integration

package gapi

import (
	"fmt"
	"net/url"
	"testing"
)

func TestAnnotationsIntegration(t *testing.T) {
	client, err := New("admin:admin", "http://localhost:3000")
	if err != nil {
		t.Error(err)
	}

	_, err = client.Annotations(url.Values{})
	if err != nil {
		t.Error(err)
	}
}

func TestNewAnnotationIntegration(t *testing.T) {
	client, err := New("admin:admin", "http://localhost:3000")
	if err != nil {
		t.Error(err)
	}

	_, err = client.NewAnnotation(&Annotation{
		Text: "integration-test-new",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateAnnotationIntegration(t *testing.T) {
	client, err := New("admin:admin", "http://localhost:3000")
	if err != nil {
		t.Error(err)
	}

	id, err := client.NewAnnotation(&Annotation{
		Text: "integration-test-update",
	})
	if err != nil {
		t.Error(err)
	}

	message, err := client.UpdateAnnotation(id, &Annotation{
		Text: "integration-test-updated",
	})
	if err != nil {
		t.Error(err)
	}

	expected := "Annotation updated"
	if message != expected {
		t.Error(fmt.Sprintf("expected UpdateAnnotation message to be %s; got %s", expected, message))
	}
}

func TestDeleteAnnotationIntegration(t *testing.T) {
	client, err := New("admin:admin", "http://localhost:3000")
	if err != nil {
		t.Error(err)
	}

	id, err := client.NewAnnotation(&Annotation{
		Text: "integration-test-delete",
	})
	if err != nil {
		t.Error(err)
	}

	message, err := client.DeleteAnnotation(id)
	if err != nil {
		t.Error(err)
	}

	expected := "Annotation deleted"
	if message != expected {
		t.Error(fmt.Sprintf("expected DeleteAnnotation message to be %s; got %s", expected, message))
	}
}
