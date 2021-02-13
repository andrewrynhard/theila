package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/talos-systems/theila/pkg/client"
)

const (
	LayoutTemplateName = "layout"
	CurrentContextKey  = "CurrentContext"
	ContentKey         = "Content"
	AsideKey           = "Aside"
)

type Renderer struct {
	Template *template.Template
	Client   *client.Client
}

func (r *Renderer) Render(name string, content interface{}) (string, error) {
	var buf bytes.Buffer

	err := r.execute(&buf, name, content)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *Renderer) RenderPartial(w http.ResponseWriter, status int, name string, content map[string]interface{}) {
	w.WriteHeader(status)

	err := r.execute(w, name, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Renderer) RenderFull(w http.ResponseWriter, status int, name string, content map[string]interface{}) {
	w.WriteHeader(status)

	var buf bytes.Buffer

	err := r.execute(&buf, name, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	config, err := r.Client.KubernetesClientConfig.ConfigAccess().GetStartingConfig()
	if err != nil {
	}

	data := map[string]interface{}{
		CurrentContextKey: config.CurrentContext,
		ContentKey:        template.HTML(buf.String()),
		AsideKey:          "",
		"Current":         content["Current"],
	}

	err = r.execute(w, LayoutTemplateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Renderer) execute(wr io.Writer, name string, data interface{}) error {
	tmpl := r.Template.Lookup(name)

	if tmpl == nil {
		return fmt.Errorf("template %q is undefined", name)
	}

	err := tmpl.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
