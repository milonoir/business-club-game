package ui

type CompanyProvider struct {
	companies     []string
	company2Color map[string]string
}

func NewCompanyProvider() *CompanyProvider {
	return &CompanyProvider{
		companies:     make([]string, 0, 4),
		company2Color: make(map[string]string, 4),
	}
}

func (p *CompanyProvider) AddCompany(name, color string) {
	p.companies = append(p.companies, name)
	p.company2Color[name] = color
}

func (p *CompanyProvider) Companies() []string {
	return p.companies
}

func (p *CompanyProvider) CompanyByIndex(index int) string {
	if index < 0 || index > len(p.companies)-1 {
		return ""
	}
	return p.companies[index]
}

func (p *CompanyProvider) ColorByCompanyIndex(index int) string {
	if index < 0 || index > len(p.companies)-1 {
		return "white"
	}
	return p.company2Color[p.companies[index]]
}

func (p *CompanyProvider) ColorByCompany(company string) string {
	return p.company2Color[company]
}
