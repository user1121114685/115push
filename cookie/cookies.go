package cookie

type Cookie []struct {
	Domain   string `json:"domain"`
	HostOnly bool   `json:"hostOnly"`
	HTTPOnly bool   `json:"httpOnly"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	SameSite string `json:"sameSite"`
	Secure   bool   `json:"secure"`
	Session  bool   `json:"session"`
	StoreID  string `json:"storeId"`
	Value    string `json:"value"`
	ID       int    `json:"id"`
}
