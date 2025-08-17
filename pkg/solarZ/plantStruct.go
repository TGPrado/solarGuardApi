package solarz

type ApiResponse struct {
	Finished                 bool       `json:"finished"`
	Error                    bool       `json:"error"`
	FirstTime                bool       `json:"firstTime"`
	Credencial               Credencial `json:"credencial"`
	Data                     Data       `json:"data"`
	RecommendedLimitExceeded bool       `json:"recommendedLimitExceeded"`
	RecommendedLimitWarning  *string    `json:"recommendedLimitWarning"`
	Importadas               int        `json:"importadas"`
	Total                    int        `json:"total"`
	IntegradorFreeSpace      int64      `json:"integradorFreeSpace"`
	UpdatedAt                string     `json:"updatedAt"`
}

type Credencial struct {
	UserApi           string  `json:"userApi"`
	APIName           string  `json:"apiName"`
	URLIcon           string  `json:"urlIcon"`
	ID                int     `json:"id"`
	IDApi             int     `json:"idApi"`
	Blocked           bool    `json:"blocked"`
	APIPossuiSuporte  bool    `json:"apiPossuiSuporte"`
	DeprecatedMessage *string `json:"deprecatedMessage"`
}

type Data struct {
	Content          []Content `json:"content"`
	Pageable         Pageable  `json:"pageable"`
	Last             bool      `json:"last"`
	TotalElements    int       `json:"totalElements"`
	TotalPages       int       `json:"totalPages"`
	First            bool      `json:"first"`
	Size             int       `json:"size"`
	Number           int       `json:"number"`
	Sort             Sort      `json:"sort"`
	NumberOfElements int       `json:"numberOfElements"`
	Empty            bool      `json:"empty"`
}

type Content struct {
	ID                           int64   `json:"id"`
	IDExterno                    string  `json:"idExterno"`
	ExternalUser                 *string `json:"externalUser"`
	Denominacao                  string  `json:"denominacao"`
	DenominacaoSolarz            *string `json:"denominacaoSolarz"`
	ImportadaPeloUsuarioAtual    bool    `json:"importadaPeloUsuarioAtual"`
	ImportadaPelaCredencialAtual bool    `json:"importadaPelaCredencialAtual"`
	Importando                   bool    `json:"importando"`
	UsinaJaImportadaID           *int64  `json:"usinaJaImportadaId"`
	HasWarning                   bool    `json:"hasWarning"`
	WarningTitle                 *string `json:"warningTitle"`
	WarningTooltip               *string `json:"warningTooltip"`
}

type Pageable struct {
	Sort       Sort `json:"sort"`
	Offset     int  `json:"offset"`
	PageSize   int  `json:"pageSize"`
	PageNumber int  `json:"pageNumber"`
	Paged      bool `json:"paged"`
	Unpaged    bool `json:"unpaged"`
}

type Sort struct {
	Empty    bool `json:"empty"`
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
}
