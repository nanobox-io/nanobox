package server



var registeredRPCs = []interface{}{}

func Register(i interface{}) {
	registeredRPCs = append(registeredRPCs, i)	
}
