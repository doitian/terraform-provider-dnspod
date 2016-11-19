package client

type Record struct {
	Id *interface{} `json:"id"`

	// Used in Record.Modify
	RecordId string `json:"record_id"`

	DomainId   string `json:"domain_id"`
	SubDomain  string `json:"sub_domain"`
	RecordType string `json:"record_type"`
	RecordLine string `json:"record_line"`
	Value      string `json:"value"`
	Mx         string `json:"mx"`
	Ttl        string `json:"ttl"`
	Weight     string `json:"weight"`
}

type Domain struct {
	Id     string `json:"id"`
	Domain string `json:"domain"`

	// In Domain.Info, the domain is returned as name instead of domain
	Name string `json:"name"`
}

type DomainCreateRequest struct {
	Domain string `json:"domain"`
}

type DomainCreateResponse struct {
	GeneralResponse

	Domain Domain `json:"domain"`
}

type DomainInfoRequest struct {
	DomainId string `json:"domain_id"`
}

type DomainInfoResponse struct {
	GeneralResponse

	Domain Domain `json:"domain"`
}

type DomainRemoveRequest struct {
	DomainId string `json:"domain_id"`
}

type DomainRemoveResponse struct {
	GeneralResponse
}

type RecordCreateRequest Record

type RecordCreateResponse struct {
	GeneralResponse

	Record Record `json:"record"`
}

type RecordModifyRequest Record

type RecordModifyResponse struct {
	GeneralResponse

	Record Record `json:"record"`
}

type RecordInfoRequest struct {
	DomainId string `json:"domain_id"`
	RecordId string `json:"record_id"`
}

type RecordInfoResponse struct {
	GeneralResponse

	Record Record `json:"record"`
}

type RecordRemoveRequest struct {
	DomainId string `json:"domain_id"`
	RecordId string `json:"record_id"`
}

type RecordRemoveResponse struct {
	GeneralResponse
}
