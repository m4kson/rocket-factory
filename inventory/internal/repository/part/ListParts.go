package part

import (
	"context"

	"github.com/m4kson/rocket-factory/inventory/internal/model"
	"github.com/m4kson/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]model.Part, 0, len(r.parts))
	for _, p := range r.parts {
		if len(filter.Ids) > 0 {
			matched := false
			for _, u := range filter.Ids {
				if p.PartId == u {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// names
		if len(filter.Names) > 0 {
			matched := false
			for _, n := range filter.Names {
				if p.Name == n {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// categories
		if len(filter.Categories) > 0 {
			matched := false
			for _, c := range filter.Categories {
				if p.Category == converter.CategoryToRepoModel(c) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// manufacturer_countries
		if len(filter.ManufacturerCountries) > 0 {
			matched := false
			for _, mc := range filter.ManufacturerCountries {
				if p.Manufacturer.Country == mc {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// tags (ищем пересечение тегов)
		if len(filter.Tags) > 0 {
			matched := false
			for _, ft := range filter.Tags {
				for _, pt := range p.Tags {
					if pt == ft {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
			if !matched {
				continue
			}
		}

		results = append(results, converter.PartToModel(p))
	}

	return results, nil
}
