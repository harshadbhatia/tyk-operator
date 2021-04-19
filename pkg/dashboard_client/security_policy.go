package dashboard_client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	v1 "github.com/TykTechnologies/tyk-operator/api/v1alpha1"
	"github.com/TykTechnologies/tyk-operator/pkg/universal_client"
)

type SecurityPolicy struct {
	*Client
}

// All Returns all policies from the Dashboard
func (p SecurityPolicy) All(ctx context.Context) ([]v1.SecurityPolicySpec, error) {
	res, err := p.Client.Get(ctx, toURL(endpointPolicies), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Returned error: %d", res.StatusCode)
	}

	var response PoliciesResponse
	if err := universal_client.JSON(res, &response); err != nil {
		return nil, err
	}
	return response.Policies, nil
}

// Get  find the Policy by id
func (p SecurityPolicy) Get(ctx context.Context, id string) (*v1.SecurityPolicySpec, error) {
	res, err := p.Client.Get(ctx, toURL(endpointPolicies, id), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var o v1.SecurityPolicySpec
	if err := universal_client.JSON(res, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// Create  creates a new policy using the def object
func (p SecurityPolicy) Create(ctx context.Context, def *v1.SecurityPolicySpec) error {
	res, err := p.Client.PostJSON(ctx, toURL(endpointPolicies), def)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return universal_client.Error(res)
	}
	var msg ResponseMsg
	if err := universal_client.JSON(res, &msg); err != nil {
		return err
	}
	switch strings.ToLower(msg.Status) {
	case "ok":
		def.MID = msg.Message
		return nil
	default:
		return universal_client.Error(res)
	}
}

// Update updates a resource object def
func (p SecurityPolicy) Update(ctx context.Context, def *v1.SecurityPolicySpec) error {
	res, err := p.Client.PutJSON(ctx, toURL(endpointPolicies, def.MID), def)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return universal_client.Error(res)
	}
	return universal_client.JSON(res, def)
}

// Delete deletes the resource by ID
func (p SecurityPolicy) Delete(ctx context.Context, id string) error {
	res, err := p.Client.Delete(ctx, toURL(endpointPolicies, id), nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return universal_client.Error(res)
	}
	return nil
}
