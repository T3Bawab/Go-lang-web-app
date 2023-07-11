package forms

type errors map[string][]string

func (e errors) AddError(IdTag, massage string) {
	e[IdTag] = append(e[IdTag], massage)

}

func (e errors) GetError(IdTag string) string {
	errorList, ok := e[IdTag]

	if !ok || len(errorList) == 0 {
		return ""
	}

	return errorList[0]
}
