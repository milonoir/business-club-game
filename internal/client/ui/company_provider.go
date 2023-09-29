package ui

import (
	"sync"
)

var (
	colors = []string{
		"blue",
		"orange",
		"yellow",
		"red",
	}
)

type CompanyProvider struct {
	mux       sync.Mutex
	companies []string
}

func NewCompanyProvider() *CompanyProvider {
	return &CompanyProvider{
		companies: make([]string, 0, 4),
	}
}

func (p *CompanyProvider) SetCompanies(companies []string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.companies = companies
}

func (p *CompanyProvider) Companies() []string {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.companies
}

func (p *CompanyProvider) CompanyByIndex(index int) string {
	p.mux.Lock()
	defer p.mux.Unlock()
	if index < 0 || index > len(p.companies)-1 {
		return ""
	}
	return p.companies[index]
}

func (p *CompanyProvider) ColorByIndex(index int) string {
	p.mux.Lock()
	defer p.mux.Unlock()
	if index < 0 || index > len(p.companies)-1 {
		return ""
	}
	return colors[index]
}
