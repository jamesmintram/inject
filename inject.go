package inject

import "reflect"

//service map
type ServiceMap map[string]interface{}
type statusMap map[string]bool

func create(
	serviceMap ServiceMap,
	status statusMap,
	instType reflect.Type) interface{} {

	newInst := reflect.New(instType).Interface()
	doInject(serviceMap, status, newInst)
	return newInst
}

func doInject(
	serviceMap ServiceMap,
	status statusMap,
	instance interface{}) {

	var instType = reflect.TypeOf(instance).Elem()
	var instVal = reflect.ValueOf(instance).Elem()

	if len(instType.Name()) > 0 && status[instType.Name()] {
		return
	}

	for idx := 0; idx < instType.NumField(); idx++ {
		var field = instType.Field(idx)
		var fieldTag = field.Tag

		if _, ok := fieldTag.Lookup("inject"); ok {
			var dep = field.Type.Elem()
			var depName = dep.Name()
			var service = serviceMap[depName]

			if service == nil {
				service = create(serviceMap, status, dep)
				serviceMap[depName] = service

			}

			//Assign the field here
			var fieldVal = instVal.FieldByName(field.Name)
			if fieldVal.CanSet() {
				fieldVal.Set(reflect.ValueOf(service))
			} else {
				//TODO: Some sort of error?
			}
		}
	}

	if len(instType.Name()) > 0 {
		status[instType.Name()] = true
	}
}

func Inject(instance interface{}) {
	var statusMap = make(map[string]bool)
	var serviceMap = make(map[string]interface{})

	doInject(serviceMap, statusMap, instance)
}

func InjectAll(services ...interface{}) {
	var statusMap = make(map[string]bool)
	var serviceMap = make(map[string]interface{})

	for _, inst := range services {
		var instType = reflect.TypeOf(inst).Elem()
		if len(instType.Name()) > 0 {
			serviceMap[instType.Name()] = inst
		}
	}

	for _, service := range services {
		doInject(serviceMap, statusMap, service)
	}
}
