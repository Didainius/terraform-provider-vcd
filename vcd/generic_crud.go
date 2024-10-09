package vcd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type crudConfig[O updateDeleter[O, I], I any] struct {
	// Mandatory parameters

	// entityLabel contains friendly entity name that is used for logging meaningful errors
	entityLabel string

	// Create
	getTypeFunc    func(d *schema.ResourceData) (*I, error)
	createFunc     func(config *I) (*O, error)
	stateStoreFunc func(d *schema.ResourceData, outerType *O) error

	readFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics

	// Update
	getEntityFunc func(id string) (O, error)

	// Read

	// Delete

	// // endpoint in the usual format (e.g. types.OpenApiPathVersion1_0_0 + types.OpenApiEndpointNsxtSegmentIpDiscoveryProfiles)
	// endpoint string

	// // Optional parameters

	// // endpointParams contains a slice of strings that will be used to construct the request URL. It will
	// // initially replace '%s' placeholders in the `endpoint` (if any) and will add them as suffix
	// // afterwards
	// endpointParams []string

	// // queryParameters will be passed as GET queries to the URL. Usually they are used for API filtering parameters
	// queryParameters url.Values
	// // additionalHeader can be used to pass additional headers for API calls. One of the common purposes is to pass
	// // tenant context
	// additionalHeader map[string]string
}

func create[I, O any](ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	entityLabel string,
	getTypeFunc func(d *schema.ResourceData) (*I, error),
	createFunc func(config *I) (*O, error),
	stateStoreFunc func(d *schema.ResourceData, outerType *O) error,
	readFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics) diag.Diagnostics {

	t, err := getTypeFunc(d)
	if err != nil {
		return diag.Errorf("error getting %s type: %s", entityLabel, err)
	}
	///
	createdEntity, err := createFunc(t)
	if err != nil {
		return diag.Errorf("error creating %s: %s", entityLabel, err)
	}

	err = stateStoreFunc(d, createdEntity)
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", entityLabel, err)
	}

	return readFunc(ctx, d, meta)
}

type updateDeleter[O, I any] interface {
	Update(*I) (O, error)
	Delete() error
}

type updater[O, I any] interface {
	Update(*I) (O, error)
}

func update[O updateDeleter[O, I], I any](ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	entityLabel string,
	getTypeFunc func(d *schema.ResourceData) (*I, error),
	getEntityFunc func(id string) (O, error),
	readFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics) diag.Diagnostics {

	t, err := getTypeFunc(d)
	if err != nil {
		return diag.Errorf("error getting %s type: %s", entityLabel, err)
	}
	///
	retrievedEntity, err := getEntityFunc(d.Id())
	if err != nil {
		return diag.Errorf("error getting %s: %s", entityLabel, err)
	}

	_, err = retrievedEntity.Update(t)
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", entityLabel, err)
	}

	return readFunc(ctx, d, meta)
}

func read[O any](ctx context.Context, d *schema.ResourceData, meta interface{}, entityLabel string, getEntityFunc func(id string) (*O, error), stateStoreFunc func(d *schema.ResourceData, outerType *O) error) diag.Diagnostics {
	retrievedEntity, err := getEntityFunc(d.Id())
	if err != nil {
		return diag.Errorf("error getting %s: %s", entityLabel, err)
	}

	err = stateStoreFunc(d, retrievedEntity)
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", entityLabel, err)
	}

	return nil
}

type deleter interface {
	Delete() error
}

func deleteRes[O updateDeleter[O, I], I any](ctx context.Context, d *schema.ResourceData, meta interface{}, entityLabel string, getEntityFunc func(id string) (O, error)) diag.Diagnostics {
	retrievedEntity, err := getEntityFunc(d.Id())
	if err != nil {
		return diag.Errorf("error getting %s: %s", entityLabel, err)
	}

	// if retrievedEntity == nil {
	// 	return diag.Errorf("error - nil entity %s retrieved", entityLabel)
	// }

	// err = (*retrievedEntity).Delete()
	err = retrievedEntity.Delete()
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", entityLabel, err)
	}

	return nil
}
