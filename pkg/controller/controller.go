package controller

func Run(p4InfoFilePath string, bmv2JsonFilePath string) error {
	_, err := NewP4InfoHelper(p4InfoFilePath)
	if err != nil {
		return err
	}
	return nil
}
