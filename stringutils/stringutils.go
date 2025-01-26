package stringutils 

func cleanInput(text string) []string {
    lc := strings.ToLower(text)
    trLc := strings.TrimSpace(lc)
    return strings.Fields(trLc)
}
