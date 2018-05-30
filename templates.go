// FIXME: Dislike this file a bit - what's the take on referencing
// viper config values (treat it as a global, or pass values around?)
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// compile all templates and cache them
var templates = template.Must(template.ParseGlob("./templates/*.html"))

func renderIndex(w http.ResponseWriter, config *Config) {
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, config)
}

type templateData struct {
	IDToken          string
	RefreshToken     string
	RedirectURL      string
	Claims           string
	Username         string
	Issuer           string
	ClusterName      string
	ShortDescription string
	ClientSecret     string
	ClientID         string
	K8sMasterURI     string
	K8sCaURI         string
	K8sCaPem         string
	IDPCaURI         string
	LogoURI          string
	KubectlVersion   string
}

func (cluster *Cluster) renderToken(w http.ResponseWriter,
	idToken,
	refreshToken string,
	idpCaURI string,
	logoURI string,
	kubectlVersion string,
	claims []byte) {

	var data map[string]interface{}
	err := json.Unmarshal(claims, &data)
	if err != nil {
		panic(err)
	}

	token_data := templateData{
		IDToken:          idToken,
		RefreshToken:     refreshToken,
		RedirectURL:      cluster.Redirect_URI,
		Claims:           string(claims),
		Username:         data["email"].(string),
		Issuer:           data["iss"].(string),
		ClusterName:      cluster.Name,
		ShortDescription: cluster.Short_Description,
		ClientSecret:     cluster.Client_Secret,
		ClientID:         cluster.Client_ID,
		K8sMasterURI:     cluster.K8s_Master_URI,
		K8sCaURI:         cluster.K8s_Ca_URI,
		K8sCaPem:         cluster.K8s_Ca_Pem,
		IDPCaURI:         idpCaURI,
		LogoURI:          logoURI,
		KubectlVersion:   kubectlVersion}

	err = templates.ExecuteTemplate(w, "kubeconfig.html", token_data)

	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}
}
