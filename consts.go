package zpl

const PDF = "pdf"
const PNG = "png"
const ZPL = "zpl"

func allowedDensities() []int {
	return []int{6, 8, 12, 24}
}

func allowedOutputFormats() []string {
	return []string{PDF, PNG, ZPL}
}
