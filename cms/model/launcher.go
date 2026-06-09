package model

type Build struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	ModLoader string `json:"mod_loader"`
	MCVersion string `json:"mc_version"`
	ServerID  string `json:"server_id"`
	FileHash  string `json:"file_hash"`
	FileSize  int64  `json:"file_size"`
	Changelog string `json:"changelog"`
	CreatedAt string `json:"created_at"`
}
