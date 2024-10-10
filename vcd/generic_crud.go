package vcd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type updateDeleter[O any, I any] interface {
	// *O
	Update(*I) (O, error)
	Delete() error
}

type crudConfig[O updateDeleter[O, I], I any] struct {
	// Mandatory parameters

	// entityLabel contains friendly entity name that is used for logging meaningful errors
	entityLabel string

	// Create
	getTypeFunc    func(d *schema.ResourceData) (*I, error)
	createFunc     func(config *I) (O, error)
	stateStoreFunc func(d *schema.ResourceData, outerType O) error

	readFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics

	// // Update
	getEntityFunc func(id string) (O, error)

	// Read

	// Delete
}

func createRes[O updateDeleter[O, I], I any](ctx context.Context, d *schema.ResourceData, meta interface{}, c crudConfig[O, I]) diag.Diagnostics {
	t, err := c.getTypeFunc(d)
	if err != nil {
		return diag.Errorf("error getting %s type: %s", c.entityLabel, err)
	}

	///
	createdEntity, err := c.createFunc(t)
	if err != nil {
		return diag.Errorf("error creating %s: %s", c.entityLabel, err)
	}

	err = c.stateStoreFunc(d, createdEntity)
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", c.entityLabel, err)
	}

	return c.readFunc(ctx, d, meta)
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

// func create[I, O any](ctx context.Context,
// 	d *schema.ResourceData,
// 	meta interface{},
// 	entityLabel string,
// 	getTypeFunc func(d *schema.ResourceData) (*I, error),
// 	createFunc func(config *I) (*O, error),
// 	stateStoreFunc func(d *schema.ResourceData, outerType *O) error,
// 	readFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics) diag.Diagnostics {

// 	t, err := getTypeFunc(d)
// 	if err != nil {
// 		return diag.Errorf("error getting %s type: %s", entityLabel, err)
// 	}
// 	///
// 	createdEntity, err := createFunc(t)
// 	if err != nil {
// 		return diag.Errorf("error creating %s: %s", entityLabel, err)
// 	}

// 	err = stateStoreFunc(d, createdEntity)
// 	if err != nil {
// 		return diag.Errorf("error storing %s to state: %s", entityLabel, err)
// 	}

// 	return readFunc(ctx, d, meta)
// }
