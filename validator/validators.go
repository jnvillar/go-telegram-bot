package validator

var maxMsgLength = 50

func Length(msg string) (bool, string) {
	validation := len(msg) < maxMsgLength
	if !validation{
		return false, "Alguno de los parametros tiene una longitud mayor a 50"
	}
	return true, "ok"
}

func LenghOfParameters(params []string) (bool, string){
	for _, s := range params {
		v, err := Length(s)
		if !v{
			return false, err
		}
	}
	return true, "ok"
}
