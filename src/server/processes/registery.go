package processes

type ActivitiesRegistry struct {
	Active   []*Activity
	Inactive []*Activity
}

var CommonRegistry = NewActivitiesRegistry()

func NewActivitiesRegistry() *ActivitiesRegistry {
	return &ActivitiesRegistry{
		Active:   make([]*Activity, 0),
		Inactive: make([]*Activity, 0),
	}
}

func (r *ActivitiesRegistry) DisableApp(a *Activity) {
	newActive := make([]*Activity, 0, len(r.Active))
	for _, act := range r.Active {
		if act != a {
			newActive = append(newActive, act)
		}
	}
	r.Active = newActive

	r.Inactive = append(r.Inactive, a)
}

func (r *ActivitiesRegistry) ActivateApp(a *Activity) {
	newInactive := make([]*Activity, 0, len(r.Inactive))
	for _, act := range r.Inactive {
		if act != a {
			newInactive = append(newInactive, act)
		}
	}
	r.Inactive = newInactive

	r.Active = append(r.Active, a)
}
