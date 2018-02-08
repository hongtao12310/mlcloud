package models

// LDAP ...
type LDAP struct {
    URL            string `json:"url"`
    SearchDN       string `json:"search_dn"`
    SearchPassword string `json:"search_password"`
    BaseDN         string `json:"base_dn"`
    Filter         string `json:"filter"`
    UID            string `json:"uid"`
    Scope          int    `json:"scope"`
    Timeout        int    `json:"timeout"` // in second
}

// Database ...
type Database struct {
    Type   string  `json:"type"`
    MySQL  *MySQL  `json:"mysql,omitempty"`
}

// MySQL ...
type MySQL struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password,omitempty"`
    Database string `json:"database"`
}

// SQLite ...
type SQLite struct {
    File string `json:"file"`
}

// Email ...
type Email struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    SSL      bool   `json:"ssl"`
    Identity string `json:"identity"`
    From     string `json:"from"`
}


