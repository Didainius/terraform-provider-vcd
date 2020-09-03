/*
 * Copyright 2020 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// NsxtTier0Router
type NsxtTier0Router struct {
	NsxtTier0Router *types.NsxtTier0Router
	client          *Client
}

// GetNsxtTier0RouterByName retrieves NSX-T tier 0 router by given parent NSX-T manager ID and Tier 0 router name
//
// Note. NSX-T manager ID is mandatory and must be in URN format (e.g.
// urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc)
func (vcdCli *VCDClient) GetNsxtTier0RouterByName(name, nsxtManagerId string) (*NsxtTier0Router, error) {
	if nsxtManagerId == "" {
		return nil, fmt.Errorf("no NSX-T manager ID specified")
	}

	if !isUrn(nsxtManagerId) {
		return nil, fmt.Errorf("NSX-T manager ID is not URN (e.g. 'urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc)', got: %s", nsxtManagerId)
	}

	if name == "" {
		return nil, fmt.Errorf("empty Tier 0 router name specified")
	}

	queryParameters := copyOrNewUrlValues(nil)
	queryParameters.Add("filter", "displayName=="+name)

	nsxtTier0Routers, err := vcdCli.GetAllNsxtTier0Routers(nsxtManagerId, queryParameters)
	if err != nil {
		return nil, fmt.Errorf("could not find NSX-T Tier-0 router with name '%s' for NSX-T manager with id '%s': %s",
			name, nsxtManagerId, err)
	}

	if len(nsxtTier0Routers) == 0 {
		return nil, fmt.Errorf("no NSX-T Tier-0 router with name '%s' for NSX-T manager with id '%s' found", name, nsxtManagerId)
	}

	if len(nsxtTier0Routers) > 1 {
		return nil, fmt.Errorf("more than one (%d) NSX-T Tier-0 router with name '%s' for NSX-T manager with id '%s' found",
			len(nsxtTier0Routers), name, nsxtManagerId)
	}

	return nsxtTier0Routers[0], nil
}

// GetAllNsxtTier0Routers retrieves all NSX-T Tier-0 routers using OpenAPI endpoint. Query parameters can be supplied to
// perform additional filtering. By default it injects FIQL filter _context==nsxtManagerId (e.g.
// _context==urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc) because it is mandatory to list child Tier-0
// routers.
//
// Note. IDs of Tier-0 routers do not have a standard and may look as strings when they are created using UI or as UUIDs
// when they are created using API
func (vcdCli *VCDClient) GetAllNsxtTier0Routers(nsxtManagerId string, queryParameters url.Values) ([]*NsxtTier0Router, error) {
	if !isUrn(nsxtManagerId) {
		return nil, fmt.Errorf("NSX-T manager ID is not URN (e.g. 'urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc)', got: %s", nsxtManagerId)
	}

	endpoint := types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointImportableTier0Routers
	minimumApiVersion, err := vcdCli.Client.checkOpenApiEndpointCompatibility(endpoint)
	if err != nil {
		return nil, err
	}

	urlRef, err := vcdCli.Client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	// Get all Tier-0 routers that are accessible to an organization VDC. Routers that are already associated with an
	// External Network are filtered out. The “_context” filter key must be set with the id of the NSX-T manager for which
	// we want to get the Tier-0 routers for.
	//
	// _context==urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc

	// Create a copy of queryParameters so that original queryParameters are not mutated (because a map is always a
	// reference)
	queryParams := queryParameterFilterAnd("_context=="+nsxtManagerId, queryParameters)

	typeResponses := []*types.NsxtTier0Router{{}}
	err = vcdCli.Client.OpenApiGetAllItems(minimumApiVersion, urlRef, queryParams, &typeResponses)
	if err != nil {
		return nil, err
	}

	returnObjects := make([]*NsxtTier0Router, len(typeResponses))
	for sliceIndex := range typeResponses {
		returnObjects[sliceIndex] = &NsxtTier0Router{
			NsxtTier0Router: typeResponses[sliceIndex],
			client:          &vcdCli.Client,
		}
	}

	return returnObjects, nil
}

// copyOrNewUrlValues either creates a copy of parameters or instantiates a new url.Values if nil parameters are
// supplied. It helps to avoid mutating supplied parameter when additional values must be injected internally.
func copyOrNewUrlValues(parameters url.Values) url.Values {
	parameterCopy := make(map[string][]string)

	// if supplied parameters are nil - we just return new initialized
	if parameters == nil {
		return parameterCopy
	}

	// Copy URL values
	for key, value := range parameters {
		parameterCopy[key] = value
	}

	return parameterCopy
}

// queryParameterFilterAnd is a helper to append "AND" clause to FIQL filter by using ';' (semicolon) if any values are
// already set in 'filter' value of parameters. If none existed before then 'filter' value will be set.
//
// Note. It does a copy of supplied 'parameters' value and does not mutate supplied original parameter.
func queryParameterFilterAnd(filter string, parameters url.Values) url.Values {
	newParameters := copyOrNewUrlValues(parameters)

	existingFilter := newParameters.Get("filter")
	if existingFilter == "" {
		newParameters.Set("filter", filter)
		return newParameters
	}

	newParameters.Set("filter", existingFilter+";"+filter)
	return newParameters
}

// isUuid returns true if the identifier is a bare UUID
func isUuid(identifier string) bool {
	reUuid := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	return reUuid.MatchString(identifier)
}

// isUrn validates if supplied identifier is of URN format (e.g. urn:vcloud:nsxtmanager:09722307-aee0-4623-af95-7f8e577c9ebc)
// it checks for the following criteria:
// 1. idenfifier is not empty
// 2. identifier has 4 elements separated by ';'
// 3. element 1 is 'urn' and element 4 is valid UUID
func isUrn(identifier string) bool {
	if identifier == "" {
		return false
	}

	ss := strings.Split(identifier, ":")
	if len(ss) != 4 {
		return false
	}

	if ss[0] != "urn" && !isUuid(ss[3]) {
		return false
	}

	return true
}
