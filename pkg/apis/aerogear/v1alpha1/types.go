package v1alpha1

import (
	"regexp"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	Group             = "aerogear.org"
	Version           = "v1alpha1"
	KeycloakKind      = "Keycloak"
	KeycloakVersion   = "4.1.0"
	KeycloakFinalizer = "finalizer.org.aerogear.keycloak"
)

type Config struct {
	ResyncPeriod  int
	LogLevel      string
	SyncResources bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KeycloakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Keycloak `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// crd:gen:Kind=Keycloak:Group=aerogear.org
type Keycloak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              KeycloakSpec   `json:"spec"`
	Status            KeycloakStatus `json:"status,omitempty"`
}

func (k *Keycloak) Defaults() {
	alphaNum := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	//set the defaults if not set to something else
	for _, r := range k.Spec.Realms {
		for i, u := range r.Users {
			if u.OutputSecret == "" {
				r.Users[i].OutputSecret = alphaNum.ReplaceAllString(u.UserName, "-")
			}
		}
		for i, c := range r.Clients {
			if c.ID == "" {
				c.ID = c.ClientID
				r.Clients[i] = c
			}
			if c.OutputSecret == "" {
				c.OutputSecret = alphaNum.ReplaceAllString(r.Realm+"-"+c.ClientID, "-")
			}
		}
	}

}

func (k *Keycloak) Validate() error {
	return nil
}

type KeycloakSpec struct {
	Version          string          `json:"version"`
	AdminCredentials string          `json:"adminCredentials"`
	Realms           []KeycloakRealm `json:"realms"`
}

type KeycloakRealm struct {
	*KeycloakApiRealm
	Users   []KeycloakUser   `json:"users,omitempty"`
	Clients []KeycloakClient `json:"clients,omitempty"`
}

type KeycloakApiRealm struct {
	ID                string                     `json:"id,omitempty"`
	Realm             string                     `json:"realm,omitempty"`
	Enabled           bool                       `json:"enabled"`
	DisplayName       string                     `json:"displayName"`
	Users             []KeycloakApiUser          `json:"users,omitempty"`
	Clients           []KeycloakApiClient        `json:"clients,omitempty"`
	IdentityProviders []KeycloakIdentityProvider `json:"identityProviders,omitempty"`
}

type KeycloakApiPasswordReset struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type KeycloakIdentityProvider struct {
	Alias                     string            `json:"alias,omitempty"`
	DisplayName               string            `json:"displayName"`
	InternalID                string            `json:"internalId,omitempty"`
	ProviderID                string            `json:"providerId,omitempty"`
	Enabled                   bool              `json:"enabled"`
	TrustEmail                bool              `json:"trustEmail"`
	StoreToken                bool              `json:"storeToken"`
	AddReadTokenRoleOnCreate  bool              `json:"addReadTokenRoleOnCreate"`
	FirstBrokerLoginFlowAlias string            `json:"firstBrokerLoginFlowAlias"`
	PostBrokerLoginFlowAlias  string            `json:"postBrokerLoginFlowAlias"`
	Config                    map[string]string `json:"config"`
}
type KeycloakIdentityProviderPair struct {
	KcIdentityProvider   *KeycloakIdentityProvider
	SpecIdentityProvider *KeycloakIdentityProvider
}

type KeycloakUser struct {
	*KeycloakApiUser
	OutputSecret string `json:"outputSecret, omitempty"`
}

type KeycloakApiUser struct {
	ID              string              `json:"id,omitempty"`
	UserName        string              `json:"username,omitempty"`
	FirstName       string              `json:"firstName"`
	LastName        string              `json:"lastName"`
	Email           string              `json:"email,omitempty"`
	EmailVerified   bool                `json:"emailVerified"`
	Enabled         bool                `json:"enabled"`
	RealmRoles      []string            `json:"realmRoles,omitempty"`
	ClientRoles     map[string][]string `json:"clientRoles"`
	RequiredActions []string            `json:"requiredActions,omitempty"`
	Groups          []string            `json:"groups,omitempty"`
}

type KeycloakUserPair struct {
	KcUser   *KeycloakUser
	SpecUser *KeycloakUser
}

type KeycloakProtocolMapper struct {
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Protocol        string            `json:"protocol,omitempty"`
	ProtocolMapper  string            `json:"protocolMapper,omitempty"`
	ConsentRequired bool              `json:"consentRequired,omitempty"`
	ConsentText     string            `json:"consentText"`
	Config          map[string]string `json:"config"`
}

type KeycloakClient struct {
	*KeycloakApiClient
	OutputSecret string `json:"outputSecret, omitempty"`
}

type KeycloakApiClient struct {
	ID                        string                   `json:"id,omitempty"`
	ClientID                  string                   `json:"clientId,omitempty"`
	Secret                    string                   `json:"secret"`
	Name                      string                   `json:"name"`
	BaseURL                   string                   `json:"baseUrl"`
	AdminURL                  string                   `json:"adminUrl"`
	RootURL                   string                   `json:"rootUrl"`
	Description               string                   `json:"description"`
	SurrogateAuthRequired     bool                     `json:"surrogateAuthRequired"`
	Enabled                   bool                     `json:"enabled"`
	ClientAuthenticatorType   string                   `json:"clientAuthenticatorType"`
	DefaultRoles              []string                 `json:"defaultRoles,omitempty"`
	RedirectUris              []string                 `json:"redirectUris,omitempty"`
	WebOrigins                []string                 `json:"webOrigins,omitempty"`
	NotBefore                 int                      `json:"notBefore"`
	BearerOnly                bool                     `json:"bearerOnly"`
	ConsentRequired           bool                     `json:"consentRequired"`
	StandardFlowEnabled       bool                     `json:"standardFlowEnabled"`
	ImplicitFlowEnabled       bool                     `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled bool                     `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled    bool                     `json:"serviceAccountsEnabled"`
	PublicClient              bool                     `json:"publicClient"`
	FrontchannelLogout        bool                     `json:"frontchannelLogout"`
	Protocol                  string                   `json:"protocol,omitempty"`
	Attributes                map[string]string        `json:"attributes,omitempty"`
	FullScopeAllowed          bool                     `json:"fullScopeAllowed"`
	NodeReRegistrationTimeout int                      `json:"nodeReRegistrationTimeout"`
	ProtocolMappers           []KeycloakProtocolMapper `json:"protocolMappers,omitempty"`
	UseTemplateConfig         bool                     `json:"useTemplateConfig"`
	UseTemplateScope          bool                     `json:"useTemplateScope"`
	UseTemplateMappers        bool                     `json:"useTemplateMappers"`
	Access                    map[string]bool          `json:"access"`
}
type KeycloakClientPair struct {
	KcClient   *KeycloakClient
	SpecClient *KeycloakClient
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type GenericStatus struct {
	Phase    StatusPhase `json:"phase"`
	Message  string      `json:"message"`
	Attempts int         `json:"attempts"`
	// marked as true when all work is done on it
	Ready bool `json:"ready"`
}

type KeycloakStatus struct {
	GenericStatus
}

type StatusPhase string

var (
	NoPhase                 StatusPhase = ""
	PhaseAccepted           StatusPhase = "accepted"
	PhaseComplete           StatusPhase = "complete"
	PhaseFailed             StatusPhase = "failed"
	PhaseModified           StatusPhase = "modified"
	PhaseProvisioned        StatusPhase = "provisioned"
	PhaseWaitForPodsToRun   StatusPhase = "waitingForPods"
	PhaseDeprovisioning     StatusPhase = "deprovisioning"
	PhaseDeprovisioned      StatusPhase = "deprovisioned"
	PhaseDeprovisionFailed  StatusPhase = "deprovisionFailed"
	PhaseCredentialsPending StatusPhase = "credentialsPending"
	PhaseProvision          StatusPhase = "provision"
)

// GetFinalizers gets the list of finalizers on obj
func GetFinalizers(obj runtime.Object) ([]string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	return accessor.GetFinalizers(), nil
}

// AddFinalizer adds value to the list of finalizers on obj
func AddFinalizer(obj runtime.Object, value string) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	finalizers := sets.NewString(accessor.GetFinalizers()...)
	finalizers.Insert(value)
	accessor.SetFinalizers(finalizers.List())
	return nil
}

// RemoveFinalizer removes the given value from the list of finalizers in obj, then returns a new list
// of finalizers after value has been removed.
func RemoveFinalizer(obj runtime.Object, value string) ([]string, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	finalizers := sets.NewString(accessor.GetFinalizers()...)
	finalizers.Delete(value)
	newFinalizers := finalizers.List()
	accessor.SetFinalizers(newFinalizers)
	return newFinalizers, nil
}
