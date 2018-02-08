package common


// system related configurations

const (
    ExtEndpoint                = "ext_endpoint"
    AUTHMode                   = "auth_mode"
    DatabaseType               = "database_type"

    // MySQL related
    MySQLHost                  = "mysql_host"
    MySQLPort                  = "mysql_port"
    MySQLUsername              = "mysql_username"
    MySQLPassword              = "mysql_password"
    MySQLDatabase              = "mysql_database"

    // LDAP related
    LDAPURL                    = "ldap_url"
    LDAPSearchDN               = "ldap_search_dn"
    LDAPSearchPwd              = "ldap_search_password"
    LDAPBaseDN                 = "ldap_base_dn"
    LDAPUID                    = "ldap_uid"
    LDAPFilter                 = "ldap_filter"
    LDAPScope                  = "ldap_scope"
    LDAPTimeout                = "ldap_timeout"
    TokenServiceURL            = "token_service_url"
    RegistryURL                = "registry_url"

    // Email related
    EmailHost                  = "email_host"
    EmailPort                  = "email_port"
    EmailUsername              = "email_username"
    EmailPassword              = "email_password"
    EmailFrom                  = "email_from"
    EmailSSL                   = "email_ssl"
    EmailIdentity              = "email_identity"
    ProjectCreationRestriction = "project_creation_restriction"
    VerifyRemoteCert           = "verify_remote_cert"
    MaxJobWorkers              = "max_job_workers"
    TokenExpiration            = "token_expiration"
    CfgExpiration              = "cfg_expiration"
    JobLogDir                  = "job_log_dir"

    // cluster config
    KubeAPIServer          = "kube_apiserver"
    KubeClusterName            = "cluster"

    // used to store the user cert files
    KubeCertDir                = "kube_cert_dir"

    // in-cluster config or out-of-cluster config
    KubeInCluster              = "kube_in_cluster"

    AdminKubeConfigPath        = "admin_kubeconfig"

    // define file system base path
    FSBasePath                 = "fs_base_path"

)
