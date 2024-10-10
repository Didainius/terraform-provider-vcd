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

func createResource[O updateDeleter[O, I], I any](ctx context.Context, d *schema.ResourceData, meta interface{}, c crudConfig[O, I]) diag.Diagnostics {
	t, err := c.getTypeFunc(d)
	if err != nil {
		return diag.Errorf("error getting %s type: %s", c.entityLabel, err)
	}

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

func updateResource[O updateDeleter[O, I], I any](ctx context.Context, d *schema.ResourceData, meta interface{}, c crudConfig[O, I]) diag.Diagnostics {
	t, err := c.getTypeFunc(d)
	if err != nil {
		return diag.Errorf("error getting %s type: %s", c.entityLabel, err)
	}

	retrievedEntity, err := c.getEntityFunc(d.Id())
	if err != nil {
		return diag.Errorf("error getting %s: %s", c.entityLabel, err)
	}

	_, err = retrievedEntity.Update(t)
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", c.entityLabel, err)
	}

	return c.readFunc(ctx, d, meta)
}

func readResource[O updateDeleter[O, I], I any](ctx context.Context, d *schema.ResourceData, meta interface{}, c crudConfig[O, I]) diag.Diagnostics {
	retrievedEntity, err := c.getEntityFunc(d.Id())
	if err != nil {
		return diag.Errorf("error getting %s: %s", c.entityLabel, err)
	}

	err = c.stateStoreFunc(d, retrievedEntity)
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", c.entityLabel, err)
	}

	return nil
}

func deleteResource[O updateDeleter[O, I], I any](ctx context.Context, d *schema.ResourceData, meta interface{}, c crudConfig[O, I]) diag.Diagnostics {
	retrievedEntity, err := c.getEntityFunc(d.Id())
	if err != nil {
		return diag.Errorf("error getting %s: %s", c.entityLabel, err)
	}

	err = retrievedEntity.Delete()
	if err != nil {
		return diag.Errorf("error storing %s to state: %s", c.entityLabel, err)
	}

	return nil
}
