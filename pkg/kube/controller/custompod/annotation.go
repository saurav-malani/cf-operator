package custompod

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type annotationPredicate struct {
	Annotation string
}

func (p annotationPredicate) validAnnotation(annotations map[string]string) bool {
	for key := range annotations {
		if key == p.Annotation {
			return true
		}
	}
	return false
}

// Create implements Predicate
func (p annotationPredicate) Create(e event.CreateEvent) bool {
	if p.validAnnotation(e.Meta.GetAnnotations()) {
		return true
	}
	return false
}

// Delete implements Predicate
func (p annotationPredicate) Delete(e event.DeleteEvent) bool {
	if p.validAnnotation(e.Meta.GetAnnotations()) {
		return true
	}
	return false
}

// Update implements Predicate
func (p annotationPredicate) Update(e event.UpdateEvent) bool {
	if p.validAnnotation(e.MetaNew.GetAnnotations()) {
		return true
	}
	return false
}

// Generic implements Predicate
func (p annotationPredicate) Generic(e event.GenericEvent) bool {
	if p.validAnnotation(e.Meta.GetAnnotations()) {
		return true
	}
	return false
}
